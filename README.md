# folder-diff

[![Go Report Card](https://goreportcard.com/badge/github.com/chiyiangel/folder-diff)](https://goreportcard.com/report/github.com/chiyiangel/folder-diff)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

folder-diff is a command-line tool written in Go for comparing all file differences between two folders. It supports both hash-based and filename-based comparison, and generates a visual HTML report.

## Features

- **File Comparison**: Recursively compare all files in two folders, including subdirectories
- **Multiple Comparison Modes**:
  - Hash-based: Use cryptographic hash functions (e.g., MD5, SHA256) to compare file contents
  - Filename-based: Compare files based on names only, faster execution
- **HTML Report**:
  - Visualize folder structure using file trees
  - Color-coded differences:
    - Green: Identical files
    - Red: Modified files
    - Blue: New files
    - Gray: Deleted files
  - Interactive features: Expand/collapse folders, click files for detailed information
- **Command-line Interface**:
  - Multiple options and flags
  - Exclude specific files or directories
  - Verbose output mode

## Installation

### Using Go

```bash
go install github.com/chiyiangel/folder-diff@latest
```

### From Source

```bash
git clone https://github.com/chiyiangel/folder-diff.git
cd folder-diff
make build
```

## Usage

Basic usage:

```bash
folder-diff /path/to/folder1 /path/to/folder2
```

Options:

- `--mode=hash`: Use hash-based comparison (default)
- `--mode=filename`: Use filename-based comparison
- `--exclude`: Exclude specific files or directories
- `--verbose`: Show verbose output

Example:

```bash
folder-diff /path/to/folder1 /path/to/folder2 --mode=hash --exclude=*.tmp
```

## Contributing

We welcome issues and pull requests. Please ensure consistent code style and pass all tests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
