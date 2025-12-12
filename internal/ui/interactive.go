package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/andiq123/sharpify/internal/backup"
	"github.com/andiq123/sharpify/internal/config"
	"github.com/andiq123/sharpify/internal/rules"
	"github.com/andiq123/sharpify/internal/scanner"
	"github.com/andiq123/sharpify/internal/transformer"
	"github.com/charmbracelet/huh"
)

type InteractiveMode struct {
	registry  *transformer.RuleRegistry
	scanner   *scanner.CSharpScanner
	config    *config.Config
	backupMgr *backup.Manager
}

func NewInteractive() *InteractiveMode {
	return &InteractiveMode{
		registry: transformer.NewRegistry(),
		scanner:  scanner.New(),
		config:   config.Load(),
	}
}

func (im *InteractiveMode) Run() error {
	fmt.Println(Banner())
	im.printStatus()

	for {
		action := im.showMainMenu()

		switch action {
		case "run":
			im.runImprove()
		case "settings":
			im.showSettings()
		case "rules":
			im.showRules()
		case "quit":
			return nil
		}
	}
}

func (im *InteractiveMode) printStatus() {
	version := im.config.GetVersion()
	mode := "safe"
	if !im.config.SafeOnly {
		mode = "all"
	}
	bkp := "off"
	if im.config.BackupEnabled {
		bkp = "on"
	}

	fmt.Printf("  %s | %s mode | backup %s\n\n",
		SubtitleStyle.Render(version.String()),
		mode,
		bkp)
}

func (im *InteractiveMode) showMainMenu() string {
	var action string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Sharpify").
				Options(
					huh.NewOption("Run", "run"),
					huh.NewOption("Settings", "settings"),
					huh.NewOption("Rules", "rules"),
					huh.NewOption("Exit", "quit"),
				).
				Value(&action),
		),
	)

	_ = form.Run()
	return action
}

func (im *InteractiveMode) runImprove() {
	// Determine working directory - use saved path or current directory
	workingDir := im.config.WorkingPath
	if workingDir == "" {
		workingDir, _ = os.Getwd()
	}

	// Scan for C# files
	files, err := im.scanner.Scan(workingDir)
	if err != nil {
		fmt.Println(ErrorStyle.Render("Scan failed: " + err.Error()))
		return
	}

	// Feature 1: If no files found, prompt for alternative path
	if len(files) == 0 {
		fmt.Println(WarningStyle.Render("No C# files found in: " + workingDir))

		var selectPath bool
		_ = huh.NewConfirm().
			Title("Select a different path?").
			Value(&selectPath).
			Run()

		if !selectPath {
			return
		}

		var newPath string
		_ = huh.NewInput().
			Title("Enter path to C# project").
			Placeholder(workingDir).
			Value(&newPath).
			Run()

		if newPath == "" {
			fmt.Println("Cancelled")
			return
		}

		// Expand ~ to home directory
		if len(newPath) > 0 && newPath[0] == '~' {
			home, _ := os.UserHomeDir()
			newPath = filepath.Join(home, newPath[1:])
		}

		// Resolve to absolute path
		absPath, err := filepath.Abs(newPath)
		if err != nil {
			fmt.Println(ErrorStyle.Render("Invalid path: " + err.Error()))
			return
		}

		// Check if path exists
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			fmt.Println(ErrorStyle.Render("Path does not exist: " + absPath))
			return
		}

		// Scan the new path
		files, err = im.scanner.Scan(absPath)
		if err != nil {
			fmt.Println(ErrorStyle.Render("Scan failed: " + err.Error()))
			return
		}

		if len(files) == 0 {
			fmt.Println(WarningStyle.Render("No C# files found in: " + absPath))
			return
		}

		// Save the working path for next time
		workingDir = absPath
		im.config.WorkingPath = absPath
		_ = im.config.Save()
		fmt.Println(SuccessStyle.Render("Path saved for future sessions"))
	}

	fmt.Printf("\nFound %d file(s) in %s\n", len(files), workingDir)

	// Feature 2: Show rules preview and allow toggling before run
	availableRules := im.registry.GetByVersion(im.config.GetVersion(), im.config.SafeOnly)

	if len(availableRules) == 0 {
		fmt.Println(WarningStyle.Render("No rules available for current settings"))
		return
	}

	// Build options for multi-select with current enabled/disabled state
	var ruleOptions []huh.Option[string]
	var selectedRules []string

	for _, r := range availableRules {
		ruleOptions = append(ruleOptions, huh.NewOption(r.Name()+" - "+r.Description(), r.Name()))
		// Pre-select rules that are NOT disabled
		if !im.config.IsRuleDisabled(r.Name()) {
			selectedRules = append(selectedRules, r.Name())
		}
	}

	fmt.Println()
	fmt.Println(TitleStyle.Render("Configure Rules"))
	fmt.Println(SubtitleStyle.Render("Select which rules to apply (space to toggle, enter to confirm)"))

	_ = huh.NewMultiSelect[string]().
		Title("Enabled Rules").
		Options(ruleOptions...).
		Value(&selectedRules).
		Run()

	// Feature 3: Update and persist rule enabled/disabled state
	for _, r := range availableRules {
		isSelected := false
		for _, sel := range selectedRules {
			if sel == r.Name() {
				isSelected = true
				break
			}
		}
		// SetRuleDisabled(name, disabled) - if selected, it's NOT disabled
		im.config.SetRuleDisabled(r.Name(), !isSelected)
	}
	_ = im.config.Save()

	// Filter to only enabled rules
	enabledRules := im.config.GetEnabledRules(availableRules)

	if len(enabledRules) == 0 {
		fmt.Println(WarningStyle.Render("No rules enabled - nothing to do"))
		return
	}

	fmt.Printf("\nApplying %d rule(s)...\n", len(enabledRules))

	t := transformer.New(enabledRules)
	results := t.TransformAll(files)

	var changed []transformer.Result
	for _, r := range results {
		if r.Changed {
			changed = append(changed, r)
		}
	}

	if len(changed) == 0 {
		fmt.Println(SuccessStyle.Render("All files up to date"))
		return
	}

	fmt.Printf("\n%d file(s) to update:\n", len(changed))
	for _, r := range changed {
		rel, _ := filepath.Rel(workingDir, r.File.Path)
		fmt.Printf("  %s\n", rel)
		for _, rule := range r.AppliedRules {
			fmt.Printf("    %s\n", RuleStyle.Render("+ "+rule.Description))
		}
	}

	var confirm bool
	_ = huh.NewConfirm().
		Title("Apply changes?").
		Value(&confirm).
		Run()

	if !confirm {
		fmt.Println("Cancelled")
		return
	}

	if im.config.BackupEnabled {
		im.backupMgr = backup.New(workingDir)
		for _, r := range changed {
			_ = im.backupMgr.Backup(r.File.Path, r.File.Content)
		}
		fmt.Printf("Backup: %s\n", im.backupMgr.BackupDir())
	}

	for _, r := range changed {
		_ = os.WriteFile(r.File.Path, []byte(r.NewContent), 0644)
	}

	fmt.Println(SuccessStyle.Render(fmt.Sprintf("Updated %d file(s)", len(changed))))
}

