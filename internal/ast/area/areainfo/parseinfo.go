package areainfo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hyloblog/hyloblog-ssg/internal/assert"
)

type ParseInfo struct {
	ign         map[string]bool
	gitdir      string
	chromastyle string
}

func NewParseInfo(chromastyle string) *ParseInfo {
	return &ParseInfo{map[string]bool{}, "", chromastyle}
}

func (info *ParseInfo) Descend(dir, ignorefile string) (*ParseInfo, error) {
	ign, err := augmentign(info.ign, filepath.Join(dir, ignorefile))
	if err != nil {
		return nil, fmt.Errorf("cannot augment ign: %w", err)
	}
	gitdir, err := augmentgitdir(info.gitdir, dir)
	if err != nil {
		return nil, fmt.Errorf("cannot check for gitdir: %w", err)
	}
	return &ParseInfo{ign, gitdir, info.chromastyle}, nil
}

func augmentign(oldign map[string]bool, path string) (map[string]bool, error) {
	ign, err := parseignorefile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot get from file: %w", err)
	}
	for k, old := range oldign {
		if _, ok := ign[k]; !ok {
			ign[k] = old
		}
	}
	return ign, nil
}

func parseignorefile(path string) (map[string]bool, error) {
	lines, err := getignorelines(path)
	if err != nil {
		return nil, fmt.Errorf("cannot get ignore lines: %w", err)
	}
	ignore := map[string]bool{}
	for _, line := range lines {
		inst := parseinstruction(line)
		ignore[inst.name] = inst.shouldignore
	}
	return ignore, nil
}

func getignorelines(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("cannot read file: %w", err)
	}
	trimmed := strings.TrimSpace(string(b))
	if len(trimmed) == 0 {
		return []string{}, nil
	}
	return strings.Split(trimmed, "\n"), nil
}

type instruction struct {
	name         string
	shouldignore bool
}

func parseinstruction(s string) *instruction {
	assert.Assert(len(s) > 0)
	if s[0] == '!' {
		return &instruction{s[1:], false}
	}
	return &instruction{s, true}
}

func augmentgitdir(gitdir, path string) (string, error) {
	if gitdir != "" {
		return gitdir, nil
	}
	gitdir = filepath.Join(path, ".git")
	is, err := isgitdir(gitdir)
	if err != nil {
		return "", fmt.Errorf("cannot check if gitdir: %w", err)
	}
	if is {
		return gitdir, nil
	}
	return "", nil
}

func isgitdir(gitdir string) (bool, error) {
	stat, err := os.Stat(gitdir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return stat.IsDir(), nil
}

func (info *ParseInfo) ShouldIgnore(name string) bool {
	shouldignore, ok := info.ign[name]
	return ok && shouldignore
}

func (info *ParseInfo) GitDir() (string, bool) {
	return info.gitdir, info.gitdir != ""
}

func (info *ParseInfo) ChromaStyle() string {
	return info.chromastyle
}
