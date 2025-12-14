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
	cfg := config.Load()
	cfg.WorkingPath = ""

	ctx, cancel := context.WithCancel(context.Background())

	im := &InteractiveMode{
		registry: transformer.NewRegistry(),
		scanner:  scanner.New(),
		config:   cfg,
		ctx:      ctx,
		cancel:   cancel,
	}

	im.setupSignalHandler()
	return im
}

func (im *InteractiveMode) setupSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\n" + SubtitleStyle.Render("👋 Goodbye!"))
		im.cancel()
		os.Exit(0)
	}()
}

func (im *InteractiveMode) Run() error {
	fmt.Println(Banner())
	fmt.Println()

	for {
		select {
		case <-im.ctx.Done():
			return nil
		default:
		}

		im.printStatusBar()
		action := im.showMainMenu()

		switch action {
		case "quick":
			im.runQuick()
		case "custom":
			im.runCustom()
		case "settings":
			im.showSettings()
		case "rules":
			im.showRules()
		case "quit", "":
			fmt.Println("\n" + SubtitleStyle.Render("👋 Goodbye!"))
			return nil
		}
	}
}

func (im *InteractiveMode) printStatusBar() {
	version := im.config.GetVersion()
	ruleCount := len(im.registry.GetByVersion(version, im.config.SafeOnly))

	mode := SuccessStyle.Render("safe")
	if !im.config.SafeOnly {
		mode = WarningStyle.Render("all rules")
	}

	fmt.Printf("  %s  •  %s  •  %s\n\n",
		AccentStyle.Render(version.String()),
		mode,
		SubtitleStyle.Render(fmt.Sprintf("%d rules available", ruleCount)),
	)
}

func (im *InteractiveMode) showMainMenu() string {
	var action string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do?").
				Options(
					huh.NewOption("⚡ Quick Run  →  Improve code with smart defaults", "quick"),
					huh.NewOption("🎯 Custom Run  →  Choose specific rules to apply", "custom"),
					huh.NewOption("⚙  Settings  →  Configure C# version & options", "settings"),
					huh.NewOption("📋 Rules  →  Browse available transformations", "rules"),
					huh.NewOption("✕  Exit", "quit"),
				).
				Value(&action),
		),
	)

	if err := form.Run(); err != nil {
		return "quit"
	}
	return action
}

func (im *InteractiveMode) runQuick() {
	path := im.selectPath()
	if path == "" {
		return
	}

	files, err := im.scanFiles(path)
	if err != nil || len(files) == 0 {
		return
	}

	// Use all enabled safe rules
	allRules := im.registry.GetByVersion(im.config.GetVersion(), im.config.SafeOnly)
	enabledRules := im.config.GetEnabledRules(allRules)

	if len(enabledRules) == 0 {
		fmt.Println(Warn("No rules enabled. Go to Settings to configure."))
		return
	}

	fmt.Printf("\n  %s Applying %d rules...\n", InfoStyle.Render("⚡"), len(enabledRules))

	im.applyTransformations(path, files, enabledRules)
}

func (im *InteractiveMode) runCustom() {
	path := im.selectPath()
	if path == "" {
		return
	}

	files, err := im.scanFiles(path)
	if err != nil || len(files) == 0 {
		return
	}

	// Show rule selector
	allRules := im.registry.GetByVersion(im.config.GetVersion(), false)
	selectedRules := im.selectRules(allRules)

	if len(selectedRules) == 0 {
		fmt.Println(Warn("No rules selected."))
		return
	}

	im.applyTransformations(path, files, selectedRules)
}

