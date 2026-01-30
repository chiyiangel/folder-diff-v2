# folder-diff

[![Go Report Card](https://goreportcard.com/badge/github.com/chiyiangel/folder-diff-v2)](https://goreportcard.com/report/github.com/chiyiangel/folder-diff-v2)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A Terminal User Interface (TUI) tool for comparing file differences between two folders. Built with Go and [tview](https://github.com/rivo/tview).

## Features

- **Interactive TUI**: Dual-pane layout showing source and target directories side by side
- **Tree View**: Hierarchical display of folder structures with expand/collapse support
- **Color-Coded Status**:
  - 🟢 Green (✓) - Identical files
  - 🔴 Red (~) - Modified files
  - 🔵 Blue (+) - New files (target only)
  - ⚫ Gray (-) - Deleted files (source only)
- **Comparison Modes**:
  - `hash`: Compare file contents using SHA256 (default)
  - `filename`: Compare by filename only (faster)
- **Pattern Exclusion**: Skip files/directories matching specified patterns
- **Keyboard Navigation**: Full keyboard support for efficient browsing

## Installation

### From Source

```bash
git clone https://github.com/chiyiangel/folder-diff-v2.git
cd folder-diff-v2
go build -o folder-diff ./cmd/folder-diff/
```

### Using Go Install

```bash
go install github.com/chiyiangel/folder-diff-v2/cmd/folder-diff@latest
```

## Usage

```bash
folder-diff [options] <source_dir> <target_dir>
```

### Options

| Option | Description |
|--------|-------------|
| `--mode=hash` | Compare by file content hash (default) |
| `--mode=filename` | Compare by filename only |
| `--exclude=PATTERNS` | Comma-separated patterns to exclude |
| `--verbose` | Show verbose output during scanning |

### Examples

```bash
# Basic comparison
folder-diff /path/to/source /path/to/target

# Filename-only comparison (faster)
folder-diff --mode=filename /path/to/source /path/to/target

# Exclude certain files
folder-diff --exclude=*.tmp,*.log,node_modules /path/to/source /path/to/target

# Verbose mode
folder-diff --verbose /path/to/source /path/to/target
```

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate up/down in tree |
| `←` / `→` | Switch between source/target panels |
| `Tab` | Toggle panel focus |
| `Space` / `Enter` | Expand/collapse folder |
| `d` | Jump to next difference |
| `h` / `?` | Show help |
| `q` / `Esc` | Quit application |

## Screenshot

```
┌─────────────────────────────────────────────────────────────────┐
│                    📁 Folder Diff - TUI Mode                    │
├────────────────────────────┬────────────────────────────────────┤
│ Source: /path/to/source    │ Target: /path/to/target            │
├────────────────────────────┼────────────────────────────────────┤
│ 📁 folder1/             ✓  │ 📁 folder1/                     ✓  │
│   📄 file1.txt          ✓  │   📄 file1.txt                  ✓  │
│   📄 file2.txt          ~  │   📄 file2.txt                  ~  │
│   📄 file3.txt          -  │                                    │
│ 📁 folder2/             +  │ 📁 folder2/                     +  │
│                            │   📄 newfile.txt                +  │
├────────────────────────────┴────────────────────────────────────┤
│ ↑↓ Navigate  ←→ Switch  Space Expand  d Next Diff  q Quit      │
└─────────────────────────────────────────────────────────────────┘
```

## Project Structure

```
folder-diff-v2/
├── cmd/folder-diff/
│   └── main.go           # Application entry point
├── internal/
│   ├── compare/
│   │   ├── types.go      # Data structures
│   │   └── comparator.go # Comparison logic
│   ├── scanner/
│   │   └── scanner.go    # Directory scanning
│   └── tui/
│       ├── app.go        # TUI application controller
│       ├── layout.go     # UI layout management
│       ├── tree.go       # Tree view component
│       └── keybinds.go   # Keyboard handling
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Contributing

Issues and pull requests are welcome. Please ensure code style consistency and pass all tests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
