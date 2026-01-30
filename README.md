# folder-diff

[![CI](https://github.com/chiyiangel/folder-diff-v2/actions/workflows/ci.yml/badge.svg)](https://github.com/chiyiangel/folder-diff-v2/actions/workflows/ci.yml)
[![Release](https://github.com/chiyiangel/folder-diff-v2/actions/workflows/release.yml/badge.svg)](https://github.com/chiyiangel/folder-diff-v2/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/chiyiangel/folder-diff-v2)](https://goreportcard.com/report/github.com/chiyiangel/folder-diff-v2)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Latest Release](https://img.shields.io/github/v/release/chiyiangel/folder-diff-v2)](https://github.com/chiyiangel/folder-diff-v2/releases/latest)

A Terminal User Interface (TUI) tool for comparing file differences between two folders with **synchronized dual-pane navigation**. Built with Go and [tview](https://github.com/rivo/tview).

## Features

- **Synchronized Navigation**: Both panels move together - always showing the same path for easy comparison
- **Interactive TUI**: Dual-pane layout showing source and target directories side by side
- **Tree View**: Hierarchical display of folder structures with expand/collapse support
- **Smart Placeholders**: Shows `[Not exists]` for files that only exist in one directory
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
| `↑` / `k` | Move up (both panels) |
| `↓` / `j` | Move down (both panels) |
| `Space` / `Enter` | Expand/collapse folder |
| `d` | Jump to next difference |
| `h` / `?` | Show help |
| `q` / `Esc` | Quit application |

## Screenshot

```
┌─────────────────────────────────────────────────────────────────┐
│                📁 Folder Diff - Synchronized View                │
├────────────────────────────┬────────────────────────────────────┤
│ Source: /path/to/source    │ Target: /path/to/target            │
├────────────────────────────┼────────────────────────────────────┤
│ > 📄 file1.txt          ✓  │ > 📄 file1.txt                  ✓  │
│   📄 file2.txt          ~  │   📄 file2.txt                  ~  │
│   📄 file3.txt          -  │   [Not exists]                  -  │
│   [Not exists]          +  │   📄 file4.txt                  +  │
│   📁 subdir/            ✓  │   📁 subdir/                    ✓  │
│     📄 sub1.txt         ✓  │     📄 sub1.txt                 ✓  │
│     📄 sub2.txt         -  │     [Not exists]                -  │
├────────────────────────────┴────────────────────────────────────┤
│ ↑↓ Navigate  Space Expand  d Next Diff  h Help  q Quit          │
└─────────────────────────────────────────────────────────────────┘
```

**Note**: The `>` indicates current selection. Both panels are synchronized - 
navigation affects both sides simultaneously, making it easy to spot differences.

## How It Works

### Synchronized Navigation
Unlike traditional file comparison tools with independent panel navigation, `folder-diff` 
synchronizes both panels:

1. **Unified Path**: Both panels always show the same relative path
2. **Smart Placeholders**: Missing files show as `[Not exists]` instead of gaps
3. **Consistent View**: Expand/collapse affects both sides together
4. **Easy Comparison**: Files at the same position can be directly compared

### Comparison Logic

- **Hash Mode** (default): Calculates SHA256 hash for each file to detect content changes
- **Filename Mode**: Only compares filenames and paths (faster for large directories)

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
│       ├── layout.go     # Synchronized UI layout
│       └── sync.go       # Synchronized tree building
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Contributing

Issues and pull requests are welcome. Please ensure code style consistency and pass all tests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
