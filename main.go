package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/andiq123/sharpify/cmd"
	"github.com/andiq123/sharpify/internal/ui"
)

var version = "1.0.0"

func main() {
	
	batch := flag.Bool("batch", false, "Run in batch mode (non-interactive)")
	batchShort := flag.Bool("b", false, "Run in batch mode (non-interactive)")
	dryRun := flag.Bool("dry-run", false, "Preview changes without modifying files")
	rulesFlag := flag.String("rules", "", "Comma-separated list of rules to apply (default: all)")
	verbose := flag.Bool("verbose", false, "Show detailed output")
	listRules := flag.Bool("list-rules", false, "List all available transformation rules")
	showVersion := flag.Bool("version", false, "Show version")
	help := flag.Bool("help", false, "Show help")

	flag.Usage = printUsage

	flag.Parse()

	if *help {
		printUsage()
		return
	}

	if *showVersion {
		fmt.Printf("sharpify v%s\n", version)
		return
	}

	if *listRules {
		cmd.ListRules()
		return
	}

	
	if *batch || *batchShort {
		runBatch(dryRun, rulesFlag, verbose)
		return
	}

	
	im := ui.NewInteractive()
	if err := im.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runBatch(dryRun *bool, rulesFlag *string, verbose *bool) {
	path := "."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	var rules []string
	if *rulesFlag != "" {
		rules = strings.Split(*rulesFlag, ",")
		for i := range rules {
			rules[i] = strings.TrimSpace(rules[i])
		}
	}

	cfg := cmd.Config{
		Path:    path,
		DryRun:  *dryRun,
		Rules:   rules,
		Verbose: *verbose,
	}

	if err := cmd.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`sharpify - Modernize legacy C# code

Usage:
  sharpify [flags] [path]
  sharpify                             # Run interactive mode (default)
  sharpify -b ./src                    # Batch mode on ./src

Arguments:
  path    Path to C# file or directory (default: current directory)

Flags:
  -b, --batch      Run in batch mode (non-interactive)
  --dry-run        Preview changes without modifying files
  --rules          Comma-separated list of rules to apply (default: all)
  --verbose        Show detailed output
  --list-rules     List all available transformation rules
  --version        Show version
  --help           Show this help

Examples:
  sharpify                             # Interactive mode (default)
  sharpify -b .                        # Batch: improve all C# files
  sharpify -b --dry-run ./src          # Preview changes
  sharpify -b --rules file-scoped-namespace,pattern-matching ./MyProject

Available Rules:
  file-scoped-namespace    Convert to file-scoped namespaces (C# 10+)
  var-pattern              Use var for obvious type declarations
  target-typed-new         Use target-typed new (C# 9+)
  null-coalescing          Use ?? and ??= operators
  expression-body          Use expression-bodied members (C# 6+)
  string-interpolation     Use string interpolation (C# 6+)
  pattern-matching         Use pattern matching (C# 7+)
  collection-expression    Use collection expressions (C# 12+)
  index-range              Use ^index and ranges (C# 8+)`)
	fmt.Println()
}