func (im *InteractiveMode) selectPath() string {
	currentDir, _ := os.Getwd()

	fmt.Println()
	fmt.Println(TitleStyle.Render("📂 Project Path"))
	fmt.Println(Tip("Press Enter for current directory, or type a path"))
	fmt.Println()

	var inputPath string
	err := huh.NewInput().
		Title("Path").
		Placeholder(currentDir).
		Value(&inputPath).
		Run()

	if err != nil {
		return ""
	}

	workingDir := currentDir
	if inputPath != "" {
		if len(inputPath) > 0 && inputPath[0] == '~' {
			home, _ := os.UserHomeDir()
			inputPath = filepath.Join(home, inputPath[1:])
		}

		absPath, err := filepath.Abs(inputPath)
		if err != nil {
			fmt.Println(Fail("Invalid path: " + err.Error()))
			return ""
		}

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			fmt.Println(Fail("Path does not exist"))
			return ""
		}

		workingDir = absPath
	}

	return workingDir
}

func (im *InteractiveMode) scanFiles(path string) ([]scanner.FileInfo, error) {
	fmt.Printf("\n  %s Scanning for C# files...\n", InfoStyle.Render("🔍"))

	files, err := im.scanner.Scan(path)
	if err != nil {
		fmt.Println(Fail("Scan failed: " + err.Error()))
		return nil, err
	}

	if len(files) == 0 {
		fmt.Println(Warn("No C# files found"))
		fmt.Println(Tip("Make sure the path contains .cs files"))
		return nil, nil
	}

	fmt.Println(Success(fmt.Sprintf("Found %d C# file(s)", len(files))))
	return files, nil
}

func (im *InteractiveMode) selectRules(allRules []rules.Rule) []rules.Rule {
	fmt.Println()
	fmt.Println(TitleStyle.Render("🎯 Select Rules"))
	fmt.Println(Tip("Space to toggle, Enter to confirm"))
	fmt.Println()

	var ruleOptions []huh.Option[string]
	var selectedNames []string

	for _, r := range allRules {
		vr := r.(rules.VersionedRule)

		icon := SuccessStyle.Render("✓")
		if !vr.IsSafe() {
			icon = WarningStyle.Render("⚠")
		}

		label := fmt.Sprintf("%s %s", icon, r.Name())
		ruleOptions = append(ruleOptions, huh.NewOption(label, r.Name()))

		if !im.config.IsRuleDisabled(r.Name()) {
			if !im.config.SafeOnly || vr.IsSafe() {
				selectedNames = append(selectedNames, r.Name())
			}
		}
	}

	_ = huh.NewMultiSelect[string]().
		Title(fmt.Sprintf("Available Rules (%d)", len(allRules))).
		Options(ruleOptions...).
		Value(&selectedNames).
		Run()

	// Update config
	for _, r := range allRules {
		isSelected := false
		for _, name := range selectedNames {
			if name == r.Name() {
				isSelected = true
				break
			}
		}
		im.config.SetRuleDisabled(r.Name(), !isSelected)
	}
	_ = im.config.Save()

	return im.config.GetEnabledRules(allRules)
}

func (im *InteractiveMode) applyTransformations(workingDir string, files []scanner.FileInfo, enabledRules []rules.Rule) {
	t := transformer.New(enabledRules)
	results := t.TransformAll(files)

	var changed []transformer.Result
	for _, r := range results {
		if r.Changed {
			changed = append(changed, r)
		}
	}

	if len(changed) == 0 {
		fmt.Println()
		fmt.Println(Success("All files are already up to date!"))
		return
	}

	// Show summary
	fmt.Println()
	fmt.Println(Divider())
	fmt.Println(TitleStyle.Render(fmt.Sprintf("📝 %d file(s) to update", len(changed))))
	fmt.Println()

	for _, r := range changed {
		rel, _ := filepath.Rel(workingDir, r.File.Path)
		fmt.Printf("  %s %s\n", FileStyle.Render("→"), rel)
		for _, rule := range r.AppliedRules {
			fmt.Printf("    %s\n", RuleStyle.Render("+ "+rule.Description))
		}
	}

	fmt.Println()

	// Confirm
	var confirm bool
	err := huh.NewConfirm().
		Title("Apply these changes?").
		Affirmative("✓ Yes, apply").
		Negative("✗ Cancel").
		Value(&confirm).
		Run()

	if err != nil || !confirm {
		fmt.Println(SubtitleStyle.Render("Cancelled"))
		return
	}

	// Backup if enabled
	if im.config.BackupEnabled {
		im.backupMgr = backup.New(workingDir)
		for _, r := range changed {
			_ = im.backupMgr.Backup(r.File.Path, r.File.Content)
		}
		fmt.Println(InfoStyle.Render("📦 Backup created: ") + SubtitleStyle.Render(im.backupMgr.BackupDir()))
	}

	// Apply changes
	for _, r := range changed {
		_ = os.WriteFile(r.File.Path, []byte(r.NewContent), 0644)
	}

	fmt.Println()
	fmt.Println(Success(fmt.Sprintf("Updated %d file(s) successfully!", len(changed))))
}

