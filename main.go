package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const version = "1.0.0"

// ── Styles ─────────────────────────────────────────────
var (
	// Header
	logoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00d4ff")).
			Background(lipgloss.Color("#1a1a2e")).
			Padding(0, 1)

	versionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555")).
			Padding(0, 1)

	// Current context
	currentLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888"))

	currentValueStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#50fa7b"))

	// Search
	searchActiveStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#f1fa8c"))

	searchPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#555")).
				Italic(true)

	// List items
	selectedItemStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00d4ff"))

	normalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#999"))

	activeItemStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#50fa7b"))

	// Decorations
	aliasStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	activeTag    = lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")).Render("●")
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#555"))
	successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555"))
	counterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#666"))

	// Box
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#333")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#555"))
)

// ── Config (aliases) ───────────────────────────────────
type config struct {
	Aliases map[string]string `json:"aliases"`
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ksw.json")
}

func loadConfig() config {
	c := config{Aliases: make(map[string]string)}
	data, err := os.ReadFile(configPath())
	if err != nil {
		return c
	}
	_ = json.Unmarshal(data, &c)
	if c.Aliases == nil {
		c.Aliases = make(map[string]string)
	}
	return c
}

func saveConfig(c config) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0644)
}

// ── Fuzzy matching ─────────────────────────────────────
type scored struct {
	index int
	score int
}

// fuzzyMatch returns a score > 0 if pattern fuzzy-matches str.
// Higher score = better match. 0 = no match.
func fuzzyMatch(str, pattern string) int {
	str = strings.ToLower(str)
	pattern = strings.ToLower(pattern)

	pLen := utf8.RuneCountInString(pattern)
	if pLen == 0 {
		return 1
	}

	sRunes := []rune(str)
	pRunes := []rune(pattern)
	sLen := len(sRunes)

	// Check if all pattern chars exist in order
	pi := 0
	for si := 0; si < sLen && pi < pLen; si++ {
		if sRunes[si] == pRunes[pi] {
			pi++
		}
	}
	if pi < pLen {
		return 0 // not all chars matched
	}

	// Score: bonus for consecutive matches, word boundary matches, and early matches
	score := 0
	pi = 0
	consecutive := 0
	for si := 0; si < sLen && pi < pLen; si++ {
		if sRunes[si] == pRunes[pi] {
			pi++
			consecutive++
			score += 10 + consecutive*5 // consecutive bonus

			// Word boundary bonus (after /, -, _, or start)
			if si == 0 || sRunes[si-1] == '/' || sRunes[si-1] == '-' || sRunes[si-1] == '_' {
				score += 20
			}
			// Early match bonus
			score += max(0, 5-si)
		} else {
			consecutive = 0
		}
	}

	// Exact substring bonus
	if strings.Contains(str, pattern) {
		score += 50
	}

	return score
}

// ── Kubeconfig helpers ─────────────────────────────────
func getContexts() ([]string, error) {
	cmd := exec.Command("kubectl", "config", "get-contexts", "-o", "name")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get contexts: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var contexts []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			contexts = append(contexts, l)
		}
	}
	return contexts, nil
}

func getCurrentContext() string {
	cmd := exec.Command("kubectl", "config", "current-context")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func switchContext(name string) error {
	cmd := exec.Command("kubectl", "config", "use-context", name)
	return cmd.Run()
}

// ── Model ──────────────────────────────────────────────
type model struct {
	contexts       []string
	filtered       []int
	cursor         int
	scrollOffset   int
	current        string
	chosen         string
	search         string
	cfg            config
	terminalHeight int
	quitting       bool
}

func initialModel(contexts []string, current string, cfg config) model {
	m := model{
		contexts:       contexts,
		current:        current,
		cfg:            cfg,
		terminalHeight: 24,
	}
	m.resetFilter()
	for i, idx := range m.filtered {
		if contexts[idx] == current {
			m.cursor = i
			break
		}
	}
	m.ensureVisible()
	return m
}

func (m *model) resetFilter() {
	m.filtered = make([]int, len(m.contexts))
	for i := range m.contexts {
		m.filtered[i] = i
	}
	m.scrollOffset = 0
}

func (m *model) applyFilter() {
	if m.search == "" {
		m.resetFilter()
		return
	}

	query := m.search

	// Build searchable strings: context name + any aliases pointing to it
	reverseAlias := make(map[string][]string)
	for alias, ctx := range m.cfg.Aliases {
		reverseAlias[ctx] = append(reverseAlias[ctx], alias)
	}

	var results []scored
	for i, ctx := range m.contexts {
		// Match against context name
		searchable := ctx
		if aliases, ok := reverseAlias[ctx]; ok {
			searchable += " " + strings.Join(aliases, " ")
		}
		score := fuzzyMatch(searchable, query)
		if score > 0 {
			results = append(results, scored{index: i, score: score})
		}
	}

	// Sort by score descending
	sort.Slice(results, func(a, b int) bool {
		return results[a].score > results[b].score
	})

	m.filtered = nil
	for _, r := range results {
		m.filtered = append(m.filtered, r.index)
	}
	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
}

func (m *model) maxVisible() int {
	headerLines := 8
	v := m.terminalHeight - headerLines - 2
	if v < 3 {
		v = 3
	}
	return v
}

func (m *model) ensureVisible() {
	mv := m.maxVisible()
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	} else if m.cursor >= m.scrollOffset+mv {
		m.scrollOffset = m.cursor - mv + 1
	}
}

