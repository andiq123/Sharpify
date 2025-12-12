package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	backup := "off"
	if im.config.BackupEnabled {
		backup = "on"
	}

	fmt.Printf("  %s | %s mode | backup %s\n\n",
		SubtitleStyle.Render(version.String()),
		mode,
		backup)
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

	form.Run()
	return action
}

func (im *InteractiveMode) runImprove() {
	cwd, _ := os.Getwd()

	files, err := im.scanner.Scan(cwd)
	if err != nil {
		fmt.Println(ErrorStyle.Render("Scan failed: " + err.Error()))
		return
	}

	if len(files) == 0 {
		fmt.Println(WarningStyle.Render("No C# files found"))
		return
	}

	fmt.Printf("\nFound %d file(s)\n", len(files))

	availableRules := im.registry.GetByVersion(im.config.GetVersion(), im.config.SafeOnly)
	t := transformer.New(availableRules)
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
		rel, _ := filepath.Rel(cwd, r.File.Path)
		fmt.Printf("  %s\n", rel)
		for _, rule := range r.AppliedRules {
			fmt.Printf("    %s\n", RuleStyle.Render("+ "+rule.Description))
		}
	}

	var confirm bool
	huh.NewConfirm().
		Title("Apply changes?").
		Value(&confirm).
		Run()

	if !confirm {
		fmt.Println("Cancelled")
		return
	}

	if im.config.BackupEnabled {
		im.backupMgr = backup.New(cwd)
		for _, r := range changed {
			im.backupMgr.Backup(r.File.Path, r.File.Content)
		}
		fmt.Printf("Backup: %s\n", im.backupMgr.BackupDir())
	}

	for _, r := range changed {
		os.WriteFile(r.File.Path, []byte(r.NewContent), 0644)
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
	im.config.Save()

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
	huh.NewConfirm().Title("").Affirmative("OK").Negative("").Value(&cont).Run()
}

func (im *InteractiveMode) showDiff(old, new string) {
	oldLines := strings.Split(old, "\n")
	newLines := strings.Split(new, "\n")

	maxLines := len(oldLines)
	if len(newLines) > maxLines {
		maxLines = len(newLines)
	}

	shown := 0
	for i := 0; i < maxLines && shown < 30; i++ {
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
				fmt.Println(DiffRemoveStyle.Render("- " + oldLine))
				shown++
			}
			if newLine != "" {
				fmt.Println(DiffAddStyle.Render("+ " + newLine))
				shown++
			}
		}
	}
}
