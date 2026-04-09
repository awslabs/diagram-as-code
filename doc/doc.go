// Package doc provides embedded documentation files for the diagram-as-code project.
package doc

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed *.md advanced/*.md
var FS embed.FS

// ListMarkdownFiles returns a list of all embedded markdown file paths.
func ListMarkdownFiles() ([]string, error) {
	var files []string
	err := fs.WalkDir(FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.EqualFold(filepath.Ext(path), ".md") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
