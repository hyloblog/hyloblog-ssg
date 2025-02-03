package page

import (
	"io"

	"github.com/hyloblog/hyloblog-ssg/internal/ast/area/sitefile"
	"github.com/hyloblog/hyloblog-ssg/internal/theme"
)

type Page interface {
	Title() (string, error)
	Link(path string, pi PageInfo) (string, error)

	GenerateIndex(w io.Writer, posts []Post, pi PageInfo) error
	Generate(w io.Writer, pi PageInfo, index Page) error
	GenerateWithoutIndex(w io.Writer, pi PageInfo) error
	GenerateEmailHtml(w io.Writer, pi PageInfo) error
	GenerateEmailText(w io.Writer) error

	IsPost() bool
	AsPost(category, link string) *Post

	ToResource(
		pagepath, emailhtmlpath, emailtextpath string,
	) (sitefile.Resource, error)
}

type PageInfo interface {
	Theme() *theme.Theme
	Head() string
	Foot() string
	Root() string
	DynamicLinks() bool
}