func (im *InteractiveMode) showSettings() {
	fmt.Println()
	fmt.Println(TitleStyle.Render("⚙ Settings"))
	fmt.Println()

	var version string
	var safeOnly bool = im.config.SafeOnly
	var backupEnabled bool = im.config.BackupEnabled

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("C# Version").
				Description("Select your target C# version").
				Options(
					huh.NewOption("C# 6   (.NET 4.6+)", "6"),
					huh.NewOption("C# 7   (.NET Core 2.0+)", "7"),
					huh.NewOption("C# 8   (.NET Core 3.0+)", "8"),
					huh.NewOption("C# 9   (.NET 5.0+)", "9"),
					huh.NewOption("C# 10  (.NET 6.0+)", "10"),
					huh.NewOption("C# 11  (.NET 7.0+)", "11"),
					huh.NewOption("C# 12  (.NET 8.0+)", "12"),
					huh.NewOption("C# 13  (.NET 9.0+)", "13"),
				).
				Value(&version),
			huh.NewConfirm().
				Title("Safe Mode").
				Description("Only apply safe transformations (recommended)").
				Value(&safeOnly),
			huh.NewConfirm().
				Title("Create Backups").
				Description("Backup files before modifying").
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

	fmt.Println()
	fmt.Println(Success("Settings saved!"))
	fmt.Println()
}

func (im *InteractiveMode) showRules() {
	groups := im.registry.GroupByVersion()
	versions := []rules.CSharpVersion{
		rules.CSharp6, rules.CSharp7, rules.CSharp8,
		rules.CSharp9, rules.CSharp10, rules.CSharp11,
		rules.CSharp12, rules.CSharp13,
	}

	fmt.Println()
	fmt.Println(TitleStyle.Render("📋 Available Rules"))
	fmt.Println()

	totalRules := 0
	for _, v := range versions {
		ruleList, ok := groups[v]
		if !ok || len(ruleList) == 0 {
			continue
		}

		// Version header with count
		isActive := v <= im.config.GetVersion()
		header := fmt.Sprintf("%s (%d rules)", v.String(), len(ruleList))
		if isActive {
			fmt.Println(AccentStyle.Render("▸ " + header))
		} else {
			fmt.Println(SubtitleStyle.Render("  " + header))
		}

		if isActive {
			for _, r := range ruleList {
				vr := r.(rules.VersionedRule)
				icon := SuccessStyle.Render("✓")
				if !vr.IsSafe() {
					icon = WarningStyle.Render("⚠")
				}
				fmt.Printf("    %s %-28s %s\n", icon, r.Name(), SubtitleStyle.Render(r.Description()))
			}
		}
		fmt.Println()
		totalRules += len(ruleList)
	}

	fmt.Println(SubtitleStyle.Render(fmt.Sprintf("  Total: %d rules", totalRules)))
	fmt.Println()

	var cont bool
	_ = huh.NewConfirm().Title("").Affirmative("← Back").Negative("").Value(&cont).Run()
}
