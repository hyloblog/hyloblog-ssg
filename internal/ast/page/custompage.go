package page

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/hyloblog/hyloblog-ssg/internal/assert"
	"github.com/hyloblog/hyloblog-ssg/internal/ast/area/sitefile"
)

type custompage struct {
	template string
	data     map[string]string
}

func CustomPage(template string, data map[string]string) *custompage {
	return &custompage{template, data}
}

func (pg *custompage) Title() (string, error) {
	return "", fmt.Errorf("custompage has no title")
}

func (pg *custompage) Link(path string, pi PageInfo) (string, error) {
	url, err := filepath.Rel(
		pi.Root(),
		rightextpath(path, pi.DynamicLinks()),
	)
	if err != nil {
		return "", fmt.Errorf("cannot get relative path: %w", err)
	}
	return "/" + url, nil
}

func (pg *custompage) GenerateIndex(
	w io.Writer, posts []Post, pi PageInfo,
) error {
	return fmt.Errorf("custom page cannot generate index")
}

func (pg *custompage) GenerateWithoutIndex(w io.Writer, pi PageInfo) error {
	return pi.Theme().ExecuteCustom(w, pg.template, pg.data)
}

func (pg *custompage) Generate(w io.Writer, pi PageInfo, index Page) error {
	assert.Assert(index != nil)
	indexppg, ok := index.(*parsedpage)
	assert.Assert(ok)

	return pi.Theme().ExecuteCustom(
		w,
		pg.template,
		pg.datawithmoremap(map[string]string{
			"SiteTitle": indexppg.title,
		}),
	)
}

func (pg *custompage) GenerateEmailHtml(w io.Writer, pi PageInfo) error {
	return fmt.Errorf("custom page cannot generate email")
}

func (pg *custompage) GenerateEmailText(w io.Writer) error {
	return fmt.Errorf("custom page cannot generate email")
}

func (pg *custompage) datawithmoremap(more map[string]string) map[string]string {
	m := map[string]string{}
	for k, v := range pg.data {
		m[k] = v
	}
	for k, v := range more {
		m[k] = v
	}
	return m
}

func (pg *custompage) IsPost() bool { return false }

func (pg *custompage) AsPost(_, _ string) *Post {
	assert.Assert(false)
	return nil
}

func (pg *custompage) ToResource(path, _, _ string) (sitefile.Resource, error) {
	return sitefile.NewNonPostResource(path), nil
}
