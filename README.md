# ksw

Interactive Kubernetes context switcher for your terminal. Built in Go.

Navigate with arrow keys, fuzzy search by typing, press Enter to switch. Single binary, no runtime dependencies.

> Named `ksw` to avoid conflict with macOS built-in `kswitch` (Kerberos).

## Install

### Homebrew

```bash
brew tap YonierGomez/kswitch
brew install kswitch
```

This installs the `ksw` binary.

### From source

```bash
go install github.com/YonierGomez/kswitch@latest
```

### Manual

Download the binary from [Releases](https://github.com/YonierGomez/kswitch/releases) and place it in your PATH.

## Usage

```bash
ksw                  # Interactive selector (fuzzy search)
ksw <name>           # Switch directly
ksw @<alias>         # Switch using alias
ksw alias dev <ctx>  # Create alias
ksw alias ls         # List aliases
ksw alias rm dev     # Remove alias
ksw -l               # List contexts
ksw -v               # Version
ksw -h               # Help
```

### Navigation

| Key          | Action                        |
|--------------|-------------------------------|
| Type         | Fuzzy filter in real time     |
| `↑` / `↓`   | Move up / down                |
| `Home`/`End` | Go to top / bottom           |
| `PgUp/PgDn`  | Jump 10 items                |
| `Backspace`  | Delete filter character       |
| `Enter`      | Switch to highlighted context |
| `Esc`        | Clear filter / Quit           |
| `Ctrl+C`     | Quit                          |

## Requirements

- `kubectl` installed and configured

## License

MIT
