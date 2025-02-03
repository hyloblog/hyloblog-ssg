package areainfo

import (
	"github.com/hyloblog/hyloblog-ssg/internal/assert"
	"github.com/hyloblog/hyloblog-ssg/internal/ast/page"
	"github.com/hyloblog/hyloblog-ssg/internal/theme"
)

type GenInfo struct {
	rootdir    string
	index      page.Page
	purpose    Purpose
	head, foot string
	theme      *theme.Theme
}

func (info *GenInfo) copy() *GenInfo {
	return &GenInfo{
		theme:   info.theme,
		rootdir: info.rootdir,
		index:   info.index,
		purpose: info.purpose,
		head:    info.head,
		foot:    info.foot,
	}
}

func NewGenInfo(theme *theme.Theme, rootdir string, purpose Purpose) *GenInfo {
	return &GenInfo{
		theme:   theme,
		rootdir: rootdir,
		purpose: purpose,
	}
}

func (info *GenInfo) WithNewIndex(index page.Page) *GenInfo {
	assert.Assert(index != nil)
	gi := info.copy()
	gi.index = index
	return gi
}

func (info *GenInfo) WithHeadFoot(head, foot string) *GenInfo {
	gi := info.copy()
	gi.head = head
	gi.foot = foot
	return gi
}

func (info *GenInfo) GetIndex() (page.Page, bool) {
	return info.index, info.index != nil
}

func (info *GenInfo) DynamicLinks() bool {
	switch info.purpose {
	case PurposeStaticServe:
		return false
	case PurposeDynamicServe, PurposeBind:
		return true
	default:
		assert.Assert(false)
		return false
	}
}

func (info *GenInfo) Theme() *theme.Theme { return info.theme }
func (info *GenInfo) Root() string        { return info.rootdir }
func (info *GenInfo) Head() string        { return info.head }
func (info *GenInfo) Foot() string        { return info.foot }
func (info *GenInfo) Binding() bool       { return info.purpose == PurposeBind }

type Purpose int

const (
	PurposeStaticServe Purpose = iota
	PurposeDynamicServe
	PurposeBind
)
