package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// FileInfo holds information about a C# file
type FileInfo struct {
	Path    string
	Content string
}

// Scanner defines the interface for file scanning (Interface Segregation)
type Scanner interface {
	Scan(root string) ([]FileInfo, error)
}

// CSharpScanner scans for C# files (Single Responsibility)
type CSharpScanner struct {
	extensions []string
}

// New creates a new CSharpScanner
func New() *CSharpScanner {
	return &CSharpScanner{
		extensions: []string{".cs"},
	}
}

// Scan walks the directory tree and finds all C# files
func (s *CSharpScanner) Scan(root string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Skip common non-source directories
			name := info.Name()
			if name == "bin" || name == "obj" || name == ".git" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		if s.isCSharpFile(path) {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			files = append(files, FileInfo{
				Path:    path,
				Content: string(content),
			})
		}

		return nil
	})

	return files, err
}

func (s *CSharpScanner) isCSharpFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, e := range s.extensions {
		if ext == e {
			return true
		}
	}
	return false
}
