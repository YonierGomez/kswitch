# ksw

Interactive Kubernetes context switcher for your terminal. Built in Go.

Navigate with arrow keys, fuzzy search by typing, pin your favorites, and switch contexts in milliseconds. Single binary, no runtime dependencies.

> Named `ksw` to avoid conflict with macOS built-in `kswitch` (Kerberos).

## Install

### Homebrew (recommended)

```bash
brew tap YonierGomez/kswitch
brew install kswitch
```

### From source

```bash
go install github.com/YonierGomez/kswitch@latest
```

### Manual

Download the binary from [Releases](https://github.com/YonierGomez/kswitch/releases) and place it in your PATH.

## Usage

```bash
ksw                          # Interactive selector (fuzzy search)
ksw <name>                   # Switch directly (short name ok: ksw payments-dev)
ksw -                        # Switch to previous context
ksw @<alias>                 # Switch using alias
ksw history                  # Show recent context history
ksw pin <name>               # Pin a context to the top of the list
ksw pin rm <name>            # Unpin a context
ksw pin ls                   # List pinned contexts
ksw rename <old> <new>       # Rename a context in kubeconfig
ksw alias <name> <context>   # Create alias for a context
ksw alias rm <name>          # Remove an alias
ksw alias ls                 # List all aliases
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

Switch back to the last context instantly:

```bash
ksw -
# ✔ Switched to arn:.../eks-payments-qa
```

### History

```bash
ksw history
#   Recent contexts:
#   1  arn:.../eks-payments-dev ●
#   2  arn:.../eks-payments-qa
#   3  arn:.../eks-payments-pdn
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
  "previous": "arn:aws:eks:us-east-1:444455556666:cluster/eks-payments-qa"
}
```

## Requirements

- `kubectl` installed and configured

## License

MIT — Built by [Yonier Gómez](https://www.yonier.com)
