package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andiq123/sharpify/internal/scanner"
	"github.com/andiq123/sharpify/internal/transformer"
)


type Config struct {
	Path      string
	DryRun    bool
	Rules     []string
	Verbose   bool
	Recursive bool
}


func Run(cfg Config) error {
	
	path, err := filepath.Abs(cfg.Path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path not found: %w", err)
	}

	
	s := scanner.New()

	
	var files []scanner.FileInfo
	if info.IsDir() {
		files, err = s.Scan(path)
	} else {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		files = []scanner.FileInfo{{Path: path, Content: string(content)}}
	}
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No C# files found")
		return nil
	}

	fmt.Printf("Found %d C# file(s)\n", len(files))

	
	registry := transformer.NewRegistry()

	
	var rules = registry.All()
	if len(cfg.Rules) > 0 {
		rules = registry.GetByNames(cfg.Rules)
	}

	if cfg.Verbose {
		fmt.Printf("Using %d rule(s):\n", len(rules))
		for _, r := range rules {
			fmt.Printf("  - %s: %s\n", r.Name(), r.Description())
		}
		fmt.Println()
	}

	
	t := transformer.New(rules)

	
	results := t.TransformAll(files)

	
	changedCount := 0
	for _, result := range results {
		if !result.Changed {
			continue
		}

		changedCount++
		relPath, _ := filepath.Rel(path, result.File.Path)
		if relPath == "" || strings.HasPrefix(relPath, "..") {
			relPath = result.File.Path
		}

		fmt.Printf("\n%s:\n", relPath)
		for _, rule := range result.AppliedRules {
			fmt.Printf("  ✓ %s\n", rule.Description)
		}

		if !cfg.DryRun {
			err := os.WriteFile(result.File.Path, []byte(result.NewContent), 0644)
			if err != nil {
				fmt.Printf("  ✗ Failed to write: %v\n", err)
				continue
			}
			fmt.Printf("  → File updated\n")
		}
	}

	fmt.Printf("\n%d file(s) %s\n", changedCount, modeText(cfg.DryRun))

	return nil
}

func modeText(dryRun bool) string {
	if dryRun {
		return "would be modified (dry-run)"
	}
	return "modified"
}


func ListRules() {
	registry := transformer.NewRegistry()
	rules := registry.All()

	fmt.Println("Available transformation rules:")
	fmt.Println()
	for _, r := range rules {
		fmt.Printf("  %s\n    %s\n\n", r.Name(), r.Description())
	}
}
