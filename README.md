# ksw

**AI-powered** Kubernetes context switcher for your terminal. Built in Go.

Switch contexts with natural language, manage groups, pins and aliases — all by just telling the AI what you need. Or use the blazing-fast interactive TUI with fuzzy search. Single binary, no runtime dependencies.

> Available for **macOS** and **Linux** (amd64 & arm64). Named `ksw` to avoid conflict with macOS built-in `kswitch` (Kerberos).

## 🤖 AI — Natural Language Context Management

Talk to your clusters. `ksw ai` understands what you mean and executes it.

```bash
# Switch contexts with natural language
ksw ai "switch to payments dev"
# ✔ Switched to arn:aws:eks:us-east-1:111122223333:cluster/eks-payments-dev

# Create groups, pins, aliases — just ask
ksw ai "create a group called backend with payments and orders dev"
# ✔ Group 'backend' created (2 contexts)

ksw ai "pin nequi dev"
# ✔ Pinned ★ arn:.../eks-nequi-dev

# Ask questions about your setup
ksw ai "list my pins and groups as a table"
# (AI builds a formatted table from your current state)

# Conversational memory — it remembers context
ksw ai "switch to sufi qa"
ksw ai "now the same but in dev"
# ✔ Switched to arn:.../eks-sufi-dev

ksw ai "go back to the previous one"
# ✔ Switched to arn:.../eks-sufi-qa

# Multi-action — multiple tasks in one prompt
ksw ai "list my pins and groups"
# (executes both commands in a single call)

# Delete, rename, explore — anything you can do manually
ksw ai "delete the sufi group"
ksw ai "rename payments-dev to pay-dev"
ksw ai "what have I done so far?"
# (shows your recent history from conversational memory)
```

### Supported AI Providers

| Provider | Models | Auth |
|----------|--------|------|
| OpenAI | gpt-4o, gpt-4o-mini, etc. | API Key |
| Claude (Anthropic) | claude-sonnet-4-20250514, etc. | API Key |
| Gemini (Google) | gemini-2.0-flash, etc. | API Key |
| AWS Bedrock | Any Bedrock model (Claude, Llama, etc.) | AWS Profile, Access Keys, or Env vars |

### AI Configuration

```bash
# Interactive setup wizard
ksw ai config
# → Select provider (openai / claude / gemini / bedrock)
# → Choose model
# → Enter credentials
# → Done!
```

### AI Features

- **Natural language** — switch, create, delete, list, rename — just describe what you want
- **Conversational memory** — remembers your last 10 interactions, understands "the previous one", "same but in qa"
- **Multi-action** — execute multiple tasks in a single prompt
- **Smart formatting** — ask for tables, summaries, or any custom format
- **Response cache** — 30s TTL avoids duplicate LLM calls for repeated queries
- **Full state awareness** — AI knows your current context, groups, pins, aliases, and history
- **Pre-filtering** — extracts keywords locally to narrow candidates before calling the LLM
- **Retry with backoff** — handles rate limits (429) and server errors gracefully

## Install

### One-line installer (macOS & Linux — easiest)

```bash
curl -sL https://raw.githubusercontent.com/YonierGomez/kswitch/main/install.sh | bash
```

Automatically detects your OS and architecture (amd64/arm64) and installs to `/usr/local/bin`.

### Homebrew (macOS & Linux)

```bash
brew tap YonierGomez/kswitch
brew install kswitch
```

### Manual — Linux

```bash
# amd64 (x86_64)
curl -sL https://github.com/YonierGomez/kswitch/releases/latest/download/ksw-linux-amd64.tar.gz | tar xz
chmod +x ksw-linux-amd64
sudo mv ksw-linux-amd64 /usr/local/bin/ksw

# arm64 (AWS Graviton, Raspberry Pi, etc.)
curl -sL https://github.com/YonierGomez/kswitch/releases/latest/download/ksw-linux-arm64.tar.gz | tar xz
chmod +x ksw-linux-arm64
sudo mv ksw-linux-arm64 /usr/local/bin/ksw
```

### Manual — macOS

