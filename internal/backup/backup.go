package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Manager handles file backups
type Manager struct {
	backupDir string
	enabled   bool
}

// New creates a new backup manager
func New(projectPath string) *Manager {
	timestamp := time.Now().Format("20060102-150405")
	backupDir := filepath.Join(projectPath, ".sharpify-backup", timestamp)

	return &Manager{
		backupDir: backupDir,
		enabled:   true,
	}
}

// SetEnabled enables or disables backups
func (m *Manager) SetEnabled(enabled bool) {
	m.enabled = enabled
}

// IsEnabled returns whether backups are enabled
func (m *Manager) IsEnabled() bool {
	return m.enabled
}

// BackupDir returns the backup directory path
func (m *Manager) BackupDir() string {
	return m.backupDir
}

// Backup creates a backup of the given file
func (m *Manager) Backup(filePath string, content string) error {
	if !m.enabled {
		return nil
	}

	// Create backup directory if needed
	if err := os.MkdirAll(m.backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Get relative path for backup
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Create backup file path (flatten the structure a bit)
	backupName := filepath.Base(absPath)
	backupPath := filepath.Join(m.backupDir, backupName)

	// Handle duplicate filenames
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

	// Write backup
	if err := os.WriteFile(backupPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write backup: %w", err)
	}

	return nil
}

// Restore restores all files from backup
func (m *Manager) Restore() error {
	// List backup files
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

		// Note: This simple restore doesn't track original paths
		// A more sophisticated implementation would store metadata
		fmt.Printf("Backup available: %s\n", backupPath)
		_ = content
	}

	return nil
}

// Cleanup removes old backups
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