func (m *model) aliasFor(ctx string) string {
	for alias, target := range m.cfg.Aliases {
		if target == ctx {
			return alias
		}
	}
	return ""
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEscape:
			if m.search != "" {
				m.search = ""
				m.resetFilter()
				m.cursor = 0
			} else {
				m.quitting = true
				return m, tea.Quit
			}
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
				m.ensureVisible()
			}
		case tea.KeyDown:
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
				m.ensureVisible()
			}
		case tea.KeyHome:
			m.cursor = 0
			m.ensureVisible()
		case tea.KeyEnd:
			m.cursor = max(0, len(m.filtered)-1)
			m.ensureVisible()
		case tea.KeyPgUp:
			m.cursor = max(0, m.cursor-10)
			m.ensureVisible()
		case tea.KeyPgDown:
			m.cursor = min(len(m.filtered)-1, m.cursor+10)
			m.ensureVisible()
		case tea.KeyEnter:
			if len(m.filtered) > 0 {
				m.chosen = m.contexts[m.filtered[m.cursor]]
				return m, tea.Quit
			}
		case tea.KeyBackspace:
			if len(m.search) > 0 {
				m.search = m.search[:len(m.search)-1]
				m.applyFilter()
			}
		case tea.KeyRunes:
			m.search += string(msg.Runes)
			m.applyFilter()
			m.cursor = 0
			m.scrollOffset = 0
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting || m.chosen != "" {
		return ""
	}

	var b strings.Builder

	// ── Current context ──
	currentAlias := m.aliasFor(m.current)
	currentLabel := m.current
	if currentAlias != "" {
		currentLabel += " " + aliasStyle.Render("@"+currentAlias)
	}
	b.WriteString("  " + currentLabelStyle.Render("  current ") + currentValueStyle.Render(currentLabel) + "\n")
	b.WriteString("\n")

	// ── Search bar ──
	if m.search != "" {
		b.WriteString("  " + searchActiveStyle.Render("  ❯ "+m.search+"█") + "\n")
	} else {
		b.WriteString("  " + searchPlaceholderStyle.Render("  ❯ type to search...") + "\n")
	}

	// ── Separator ──
	b.WriteString("  " + dimStyle.Render("  ─────────────────────────────────────────") + "\n")

	if len(m.filtered) == 0 {
		b.WriteString("\n  " + dimStyle.Render("  No matching contexts") + "\n")
		return b.String()
	}

	maxVisible := m.maxVisible()

	start := m.scrollOffset
	end := start + maxVisible
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	// ── Scroll indicator top ──
	if start > 0 {
		b.WriteString("  " + dimStyle.Render(fmt.Sprintf("    ▲ %d more", start)) + "\n")
	}

	// ── List ──
	for i := start; i < end; i++ {
		ctx := m.contexts[m.filtered[i]]
		isActive := ctx == m.current
		alias := m.aliasFor(ctx)

		pointer := "   "
		var name string

		if i == m.cursor {
			pointer = " ❯ "
			name = selectedItemStyle.Render(ctx)
		} else if isActive {
			name = activeItemStyle.Render(ctx)
		} else {
			name = normalItemStyle.Render(ctx)
		}

		extras := ""
		if alias != "" {
			extras += " " + aliasStyle.Render("@"+alias)
		}
		if isActive {
			extras += " " + activeTag
		}

		b.WriteString("  " + pointer + name + extras + "\n")
	}

	// ── Scroll indicator bottom ──
	if end < len(m.filtered) {
		b.WriteString("  " + dimStyle.Render(fmt.Sprintf("    ▼ %d more", len(m.filtered)-end)) + "\n")
	}

	// ── Footer ──
	b.WriteString("\n")
	b.WriteString("  " + counterStyle.Render(fmt.Sprintf("  %d/%d", len(m.filtered), len(m.contexts))) +
		helpStyle.Render("  ↑↓ navigate · enter select · esc clear · ctrl+c quit") + "\n")

	return b.String()
}

