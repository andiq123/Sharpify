package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)


type Manager struct {
	backupDir string
	enabled   bool
}


func New(projectPath string) *Manager {
	timestamp := time.Now().Format("20060102-150405")
	backupDir := filepath.Join(projectPath, ".sharpify-backup", timestamp)

	return &Manager{
		backupDir: backupDir,
		enabled:   true,
	}
}


func (m *Manager) SetEnabled(enabled bool) {
	m.enabled = enabled
}


func (m *Manager) IsEnabled() bool {
	return m.enabled
}


func (m *Manager) BackupDir() string {
	return m.backupDir
}


func (m *Manager) Backup(filePath string, content string) error {
	if !m.enabled {
		return nil
	}

	
	if err := os.MkdirAll(m.backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	
	backupName := filepath.Base(absPath)
	backupPath := filepath.Join(m.backupDir, backupName)

	
	counter := 1
	for {
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			break
		}
		ext := filepath.Ext(backupName)
		name := backupName[:len(backupName)-len(ext)]
		backupPath = filepath.Join(m.backupDir, fmt.Sprintf("%s_%d%s", name, counter, ext))
		counter++
	}

	
	if err := os.WriteFile(backupPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write backup: %w", err)
	}

	return nil
}


func (m *Manager) Restore() error {
	
	entries, err := os.ReadDir(m.backupDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		backupPath := filepath.Join(m.backupDir, entry.Name())
		content, err := os.ReadFile(backupPath)
		if err != nil {
			return fmt.Errorf("failed to read backup file %s: %w", entry.Name(), err)
		}

		
		
		fmt.Printf("Backup available: %s\n", backupPath)
		_ = content
	}

	return nil
}


func (m *Manager) Cleanup(keepDays int) error {
	baseDir := filepath.Dir(m.backupDir)

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -keepDays)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			path := filepath.Join(baseDir, entry.Name())
			os.RemoveAll(path)
		}
	}

	return nil
}
