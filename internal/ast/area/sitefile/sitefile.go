package sitefile

import (
	"time"

	"github.com/hyloblog/hyloblog-ssg/internal/assert"
)

type Resource interface {
	Path() string
	IsPost() bool
	Post() Post
}

type Post interface {
	Title() string
	Time() (time.Time, bool)
	HtmlPath() string
	PlaintextPath() string
}

type file struct {
	path   string
	ispost bool
	post   Post
}

func NewPostResource(path string, post Post) Resource {
	return &file{path, true, post}
}

func NewNonPostResource(path string) Resource {
	return &file{path, false, nil}
}

func (f *file) Path() string { return f.path }
func (f *file) IsPost() bool { return f.ispost }
func (f *file) Post() Post {
	assert.Assert(f.IsPost())
	return f.post
}

type post struct {
	title                   string
	htmlpath, plaintextpath string
	time                    time.Time
}

func NewPost(title, htmlpath, plaintextpath string) *post {
	return &post{title, htmlpath, plaintextpath, time.Time{}}
}

func NewTimedPost(title, htmlpath, plaintextpath string, time time.Time) *post {
	return &post{title, htmlpath, plaintextpath, time}
}

func (p *post) Title() string {
	return p.title
}

func (p *post) HtmlPath() string {
	return p.htmlpath
}

func (p *post) PlaintextPath() string {
	return p.plaintextpath
}

func (p *post) Time() (time.Time, bool) {
	return p.time, !p.time.IsZero()
}