// ── Main ───────────────────────────────────────────────
func main() {
	cfg := loadConfig()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version":
			fmt.Printf("ksw v%s\n", version)
			return

		case "-h", "--help":
			fmt.Printf(`ksw v%s - Interactive Kubernetes context switcher

Usage:
  ksw                  Launch interactive selector (fuzzy search)
  ksw <name>           Switch directly to context <name>
  ksw @<alias>         Switch using an alias
  ksw alias <name> <context>  Create alias for a context
  ksw alias rm <name>         Remove an alias
  ksw alias ls                List all aliases
  ksw -l               List contexts (non-interactive)
  ksw -h               Show this help
  ksw -v               Show version

Navigation:
  Type                Filter contexts with fuzzy search
  ↑ / ↓              Move up / down
  Home / End          Go to top / bottom
  PgUp / PgDn         Jump 10 items
  Backspace           Delete last character from filter
  Enter               Switch to highlighted context
  Esc                 Clear filter / Quit
  Ctrl+C              Quit

Aliases are stored in ~/.ksw.json
`, version)
			return

		case "-l", "--list":
			contexts, err := getContexts()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			current := getCurrentContext()
			reverseAlias := make(map[string]string)
			for alias, ctx := range cfg.Aliases {
				reverseAlias[ctx] = alias
			}
			for _, ctx := range contexts {
				alias := ""
				if a, ok := reverseAlias[ctx]; ok {
					alias = aliasStyle.Render(" @" + a)
				}
				if ctx == current {
					fmt.Printf("%s%s %s\n", currentValueStyle.Render("▸ "+ctx), alias, activeTag)
				} else {
					fmt.Printf("  %s%s\n", ctx, alias)
				}
			}
			return

		case "alias":
			handleAlias(cfg)
			return

		default:
			arg := os.Args[1]

			// Handle @alias
			if strings.HasPrefix(arg, "@") {
				aliasName := arg[1:]
				target, ok := cfg.Aliases[aliasName]
				if !ok {
					fmt.Fprintf(os.Stderr, "%s Alias '%s' not found. Use 'ksw alias ls' to list.\n", warnStyle.Render("✗"), aliasName)
					os.Exit(1)
				}
				// Try exact match first, then suffix/substring match
				if err := switchContext(target); err != nil {
					contexts, cerr := getContexts()
					if cerr != nil {
						fmt.Fprintln(os.Stderr, cerr)
						os.Exit(1)
					}
					var matches []string
					for _, ctx := range contexts {
						if strings.HasSuffix(ctx, "/"+target) || strings.HasSuffix(ctx, target) || strings.Contains(ctx, target) {
							matches = append(matches, ctx)
						}
					}
					if len(matches) == 1 {
						target = matches[0]
						if err := switchContext(target); err != nil {
							fmt.Fprintf(os.Stderr, "%s Context '%s' (alias @%s) not found in kubeconfig.\n", warnStyle.Render("✗"), target, aliasName)
							os.Exit(1)
						}
					} else if len(matches) > 1 {
						fmt.Fprintf(os.Stderr, "%s Ambiguous alias @%s, matches:\n", warnStyle.Render("✗"), aliasName)
						for _, m := range matches {
							fmt.Fprintf(os.Stderr, "  %s\n", m)
						}
						os.Exit(1)
					} else {
						fmt.Fprintf(os.Stderr, "%s Context '%s' (alias @%s) not found in kubeconfig.\n", warnStyle.Render("✗"), target, aliasName)
						os.Exit(1)
					}
				}
				fmt.Printf("%s Switched to %s %s\n", successStyle.Render("✔"), target, aliasStyle.Render("@"+aliasName))
				return
			}

			if arg[0] != '-' {
				// Try exact match first, then suffix/substring match
				target := arg
				if err := switchContext(target); err != nil {
					// Exact match failed, try to find by suffix or substring
					contexts, cerr := getContexts()
					if cerr != nil {
						fmt.Fprintln(os.Stderr, cerr)
						os.Exit(1)
					}
					var matches []string
					for _, ctx := range contexts {
						if strings.HasSuffix(ctx, "/"+arg) || strings.HasSuffix(ctx, arg) || strings.Contains(ctx, arg) {
							matches = append(matches, ctx)
						}
					}
					if len(matches) == 1 {
						target = matches[0]
						if err := switchContext(target); err != nil {
							fmt.Fprintf(os.Stderr, "%s Context '%s' not found.\n", warnStyle.Render("✗"), target)
							os.Exit(1)
						}
					} else if len(matches) > 1 {
						fmt.Fprintf(os.Stderr, "%s Ambiguous context '%s', matches:\n", warnStyle.Render("✗"), arg)
						for _, m := range matches {
							fmt.Fprintf(os.Stderr, "  %s\n", m)
						}
						os.Exit(1)
					} else {
						fmt.Fprintf(os.Stderr, "%s Context '%s' not found.\n", warnStyle.Render("✗"), arg)
						os.Exit(1)
					}
				}
				fmt.Printf("%s Switched to %s\n", successStyle.Render("✔"), target)
				return
			}
			fmt.Fprintf(os.Stderr, "Unknown flag: %s. Use -h for help.\n", arg)
			os.Exit(1)
		}
	}

	// Interactive mode
	contexts, err := getContexts()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(contexts) == 0 {
		fmt.Fprintln(os.Stderr, "No contexts found in kubeconfig.")
		os.Exit(1)
	}

	current := getCurrentContext()
	m := initialModel(contexts, current, cfg)

	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	final := result.(model)
	if final.chosen != "" && final.chosen != current {
		if err := switchContext(final.chosen); err != nil {
			fmt.Fprintf(os.Stderr, "Error switching to %s: %v\n", final.chosen, err)
			os.Exit(1)
		}
		alias := final.aliasFor(final.chosen)
		extra := ""
		if alias != "" {
			extra = " " + aliasStyle.Render("@"+alias)
		}
		fmt.Printf("%s Switched to %s%s\n", successStyle.Render("✔"), final.chosen, extra)
	} else if final.chosen == current {
		fmt.Printf("%s Already on %s\n", dimStyle.Render("·"), current)
	}
}

