package ssg

import (
	"errors"
	"fmt"
	"time"

	"github.com/hyloblog/hyloblog-ssg/internal/ast/area"
	"github.com/hyloblog/hyloblog-ssg/internal/ast/area/sitefile"
	"github.com/hyloblog/hyloblog-ssg/internal/theme"
)

// A Site is the outcome of generating from a hylo directory.
type Site interface {
	// The Title of the Site.
	Title() string

	// Hash is a unique, deterministic identifier for the site such that
	// equivalent source directories should produce the same result, and
	// different ones different ones.
	Hash() string

	// The Files that constitute the Site.
	Bindings() map[string]Resource
}

type site struct {
	title, hash string
	bindings    map[string]Resource
}

func (s *site) Title() string                 { return s.title }
func (s *site) Hash() string                  { return s.hash }
func (s *site) Bindings() map[string]Resource { return s.bindings }

var (
	ErrTheme = errors.New("theme error")
)

func GenerateSiteWithBindings(
	src, target, themeName, chromastyle string,
	head, foot string,
	custompages map[string]CustomPage,
) (Site, error) {
	a, err := area.ParseArea(src, chromastyle)
	if err != nil {
		return nil, fmt.Errorf("cannot parse area: %w", err)
	}
	if err := a.Inject(toinjectmap(custompages)); err != nil {
		return nil, fmt.Errorf("injection error: %w", err)
	}
	bindings, err := a.GenerateWithBindings(target, themeName, head, foot)
	if err != nil {
		if errors.Is(err, theme.ErrNoCustomPageTemplate) {
			return nil, fmt.Errorf("%w: %w", ErrTheme, err)
		}
		return nil, fmt.Errorf("cannot generate: %w", err)
	}
	h, err := a.Hash()
	if err != nil {
		return nil, fmt.Errorf("cannot get hash: %w", err)
	}
	return &site{gettitle(a), h, tofilemap(bindings)}, nil
}

func GetSiteHash(src string) (string, error) {
	const defaultChromaStyle = "based"

	a, err := area.ParseArea(src, defaultChromaStyle)
	if err != nil {
		return "", fmt.Errorf("cannot parse area: %w", err)
	}
	return a.Hash()
}

// A Resource is any URL-accessible resource in a site.
type Resource interface {
	// Path is the path on disk to the generated file.
	Path() string

	// IsPost indicates whether or not the Resource is a post.
	IsPost() bool

	// Post contains the Post-specific data of the Resource. If the Resource is not
	// a post (i.e. !IsPost()) this will result in assertion failure.
	Post() Post
}

type Post interface {
	// Title is the title of the post.
	Title() string

	// Time is the timestamp associated with the post, if available.
	Time() (time.Time, bool)

	// HtmlPath is the path on disk to a file containing the HTML body
	// of the post without the header/footer portions that render
	// ordinarily.
	HtmlPath() string

	// PlaintextPath is the path on disk to a file containing body of
	// the post rendered in plaintext by Pandoc.
	PlaintextPath() string
}

func toinjectmap(m1 map[string]CustomPage) map[string]sitefile.CustomPage {
	m2 := map[string]sitefile.CustomPage{}
	for k, v := range m1 {
		m2[k] = v
	}
	return m2
}

func gettitle(a *area.Area) string {
	if s, err := a.Title(); err == nil {
		return s
	}
	return ""
}

func tofilemap(m1 map[string]sitefile.Resource) map[string]Resource {
	m2 := map[string]Resource{}
	for k, v := range m1 {
		m2[k] = newresource(v)
	}
	return m2
}

type resource struct {
	rsc sitefile.Resource
}

func newresource(rsc sitefile.Resource) Resource {
	return &resource{rsc}
}

func (r *resource) Path() string { return r.rsc.Path() }
func (r *resource) IsPost() bool { return r.rsc.IsPost() }
func (r *resource) Post() Post   { return r.rsc.Post() }

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
