package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andiq123/sharpify/internal/backup"
	"github.com/andiq123/sharpify/internal/rules"
	"github.com/andiq123/sharpify/internal/scanner"
	"github.com/andiq123/sharpify/internal/transformer"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// InteractiveMode runs the interactive CLI
type InteractiveMode struct {
	registry      *transformer.RuleRegistry
	scanner       *scanner.CSharpScanner
	targetVersion rules.CSharpVersion
	safeOnly      bool
	backupMgr     *backup.Manager
}

// NewInteractive creates a new interactive mode
func NewInteractive() *InteractiveMode {
	return &InteractiveMode{
		registry:      transformer.NewRegistry(),
		scanner:       scanner.New(),
		targetVersion: rules.CSharp12,
		safeOnly:      true,
	}
}

// Run starts the interactive CLI
func (im *InteractiveMode) Run() error {
	fmt.Println(Banner())
	fmt.Println(SubtitleStyle.Render("  Modernize your legacy C# code with ease\n"))

	for {
		action, err := im.showMainMenu()
		if err != nil {
			return err
		}

		switch action {
		case "scan":
			if err := im.scanAndImprove(); err != nil {
				fmt.Println(ErrorStyle.Render("Error: " + err.Error()))
			}
		case "quick":
			if err := im.quickScan(); err != nil {
				fmt.Println(ErrorStyle.Render("Error: " + err.Error()))
			}
		case "version":
			if err := im.selectTargetVersion(); err != nil {
				fmt.Println(ErrorStyle.Render("Error: " + err.Error()))
			}
		case "rules":
			im.showRulesInfo()
		case "settings":
			if err := im.showSettings(); err != nil {
				fmt.Println(ErrorStyle.Render("Error: " + err.Error()))
			}
		case "quit":
			fmt.Println(SuccessStyle.Render("\n👋 Goodbye! Happy coding!\n"))
			return nil
		}
	}
}

func (im *InteractiveMode) showMainMenu() (string, error) {
	var action string

	versionInfo := fmt.Sprintf("Target: %s (%s)", im.targetVersion.String(), im.targetVersion.DotNetVersion())
	safeInfo := "Safe transformations only"
	if !im.safeOnly {
		safeInfo = "All transformations (including experimental)"
	}

	fmt.Println(SubtitleStyle.Render(versionInfo))
	fmt.Println(SubtitleStyle.Render(safeInfo + "\n"))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do?").
				Options(
					huh.NewOption("🔍 Scan & Improve C# Project", "scan"),
					huh.NewOption("⚡ Quick Scan (current directory)", "quick"),
					huh.NewOption("🎯 Select Target C#/.NET Version", "version"),
					huh.NewOption("📋 View Available Rules", "rules"),
					huh.NewOption("⚙️  Settings", "settings"),
					huh.NewOption("🚪 Exit", "quit"),
				).
				Value(&action),
		),
	)

	err := form.Run()
	return action, err
}

func (im *InteractiveMode) selectTargetVersion() error {
	var version string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select target C# version").
				Description("Rules for this version and earlier will be applied").
				Options(
					huh.NewOption("C# 6.0  (.NET Framework 4.6+ / .NET Core 1.0+)", "6"),
					huh.NewOption("C# 7.x  (.NET Framework 4.7+ / .NET Core 2.0+)", "7"),
					huh.NewOption("C# 8.0  (.NET Core 3.0+ / .NET Standard 2.1)", "8"),
					huh.NewOption("C# 9.0  (.NET 5.0+)", "9"),
					huh.NewOption("C# 10.0 (.NET 6.0+) - LTS", "10"),
					huh.NewOption("C# 11.0 (.NET 7.0+)", "11"),
					huh.NewOption("C# 12.0 (.NET 8.0+) - LTS [Recommended]", "12"),
					huh.NewOption("C# 13.0 (.NET 9.0+) - Latest", "13"),
				).
				Value(&version),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	switch version {
	case "6":
		im.targetVersion = rules.CSharp6
	case "7":
		im.targetVersion = rules.CSharp7
	case "8":
		im.targetVersion = rules.CSharp8
	case "9":
		im.targetVersion = rules.CSharp9
	case "10":
		im.targetVersion = rules.CSharp10
	case "11":
		im.targetVersion = rules.CSharp11
	case "12":
		im.targetVersion = rules.CSharp12
	case "13":
		im.targetVersion = rules.CSharp13
	}

	// Show available rules for this version
	availableRules := im.registry.GetByVersion(im.targetVersion, im.safeOnly)
	fmt.Println(SuccessStyle.Render(fmt.Sprintf("\n✓ Target set to %s", im.targetVersion.String())))
	fmt.Println(SubtitleStyle.Render(fmt.Sprintf("  %d rules available for this version\n", len(availableRules))))

	return nil
}

