package scanner

import (
	"os"
	"path/filepath"
	"strings"
)


type FileInfo struct {
	Path    string
	Content string
}


type Scanner interface {
	Scan(root string) ([]FileInfo, error)
}


type CSharpScanner struct {
	extensions []string
}


func New() *CSharpScanner {
	return &CSharpScanner{
		extensions: []string{".cs"},
	}
}


func (s *CSharpScanner) Scan(root string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			
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
