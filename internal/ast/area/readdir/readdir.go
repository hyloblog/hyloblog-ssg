package readdir

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyloblog/hyloblog-ssg/internal/assert"
)

type Directory interface {
	Files() []File
	Directories() []File
}

type File interface {
	Path() string
}

type directory struct {
	files       []File
	directories []File
}

func (d *directory) Files() []File       { return d.files }
func (d *directory) Directories() []File { return d.directories }

func ReadDir(dir string) (Directory, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("os read error: %w", err)
	}
	return &directory{
		getfiles(dir, entries), getdirectories(dir, entries),
	}, nil
}

type file string

func NewFile(path string) File { return file(path) }

func (f file) Path() string { return string(f) }

func getfiles(dir string, entries []os.DirEntry) []File {
	var files []File
	for _, e := range entries {
		if !e.IsDir() {
			t := e.Type()
			assert.Printf(t.IsRegular(), "unknown file type %v", t)
			files = append(files, file(filepath.Join(dir, e.Name())))
		}
	}
	return files
}

func getdirectories(dir string, entries []os.DirEntry) []File {
	var dirs []File
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, file(filepath.Join(dir, e.Name())))
		}
	}
	return dirs
}