func (im *InteractiveMode) showSettings() error {
	var safeOnly bool = im.safeOnly
	var createBackup bool = true

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Safe transformations only?").
				Description("When enabled, only applies transformations guaranteed not to change logic").
				Value(&safeOnly),
			huh.NewConfirm().
				Title("Create backups before changes?").
				Description("Backup original files before applying transformations").
				Value(&createBackup),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	im.safeOnly = safeOnly

	if im.safeOnly {
		fmt.Println(SuccessStyle.Render("\n✓ Safe mode enabled - only safe transformations will be applied"))
	} else {
		fmt.Println(WarningStyle.Render("\n⚠ Experimental mode enabled - review changes carefully"))
	}

	return nil
}

func (im *InteractiveMode) scanAndImprove() error {
	// Get path
	var path string
	cwd, _ := os.Getwd()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter path to C# project or file").
				Placeholder(cwd).
				Value(&path),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if path == "" {
		path = cwd
	}

	return im.processPath(path)
}

func (im *InteractiveMode) quickScan() error {
	cwd, _ := os.Getwd()
	return im.processPath(cwd)
}

func (im *InteractiveMode) processPath(path string) error {
	// Resolve path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Check if exists
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("path not found: %w", err)
	}

	fmt.Println(SubtitleStyle.Render("\n📂 Scanning: " + absPath))
	fmt.Println(SubtitleStyle.Render(fmt.Sprintf("   Target: %s | Mode: %s\n",
		im.targetVersion.String(),
		modeString(im.safeOnly))))

	// Initialize backup manager
	im.backupMgr = backup.New(absPath)

	// Scan for files
	var files []scanner.FileInfo
	if info.IsDir() {
		files, err = im.scanner.Scan(absPath)
	} else {
		content, err := os.ReadFile(absPath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		files = []scanner.FileInfo{{Path: absPath, Content: string(content)}}
	}
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	if len(files) == 0 {
		fmt.Println(WarningStyle.Render("No C# files found in this location."))
		return nil
	}

	fmt.Println(SuccessStyle.Render(fmt.Sprintf("Found %d C# file(s)\n", len(files))))

	// Get rules for target version
	availableRules := im.registry.GetByVersion(im.targetVersion, im.safeOnly)

	// Let user select/confirm rules
	selectedRules, err := im.selectRulesFromList(availableRules)
	if err != nil {
		return err
	}

	if len(selectedRules) == 0 {
		fmt.Println(WarningStyle.Render("No rules selected. Aborting."))
		return nil
	}

	// Transform files
	t := transformer.New(selectedRules)
	results := t.TransformAll(files)

	// Filter changed files
	var changedResults []transformer.Result
	for _, r := range results {
		if r.Changed {
			changedResults = append(changedResults, r)
		}
	}

	if len(changedResults) == 0 {
		fmt.Println(SuccessStyle.Render("\n✨ All files are already up to date!\n"))
		return nil
	}

	// Show preview and confirm
	return im.reviewAndApply(changedResults, absPath)
}

func modeString(safeOnly bool) string {
	if safeOnly {
		return "Safe"
	}
	return "All"
}

