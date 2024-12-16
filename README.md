# folder-diff

[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A command-line tool developed in Go that compares file differences between two folders, providing detailed visual comparison through an interactive HTML report.

## Features

🔍 **Flexible Comparison**
- Hash-based comparison (SHA256) to detect content changes
- Filename-based comparison for quick directory structure analysis
- Recursive directory scanning
- File exclusion patterns support

📊 **Interactive HTML Report**
- Visual file tree representation
- Color-coded file status indicators
- Expandable/collapsible directories
- Detailed file information on click
- Status indicators:
  - 🟢 Identical files
  - 🔴 Modified files
  - 🔵 New files
  - ⚫ Deleted files

## Installation

### Prerequisites

- Go 1.22 or newer

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/folder-diff.git
cd folder-diff

# Build the project
make build

# For all platforms (Linux, Windows, macOS, Synology)
make build-all
```

## Usage

### Basic Command

```bash
folder-diff [options] <source_dir> <target_dir>
```

### Options

```
--mode      Comparison mode: hash or filename (default: hash)
--exclude   Comma-separated list of patterns to exclude
--verbose   Enable verbose output
```

### Examples

```bash
# Compare two directories using hash comparison
folder-diff /path/to/source /path/to/target

# Use filename-based comparison
folder-diff --mode=filename /path/to/source /path/to/target

# Exclude specific files
folder-diff --exclude=*.tmp,*.log /path/to/source /path/to/target

# Enable verbose output
folder-diff --verbose /path/to/source /path/to/target
```

## Output

The tool generates an `diff_report.html` file in the current directory. Open it in a web browser to:
- View the complete directory structure comparison
- Interact with the file tree
- See detailed file information
- Identify differences through color coding

## Build Targets

Available build targets in `make build-all`:
- Linux (amd64)
- Windows (amd64)
- macOS (amd64)
- Synology NAS (arm, x86)

## Development

### Project Structure

```
folder-diff/
├── cmd/folder-diff/    # Main application entry point
├── internal/           # Internal packages
│   ├── compare/       # Comparison logic
│   ├── report/        # HTML report generation
│   └── scanner/       # Directory scanning
├── docs/              # Documentation
└── build/             # Build artifacts
```

### Make Commands

```bash
make build      # Build for current platform
make build-all  # Build for all platforms
make test       # Run tests
make clean      # Clean build artifacts
make fmt        # Format code
make run        # Build and run
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Font Awesome](https://fontawesome.com/) for the icons used in the HTML report