```bash
# Apple Silicon (M1/M2/M3)
curl -sL https://github.com/YonierGomez/kswitch/releases/latest/download/ksw-darwin-arm64.tar.gz | tar xz
chmod +x ksw-darwin-arm64
sudo mv ksw-darwin-arm64 /usr/local/bin/ksw

# Intel
curl -sL https://github.com/YonierGomez/kswitch/releases/latest/download/ksw-darwin-amd64.tar.gz | tar xz
chmod +x ksw-darwin-amd64
sudo mv ksw-darwin-amd64 /usr/local/bin/ksw
```

### From source

```bash
go install github.com/YonierGomez/kswitch@latest
```

## Usage

```bash
# ── AI (natural language) ──
ksw ai "<query>"             # AI-powered: switch, create, list, delete — anything
ksw ai config                # Configure AI provider and credentials

# ── Interactive TUI ──
ksw                          # Interactive selector (fuzzy search)
ksw <name>                   # Switch directly (short name ok: ksw payments-dev)
ksw -                        # Switch to previous context
ksw @<alias>                 # Switch using alias

# ── History ──
ksw history                  # Show recent context history
ksw history <n>              # Switch to history entry by number

# ── Groups ──
ksw group add <name> [ctx]   # Create a group and add contexts to it
ksw group rm <name>          # Remove a group
ksw group ls                 # List all groups with their members
ksw group use <name>         # Open TUI filtered to a group
ksw group add-ctx <g> <ctx>  # Add a context to an existing group
ksw group rmi <g> <ctx>      # Remove a context from a group

# ── Pins ──
ksw pin <name>               # Pin a context to the top of the list
ksw pin rm <name>            # Unpin a context
ksw pin ls                   # List pinned contexts
ksw pin use                  # Open TUI filtered to pinned contexts only

# ── Aliases & Rename ──
ksw alias <name> <context>   # Create alias for a context
ksw alias rm <name>          # Remove an alias
ksw alias ls                 # List all aliases
ksw rename <old> <new>       # Rename a context in kubeconfig

# ── Other ──
ksw completion install       # Auto-install shell completion (~/.zshrc or ~/.bashrc)
ksw completion zsh           # Print zsh setup line
ksw completion bash          # Print bash setup line
ksw -l                       # List contexts (non-interactive)
ksw -v                       # Version
ksw -h                       # Help
```

### Interactive TUI Navigation

| Key          | Action                              |
|--------------|-------------------------------------|
| Type         | Fuzzy filter in real time           |
| `↑` / `↓`   | Move up / down                      |
| `Home`/`End` | Go to top / bottom                  |
| `PgUp/PgDn`  | Jump 10 items                       |
| `Backspace`  | Delete filter character             |
| `Enter`      | Switch to highlighted context       |
| `Ctrl+P`     | Pin / unpin current context (★)     |
| `Ctrl+T`     | Jump to first pinned context        |
| `Ctrl+F`     | Toggle pinned-only filter `[★ pinned]` |
| `Ctrl+H`     | Toggle short name view (persisted)  |
| `Esc`        | Clear filter / Quit                 |
| `Ctrl+C`     | Quit                                |

### Short-name switching

You can switch to a context using just the cluster name, without the full ARN:

```bash
ksw eks-payments-dev
# ✔ Switched to arn:aws:eks:us-east-1:111122223333:cluster/eks-payments-dev
```

If the name is ambiguous, ksw will show all matches:

```bash
ksw payments
# ✗ Ambiguous context 'payments', matches:
#   arn:.../eks-payments-dev
#   arn:.../eks-payments-qa
```

### Pinned contexts

Pin your most-used contexts so they always appear at the top of the list:

```bash
ksw pin eks-payments-dev     # Pin by short name
ksw pin ls                   # List pinned contexts
ksw pin rm eks-payments-dev  # Unpin
```

In the TUI, pinned contexts appear in **yellow** with a `★` marker. Press `Ctrl+P` to toggle pin on the current item, and `Ctrl+T` to jump to the first pinned context from anywhere in the list.

### Previous context

Switch back to the last context instantly — like `cd -` in bash:

```bash
ksw -
# ✔ Switched to arn:.../eks-payments-qa

# Toggle back and forth between two contexts
ksw payments-dev   # switch to dev
ksw -              # back to qa
ksw -              # back to dev
```

### History

Show the last 10 contexts you visited:

```bash
ksw history
#   Recent contexts:
#   1  arn:.../eks-payments-dev ●
#   2  arn:.../eks-payments-qa
#   3  arn:.../eks-payments-pdn
#   4  arn:.../eks-orders-dev
#   5  arn:.../eks-orders-pdn
```

Switch directly to any history entry by number:

```bash
ksw history 3
# ✔ Switched to arn:.../eks-payments-pdn

ksw history 5
# ✔ Switched to arn:.../eks-orders-pdn
```

### Groups

Organize contexts into named groups and open the TUI filtered to only those contexts:

```bash
# Create a group with multiple contexts (short names or glob patterns)
ksw group add payments eks-payments-dev eks-payments-qa eks-payments-pdn
# ✔ Group payments — added 3 context(s)

# Use substring match (simplest — matches anywhere in the name)
ksw group add outposts outposts
# ✔ Group outposts — added 9 context(s)

# Use glob patterns — quote to prevent shell expansion
# pattern*  → auto-wraps to *pattern* (matches anywhere)
ksw group add payments "eks-payments*"
# ✔ Group payments — added 3 context(s)
#   · arn:.../eks-payments-dev
#   · arn:.../eks-payments-qa
#   · arn:.../eks-payments-pdn

# Open TUI showing only the payments group
ksw group use payments
# [payments] label shown in header, only 3 contexts visible

# List all groups
ksw group ls
# payments (3 contexts)
#   · arn:.../eks-payments-dev
#   · arn:.../eks-payments-qa
#   · arn:.../eks-payments-pdn

# Add a context to an existing group
ksw group add-ctx payments eks-payments-staging

# Remove a context from a group
ksw group rmi payments eks-payments-staging

# Remove a group entirely
ksw group rm payments
```

### Aliases

```bash
ksw alias prod arn:aws:eks:us-east-1:111122223333:cluster/eks-payments-dev
# ✔ Alias @prod → arn:.../eks-payments-dev

ksw @prod
# ✔ Switched to arn:.../eks-payments-dev @prod

# Aliases also work with short names:
ksw alias dev eks-payments-dev
ksw @dev
# ✔ Switched to arn:.../eks-payments-dev @dev
```

### Shell completion

```bash
ksw completion install   # Auto-installs in ~/.zshrc or ~/.bashrc
# ✔ Installed zsh completion in /Users/you/.zshrc
# Run: source ~/.zshrc
```

### Rename a context

```bash
ksw rename eks-payments-dev payments-dev
# ✔ Renamed arn:.../eks-payments-dev → payments-dev
```

## Configuration

All settings are stored in `~/.ksw.json`:

```json
{
  "aliases": {
    "prod": "arn:aws:eks:us-east-1:111122223333:cluster/eks-payments-dev"
  },
  "pins": [
    "arn:aws:eks:us-east-1:111122223333:cluster/eks-payments-dev"
  ],
  "history": [
    "arn:aws:eks:us-east-1:111122223333:cluster/eks-payments-dev",
    "arn:aws:eks:us-east-1:444455556666:cluster/eks-payments-qa"
  ],
  "previous": "arn:aws:eks:us-east-1:444455556666:cluster/eks-payments-qa",
  "ai": {
    "provider": "bedrock",
    "model": "us.anthropic.claude-sonnet-4-6",
    "region": "us-east-1",
    "auth_method": "profile",
    "profile": "my-aws-profile"
  },
  "ai_memory": [
    {
      "query": "switch to payments dev",
      "action": "switch",
      "result": "eks-payments-dev",
      "time": 1709312400
    }
  ]
}
```

## Requirements

- `kubectl` installed and configured
- For `ksw ai` with AWS Bedrock: `aws` CLI installed and configured

## License

MIT — Built by [Yonier Gómez](https://www.yonier.com)
