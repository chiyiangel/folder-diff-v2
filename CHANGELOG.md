# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Terminal User Interface (TUI) with synchronized dual-pane navigation
- Cross-platform build support (Linux, macOS, Windows)
- GitHub Actions CI/CD workflows
- Automatic release packaging
- Automatic CHANGELOG generation
- Comprehensive Makefile with multiple targets

### Changed
- Replaced HTML report output with interactive TUI
- Improved file comparison visualization with color coding

### Fixed
- TUI color tag syntax for proper rendering
- Icon differentiation for non-existent files

## [v1.0.0] - 2024-01-30

### Added
- Initial release with TUI support
- File comparison using hash or filename mode
- Pattern exclusion support
- Keyboard navigation
- Synchronized navigation between source and target directories
- Smart placeholders for missing files