func handleAlias(cfg config) {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ksw alias <ls|rm|name> [context]")
		os.Exit(1)
	}

	sub := os.Args[2]

	switch sub {
	case "ls", "list":
		if len(cfg.Aliases) == 0 {
			fmt.Println(dimStyle.Render("No aliases configured. Use: ksw alias <name> <context>"))
			return
		}
		// Sort aliases for consistent output
		names := make([]string, 0, len(cfg.Aliases))
		for name := range cfg.Aliases {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			fmt.Printf("  %s → %s\n", aliasStyle.Render("@"+name), cfg.Aliases[name])
		}

	case "rm", "remove", "delete":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "Usage: ksw alias rm <name>")
			os.Exit(1)
		}
		name := os.Args[3]
		if _, ok := cfg.Aliases[name]; !ok {
			fmt.Fprintf(os.Stderr, "%s Alias '%s' not found.\n", warnStyle.Render("✗"), name)
			os.Exit(1)
		}
		delete(cfg.Aliases, name)
		if err := saveConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s Removed alias %s\n", successStyle.Render("✔"), aliasStyle.Render("@"+name))

	default:
		// ksw alias <name> <context>
		name := sub
		if len(os.Args) < 4 {
			// Show what this alias points to
			if target, ok := cfg.Aliases[name]; ok {
				fmt.Printf("  %s → %s\n", aliasStyle.Render("@"+name), target)
			} else {
				fmt.Fprintf(os.Stderr, "Usage: ksw alias <name> <context>\n")
				os.Exit(1)
			}
			return
		}
		context := os.Args[3]
		cfg.Aliases[name] = context
		if err := saveConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s Alias %s → %s\n", successStyle.Render("✔"), aliasStyle.Render("@"+name), context)
	}
}
