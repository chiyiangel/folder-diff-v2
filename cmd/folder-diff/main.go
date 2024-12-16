package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"folder-diff-v2/internal/compare"
	"folder-diff-v2/internal/report"
	"folder-diff-v2/internal/scanner"
)

func main() {
	mode := flag.String("mode", "hash", "Comparison mode: hash or filename")
	exclude := flag.String("exclude", "", "Comma-separated list of patterns to exclude")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Println("Usage: folder-diff [options] <source_dir> <target_dir>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	sourceDir := flag.Arg(0)
	targetDir := flag.Arg(1)

	var excludePatterns []string
	if *exclude != "" {
		excludePatterns = strings.Split(*exclude, ",")
	}

	scanner := scanner.NewScanner(excludePatterns)
	comparator := compare.NewComparator(compare.ComparisonMode(*mode))

	sourceFiles, err := scanner.ScanDirectory(sourceDir)
	if err != nil {
		log.Fatalf("Error scanning source directory: %v", err)
	}

	targetFiles, err := scanner.ScanDirectory(targetDir)
	if err != nil {
		log.Fatalf("Error scanning target directory: %v", err)
	}

	result := comparator.Compare(sourceFiles, targetFiles)
	result.SourceRoot = sourceDir
	result.TargetRoot = targetDir

	reporter := report.NewGenerator()
	if err := reporter.GenerateHTML(result, "diff_report.html"); err != nil {
		log.Fatalf("Error generating report: %v", err)
	}

	if *verbose {
		printVerboseResults(result)
	}
}

func printVerboseResults(result *compare.ComparisonResult) {
	fmt.Printf("Comparison completed:\n")
	fmt.Printf("Source: %s\n", result.SourceRoot)
	fmt.Printf("Target: %s\n", result.TargetRoot)
	fmt.Printf("Mode: %s\n", result.Mode)
}
