package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"folder-diff-v2/internal/compare"
	"folder-diff-v2/internal/scanner"
	"folder-diff-v2/internal/tui"
)

// Version information (set by ldflags during build)
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

func main() {
	mode := flag.String("mode", "hash", "Comparison mode: hash or filename")
	exclude := flag.String("exclude", "", "Comma-separated list of patterns to exclude")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	version := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Show version if requested
	if *version {
		fmt.Printf("folder-diff %s\n", Version)
		fmt.Printf("Commit:     %s\n", Commit)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Go Version: %s\n", GoVersion)
		os.Exit(0)
	}

	if flag.NArg() != 2 {
		fmt.Println("Usage: folder-diff [options] <source_dir> <target_dir>")
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  folder-diff /path/to/source /path/to/target")
		fmt.Println("  folder-diff --mode=filename /path/to/source /path/to/target")
		fmt.Println("  folder-diff --exclude=*.tmp,*.log /path/to/source /path/to/target")
		os.Exit(1)
	}

	sourceDir := flag.Arg(0)
	targetDir := flag.Arg(1)

	// Validate directories
	if err := validateDirectory(sourceDir); err != nil {
		log.Fatalf("Source directory error: %v", err)
	}
	if err := validateDirectory(targetDir); err != nil {
		log.Fatalf("Target directory error: %v", err)
	}

	var excludePatterns []string
	if *exclude != "" {
		excludePatterns = strings.Split(*exclude, ",")
	}

	if *verbose {
		fmt.Printf("folder-diff %s\n", Version)
		fmt.Printf("Scanning directories...\n")
		fmt.Printf("Source: %s\n", sourceDir)
		fmt.Printf("Target: %s\n", targetDir)
		fmt.Printf("Mode: %s\n", *mode)
		if len(excludePatterns) > 0 {
			fmt.Printf("Exclude patterns: %v\n", excludePatterns)
		}
	}

	// Scan directories
	s := scanner.NewScanner(excludePatterns)
	comparator := compare.NewComparator(compare.ComparisonMode(*mode))

	sourceFiles, err := s.ScanDirectory(sourceDir)
	if err != nil {
		log.Fatalf("Error scanning source directory: %v", err)
	}

	targetFiles, err := s.ScanDirectory(targetDir)
	if err != nil {
		log.Fatalf("Error scanning target directory: %v", err)
	}

	// Compare
	result := comparator.Compare(sourceFiles, targetFiles)
	result.SourceRoot = sourceDir
	result.TargetRoot = targetDir

	if *verbose {
		fmt.Printf("Found %d source files, %d target files\n", len(sourceFiles), len(targetFiles))
		fmt.Println("Starting TUI...")
	}

	// Start TUI
	app := tui.NewApp(result, sourceDir, targetDir)
	if err := app.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}

// validateDirectory checks if a path is a valid directory
func validateDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", path)
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", path)
	}
	return nil
}