func (im *InteractiveMode) selectRulesFromList(availableRules []rules.Rule) ([]rules.Rule, error) {
	var options []huh.Option[string]
	var selected []string

	// Group by version for display
	for _, r := range availableRules {
		vr := r.(rules.VersionedRule)
		safeLabel := ""
		if !vr.IsSafe() {
			safeLabel = " ⚠️"
		}
		label := fmt.Sprintf("[%s] %s - %s%s",
			vr.MinVersion().String(),
			r.Name(),
			r.Description(),
			safeLabel)
		options = append(options, huh.NewOption(label, r.Name()))
		selected = append(selected, r.Name()) // Pre-select all
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(fmt.Sprintf("Select rules to apply (target: %s)", im.targetVersion.String())).
				Description("Use space to toggle, enter to confirm").
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}

	return im.registry.GetByNames(selected), nil
}

func (im *InteractiveMode) reviewAndApply(results []transformer.Result, basePath string) error {
	fmt.Println(TitleStyle.Render(fmt.Sprintf("\n📝 %d file(s) will be modified:\n", len(results))))

	for _, result := range results {
		relPath, _ := filepath.Rel(basePath, result.File.Path)
		if relPath == "" || strings.HasPrefix(relPath, "..") {
			relPath = result.File.Path
		}

		fmt.Println(FileStyle.Render("  " + relPath))
		for _, rule := range result.AppliedRules {
			fmt.Println(RuleStyle.Render("    ✓ " + rule.Description))
		}
		fmt.Println()
	}

	// Ask what to do
	var action string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do?").
				Options(
					huh.NewOption("✅ Apply all changes (with backup)", "apply"),
					huh.NewOption("👁️  Preview changes file by file", "preview"),
					huh.NewOption("❌ Cancel", "cancel"),
				).
				Value(&action),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	switch action {
	case "apply":
		return im.applyAll(results)
	case "preview":
		return im.previewAndApply(results, basePath)
	case "cancel":
		fmt.Println(WarningStyle.Render("\nOperation cancelled."))
		return nil
	}

	return nil
}

func (im *InteractiveMode) applyAll(results []transformer.Result) error {
	// Create backups
	if im.backupMgr != nil && im.backupMgr.IsEnabled() {
		fmt.Println(SubtitleStyle.Render("\n📦 Creating backups..."))
		for _, result := range results {
			if err := im.backupMgr.Backup(result.File.Path, result.File.Content); err != nil {
				fmt.Println(WarningStyle.Render(fmt.Sprintf("  ⚠ Backup failed for %s: %v", result.File.Path, err)))
			}
		}
		fmt.Println(SuccessStyle.Render(fmt.Sprintf("  ✓ Backups saved to: %s\n", im.backupMgr.BackupDir())))
	}

	// Apply changes
	for _, result := range results {
		if err := os.WriteFile(result.File.Path, []byte(result.NewContent), 0644); err != nil {
			fmt.Println(ErrorStyle.Render(fmt.Sprintf("Failed to write %s: %v", result.File.Path, err)))
			continue
		}
	}

	fmt.Println(SuccessStyle.Render(fmt.Sprintf("\n✨ Successfully updated %d file(s)!\n", len(results))))
	return nil
}

