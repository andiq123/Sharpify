package ui

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

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
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewInteractive() *InteractiveMode {
	// Load config but clear the working path - always start fresh
	cfg := config.Load()
	cfg.WorkingPath = "" // Don't remember last path

	ctx, cancel := context.WithCancel(context.Background())

	im := &InteractiveMode{
		registry: transformer.NewRegistry(),
		scanner:  scanner.New(),
		config:   cfg,
		ctx:      ctx,
		cancel:   cancel,
	}

	// Setup Ctrl+C handler
	im.setupSignalHandler()

	return im
}

func (im *InteractiveMode) setupSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n" + SubtitleStyle.Render("Goodbye! 👋"))
		im.cancel()
		os.Exit(0)
	}()
}

func (im *InteractiveMode) Run() error {
	fmt.Println(Banner())
	fmt.Println(SubtitleStyle.Render("  Press Ctrl+C anytime to exit"))
	fmt.Println()
	im.printStatus()

	for {
		// Check if context was cancelled
		select {
		case <-im.ctx.Done():
			return nil
		default:
		}

		action := im.showMainMenu()

		switch action {
		case "run":
			im.runImprove()
		case "settings":
			im.showSettings()
		case "rules":
			im.showRules()
		case "quit", "":
			fmt.Println(SubtitleStyle.Render("Goodbye! 👋"))
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
				Title("What would you like to do?").
				Options(
					huh.NewOption("▶ Run  - Scan and improve C# code", "run"),
					huh.NewOption("⚙ Settings  - Configure version and options", "settings"),
					huh.NewOption("📋 Rules  - View available transformations", "rules"),
					huh.NewOption("✕ Exit  - Quit Sharpify", "quit"),
				).
				Value(&action),
		),
	)

	if err := form.Run(); err != nil {
		return "quit" // Ctrl+C or error - exit gracefully
	}
	return action
}

func (im *InteractiveMode) runImprove() {
	// Always ask for path - get current directory as default suggestion
	currentDir, _ := os.Getwd()

	fmt.Println()
	fmt.Println(TitleStyle.Render("📂 Select Project Path"))
	fmt.Println(SubtitleStyle.Render("Press Enter to use current directory, or type a new path"))

	var inputPath string
	err := huh.NewInput().
		Title("Path to C# project").
		Placeholder(currentDir + " (current)").
		Value(&inputPath).
		Run()

	if err != nil {
		// Ctrl+C pressed
		return
	}

	// Use current directory if empty input
	workingDir := currentDir
	if inputPath != "" {
		// Expand ~ to home directory
		if len(inputPath) > 0 && inputPath[0] == '~' {
			home, _ := os.UserHomeDir()
			inputPath = filepath.Join(home, inputPath[1:])
		}

		// Resolve to absolute path
		absPath, err := filepath.Abs(inputPath)
		if err != nil {
			fmt.Println(ErrorStyle.Render("Invalid path: " + err.Error()))
			return
		}

		// Check if path exists
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			fmt.Println(ErrorStyle.Render("Path does not exist: " + absPath))
			return
		}

		workingDir = absPath
	}

	fmt.Printf("\n%s %s\n", SubtitleStyle.Render("Scanning:"), workingDir)

	// Scan for C# files
	files, err := im.scanner.Scan(workingDir)
	if err != nil {
		fmt.Println(ErrorStyle.Render("Scan failed: " + err.Error()))
		return
	}

	if len(files) == 0 {
		fmt.Println(WarningStyle.Render("No C# files found in: " + workingDir))
		fmt.Println(SubtitleStyle.Render("Tip: Make sure the path contains .cs files"))
		return
	}

	fmt.Printf("\n%s Found %d C# file(s)\n", SuccessStyle.Render("✓"), len(files))

	// Feature 2: Show ALL rules preview and allow toggling before run
	// Get ALL version-compatible rules (not filtered by safety)
	allRules := im.registry.GetByVersion(im.config.GetVersion(), false)

	if len(allRules) == 0 {
		fmt.Println(WarningStyle.Render("No rules available for current settings"))
		return
	}

	// Build options for multi-select with visual indicators
	var ruleOptions []huh.Option[string]
	var selectedRules []string

	for _, r := range allRules {
		vr := r.(rules.VersionedRule)

		// Add visual indicator for safe vs unsafe
		var indicator string
		if vr.IsSafe() {
			indicator = SuccessStyle.Render("✓")
		} else {
			indicator = WarningStyle.Render("⚠")
		}

		label := fmt.Sprintf("%s %s - %s", indicator, r.Name(), r.Description())
		ruleOptions = append(ruleOptions, huh.NewOption(label, r.Name()))

		// Pre-selection: if safeOnly mode, only pre-select safe rules
		// Otherwise pre-select all rules that aren't explicitly disabled
		if !im.config.IsRuleDisabled(r.Name()) {
			if !im.config.SafeOnly || vr.IsSafe() {
				selectedRules = append(selectedRules, r.Name())
			}
		}
	}

	fmt.Println()
	fmt.Println(TitleStyle.Render("Configure Rules"))
	fmt.Println(SubtitleStyle.Render("✓ safe  ⚠ opt-in (may need review) | space to toggle, enter to confirm"))

	_ = huh.NewMultiSelect[string]().
		Title("Rules to Apply").
		Options(ruleOptions...).
		Value(&selectedRules).
		Run()

	// Feature 3: Update and persist rule enabled/disabled state
	for _, r := range allRules {
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
	enabledRules := im.config.GetEnabledRules(allRules)

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
	err = huh.NewConfirm().
		Title("Apply changes?").
		Affirmative("Yes, apply").
		Negative("No, cancel").
		Value(&confirm).
		Run()

	if err != nil {
		// Ctrl+C pressed
		fmt.Println(SubtitleStyle.Render("Cancelled"))
		return
	}

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
		// Ctrl+C - return to main menu
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

	fmt.Println(SubtitleStyle.Render("  Press Enter to return to main menu"))
	var cont bool
	_ = huh.NewConfirm().Title("").Affirmative("← Back").Negative("").Value(&cont).Run()
}