func (im *InteractiveMode) showSettings() {
	var version string
	var safeOnly bool = im.config.SafeOnly
	var backupEnabled bool = im.config.BackupEnabled

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("C# Version").
				Options(
					huh.NewOption("C# 6  (.NET 4.6+)", "6"),
					huh.NewOption("C# 7  (.NET Core 2.0+)", "7"),
					huh.NewOption("C# 8  (.NET Core 3.0+)", "8"),
					huh.NewOption("C# 9  (.NET 5.0+)", "9"),
					huh.NewOption("C# 10 (.NET 6.0+)", "10"),
					huh.NewOption("C# 11 (.NET 7.0+)", "11"),
					huh.NewOption("C# 12 (.NET 8.0+)", "12"),
					huh.NewOption("C# 13 (.NET 9.0+)", "13"),
				).
				Value(&version),
			huh.NewConfirm().
				Title("Safe mode").
				Description("Only apply safe transformations").
				Value(&safeOnly),
			huh.NewConfirm().
				Title("Backup").
				Description("Create backups before changes").
				Value(&backupEnabled),
		),
	)

	if err := form.Run(); err != nil {
		return
	}

	im.config.TargetVersion = version
	im.config.SafeOnly = safeOnly
	im.config.BackupEnabled = backupEnabled
	_ = im.config.Save()

	fmt.Println(SuccessStyle.Render("Settings saved"))
	im.printStatus()
}

func (im *InteractiveMode) showRules() {
	groups := im.registry.GroupByVersion()
	versions := []rules.CSharpVersion{
		rules.CSharp6, rules.CSharp7, rules.CSharp8,
		rules.CSharp9, rules.CSharp10, rules.CSharp11,
		rules.CSharp12, rules.CSharp13,
	}

	fmt.Println()
	for _, v := range versions {
		ruleList, ok := groups[v]
		if !ok || len(ruleList) == 0 {
			continue
		}

		marker := ""
		if v == im.config.GetVersion() {
			marker = " *"
		}

		fmt.Println(TitleStyle.Render(v.String() + marker))
		for _, r := range ruleList {
			vr := r.(rules.VersionedRule)
			safe := SuccessStyle.Render("safe")
			if !vr.IsSafe() {
				safe = WarningStyle.Render("opt-in")
			}
			fmt.Printf("  %-25s %s [%s]\n", r.Name(), SubtitleStyle.Render(r.Description()), safe)
		}
		fmt.Println()
	}

	var cont bool
	_ = huh.NewConfirm().Title("").Affirmative("OK").Negative("").Value(&cont).Run()
}