func (im *InteractiveMode) previewAndApply(results []transformer.Result, basePath string) error {
	appliedCount := 0

	for i, result := range results {
		relPath, _ := filepath.Rel(basePath, result.File.Path)
		if relPath == "" || strings.HasPrefix(relPath, "..") {
			relPath = result.File.Path
		}

		fmt.Println(TitleStyle.Render(fmt.Sprintf("\n[%d/%d] %s", i+1, len(results), relPath)))

		// Show applied rules
		fmt.Println(SubtitleStyle.Render("Applied transformations:"))
		for _, rule := range result.AppliedRules {
			fmt.Println(RuleStyle.Render("  ✓ " + rule.Description))
		}
		fmt.Println()

		// Show diff
		im.showDiff(result.File.Content, result.NewContent)

		// Ask what to do
		var action string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Apply changes to this file?").
					Options(
						huh.NewOption("✅ Apply", "apply"),
						huh.NewOption("⏭️  Skip", "skip"),
						huh.NewOption("❌ Cancel all remaining", "cancel"),
					).
					Value(&action),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		switch action {
		case "apply":
			// Backup first
			if im.backupMgr != nil && im.backupMgr.IsEnabled() {
				if err := im.backupMgr.Backup(result.File.Path, result.File.Content); err != nil {
					fmt.Println(WarningStyle.Render(fmt.Sprintf("  ⚠ Backup failed: %v", err)))
				}
			}

			if err := os.WriteFile(result.File.Path, []byte(result.NewContent), 0644); err != nil {
				fmt.Println(ErrorStyle.Render(fmt.Sprintf("Failed to write: %v", err)))
			} else {
				fmt.Println(SuccessStyle.Render("✓ Applied"))
				appliedCount++
			}
		case "skip":
			fmt.Println(WarningStyle.Render("⏭️  Skipped"))
		case "cancel":
			fmt.Println(WarningStyle.Render("\nOperation cancelled."))
			if appliedCount > 0 {
				fmt.Println(SubtitleStyle.Render(fmt.Sprintf("  %d file(s) were already applied", appliedCount)))
			}
			return nil
		}
	}

	fmt.Println(SuccessStyle.Render(fmt.Sprintf("\n✨ Review complete! %d file(s) updated.\n", appliedCount)))

	if im.backupMgr != nil && im.backupMgr.IsEnabled() && appliedCount > 0 {
		fmt.Println(SubtitleStyle.Render(fmt.Sprintf("📦 Backups saved to: %s\n", im.backupMgr.BackupDir())))
	}

	return nil
}

func (im *InteractiveMode) showDiff(old, new string) {
	oldLines := strings.Split(old, "\n")
	newLines := strings.Split(new, "\n")

	fmt.Println(BoxStyle.Render("Changes Preview"))
	fmt.Println()

	// Simple line-by-line diff
	maxLines := len(oldLines)
	if len(newLines) > maxLines {
		maxLines = len(newLines)
	}

	shown := 0
	for i := 0; i < maxLines && shown < 40; i++ {
		oldLine := ""
		newLine := ""
		if i < len(oldLines) {
			oldLine = oldLines[i]
		}
		if i < len(newLines) {
			newLine = newLines[i]
		}

		if oldLine != newLine {
			if oldLine != "" {
				fmt.Println(DiffRemoveStyle.Render(fmt.Sprintf("- %s", oldLine)))
				shown++
			}
			if newLine != "" {
				fmt.Println(DiffAddStyle.Render(fmt.Sprintf("+ %s", newLine)))
				shown++
			}
		}
	}

	if shown >= 40 {
		fmt.Println(SubtitleStyle.Render("... (diff truncated)"))
	}
	fmt.Println()
}

func (im *InteractiveMode) showRulesInfo() {
	fmt.Println(TitleStyle.Render("\n📋 Available Transformation Rules\n"))

	groups := im.registry.GroupByVersion()

	versions := []rules.CSharpVersion{
		rules.CSharp6, rules.CSharp7, rules.CSharp8,
		rules.CSharp9, rules.CSharp10, rules.CSharp11,
		rules.CSharp12, rules.CSharp13,
	}

	for _, version := range versions {
		ruleList, ok := groups[version]
		if !ok || len(ruleList) == 0 {
			continue
		}

		// Mark current target version
		marker := ""
		if version == im.targetVersion {
			marker = " ← current target"
		}

		versionStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			Background(lipgloss.Color("#1E1B4B")).
			Padding(0, 1)

		fmt.Println(versionStyle.Render(version.String() + " (" + version.DotNetVersion() + ")" + marker))

		for _, r := range ruleList {
			vr := r.(rules.VersionedRule)
			safeLabel := SuccessStyle.Render("✓ safe")
			if !vr.IsSafe() {
				safeLabel = WarningStyle.Render("⚠ review")
			}

			fmt.Printf("  %s %s\n", RuleStyle.Render("•"), r.Name())
			fmt.Printf("    %s [%s]\n", SubtitleStyle.Render(r.Description()), safeLabel)
		}
		fmt.Println()
	}

	// Wait for user
	var cont bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Press Enter to continue").
				Affirmative("OK").
				Negative("").
				Value(&cont),
		),
	)
	_ = form.Run()
}
