package ssg

import (
	"fmt"
	"time"

	"github.com/xr0-org/progstack-ssg/internal/ast/area"
	"github.com/xr0-org/progstack-ssg/internal/ast/area/sitefile"
)

// A File is any URL-accessible resource in a site.
type File interface {
	// Path is the path on disk to the generated File.
	Path() string

	// IsPost indicates whether or not the File is a post.
	IsPost() bool

	// Title is the title of the post.
	PostTitle() string

	// Time is the timestamp associated with the post, if available.
	PostTime() (time.Time, bool)
}

func GenerateSiteWithBindings(
	src, target, theme, chromastyle string,
	head, foot string,
	custompages map[string]CustomPage,
) (map[string]File, error) {
	a, err := area.ParseArea(src, chromastyle)
	if err != nil {
		return nil, fmt.Errorf("cannot parse area: %w", err)
	}
	if err := a.Inject(toinjectmap(custompages)); err != nil {
		return nil, fmt.Errorf("injection error: %w", err)
	}
	m1, err := a.GenerateWithBindings(target, theme, head, foot)
	if err != nil {
		return nil, fmt.Errorf("cannot generate: %w", err)
	}
	m2 := map[string]File{}
	for k, v := range m1 {
		m2[k] = v
	}
	return m2, nil
}

func toinjectmap(m1 map[string]CustomPage) map[string]sitefile.CustomPage {
	m2 := map[string]sitefile.CustomPage{}
	for k, v := range m1 {
		m2[k] = v
	}
	return m2
}

// A CustomPage is a page generated without existing in the source directory.
type CustomPage interface {
	Template() string
	Data() map[string]string
}

type custompage struct {
	template string
	data     map[string]string
}

func (pg *custompage) Template() string        { return pg.template }
func (pg *custompage) Data() map[string]string { return pg.data }

func NewSubscriberPage(url string) CustomPage {
	return &custompage{
		"subscribe.html",
		map[string]string{"FormAction": url},
	}
}

func NewMessagePage(title, message string) CustomPage {
	return &custompage{
		"message.html",
		map[string]string{"Title": title, "Message": message},
	}
}
