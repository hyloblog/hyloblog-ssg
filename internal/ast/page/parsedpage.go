package page

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hyloblog/hyloblog-ssg/internal/assert"
	"github.com/hyloblog/hyloblog-ssg/internal/ast/area/sitefile"
	"github.com/hyloblog/hyloblog-ssg/internal/ast/page/pandoc"
	"github.com/hyloblog/hyloblog-ssg/internal/theme"
)

type parsedpage struct {
	title, url string
	timing     *timing
	doc        string
	rawmd      string
	a          authoring
}

func ParsePage(path, chromastyle string) (Page, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}
	components, err := separate(string(buf))
	if err != nil {
		return nil, fmt.Errorf("cannot separate: %w", err)
	}
	mdpage, err := parsemdpage(components.content, chromastyle)
	if err != nil {
		return nil, fmt.Errorf("cannot parse content: %w", err)
	}
	m, err := parsemetadata(components.metadata)
	if err != nil {
		return nil, fmt.Errorf("cannot parse metadata: %w", err)
	}
	return &parsedpage{
		title:  mdpage.title,
		url:    m.URL,
		timing: m.timing(),
		doc:    mdpage.content,
		rawmd:  components.content,
		a:      *m.authoring(),
	}, nil
}

type components struct {
	metadata string
	content  string
}

func separate(s string) (*components, error) {
	s = strings.TrimSpace(s)
	if len(s) < 3 || s[:3] != "---" {
		return &components{"", s}, nil
	}
	s = s[3:]
	endmeta := strings.Index(s, "---")
	if endmeta == -1 {
		return nil, fmt.Errorf("unclosed metadata section")
	}
	return &components{
		strings.TrimSpace(s[:endmeta]),
		strings.TrimSpace(s[endmeta+3:]),
	}, nil
}

func (p *parsedpage) Title() (string, error) {
	return p.title, nil
}

func (pg *parsedpage) Link(path string, pi PageInfo) (string, error) {
	dynamiclinks := pi.DynamicLinks()
	if pg.url != "" {
		if dynamiclinks {
			return pg.url, nil
		}
		log.Printf("warning: %q has custom url in static mode\n", path)
	}
	url, err := filepath.Rel(pi.Root(), rightextpath(path, dynamiclinks))
	if err != nil {
		return "", fmt.Errorf("cannot get relative path: %w", err)
	}
	return "/" + url, nil
}

func rightextpath(path string, dynamiclinks bool) string {
	return replaceext(path, rightext(dynamiclinks))
}

func rightext(dynamiclinks bool) string {
	if dynamiclinks {
		return ""
	}
	return ".html"
}

func replaceext(path, newext string) string {
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)] + newext
}

func ParsePageGit(path, gitdir, chromastyle string) (Page, error) {
	pg, err := ParsePage(path, chromastyle)
	if err != nil {
		return nil, err
	}
	ppg, ok := pg.(*parsedpage)
	assert.Assert(ok)
	info, err := getgitinfo(path, gitdir)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot get timing from git: %w", err,
		)
	}
	if ppg.timing == nil {
		ppg.timing = &info.timing
	}
	ppg.a.addgitauthor(info.author)
	return ppg, nil
}

type gitinfo struct {
	timing timing
	author string
}

type timing struct {
	published, updated time.Time
}

func getgitinfo(path, gitdir string) (*gitinfo, error) {
	repo, err := git.PlainOpen(gitdir)
	if err != nil {
		return nil, fmt.Errorf("cannot open repo: %w", err)
	}
	rel, err := filepath.Rel(filepath.Join(gitdir, ".."), path)
	if err != nil {
		return nil, fmt.Errorf("cannot get relative path: %w", err)
	}
	log, err := repo.Log(&git.LogOptions{FileName: &rel})
	if err != nil {
		return nil, fmt.Errorf("cannot get log: %w", err)
	}
	defer log.Close()
	commits := getcommits(log)
	published, err := getcreated(commits)
	if err != nil {
		return nil, fmt.Errorf("cannot get created: %w", err)
	}
	author, err := getauthor(commits)
	if err != nil {
		return nil, fmt.Errorf("cannot get author: %w", err)
	}
	updated, err := getupdated(commits)
	if err != nil {
		return nil, fmt.Errorf("cannot get updated: %w", err)
	}
	return &gitinfo{timing{*published, *updated}, author}, nil
}

func getcommits(iter object.CommitIter) []object.Commit {
	var commits []object.Commit
	err := iter.ForEach(func(c *object.Commit) error {
		commits = append(commits, *c)
		return nil
	})
	assert.Assert(err == nil)
	return commits
}

func getcreated(commits []object.Commit) (*time.Time, error) {
	var t time.Time
	for _, c := range commits {
		when := c.Author.When
		if t.IsZero() || when.Before(t) {
			t = when
		}
	}
	return &t, nil
}

func getauthor(commits []object.Commit) (string, error) {
	if len(commits) == 0 {
		return "", nil
	}
	return commits[0].Author.Name, nil
}

func getupdated(commits []object.Commit) (*time.Time, error) {
	var t time.Time
	for _, c := range commits {
		when := c.Author.When
		if t.IsZero() || when.After(t) {
			t = when
		}
	}
	return &t, nil
}

func (pg *parsedpage) GenerateIndex(w io.Writer, posts []Post, pi PageInfo) error {
	return pi.Theme().ExecuteIndex(w, &theme.IndexData{
		Title:   pg.title,
		Content: pg.doc,
		Posts:   tothemeposts(posts, pg),
		Head:    pi.Head(),
		Foot:    pi.Foot(),
	})
}

type Post struct {
	title, category, link string
	timing                *timing
	a                     authoring
}

func tothemeposts(posts []Post, index *parsedpage) []theme.Post {
	sort.Slice(posts, func(i, j int) bool {
		t0, t1 := posts[i].timing, posts[j].timing
		return t0 != nil && t1 != nil &&
			t0.published.After(t1.published)
	})
	themeposts := make([]theme.Post, len(posts))
	for i, p := range posts {
		themeposts[i] = theme.Post{
			Title:    p.title,
			Category: p.category,
			Link:     p.link,
			Date:     getdate(p.timing),
			Authors:  p.a.getauthors(&index.a),
		}
	}
	return themeposts
}

func (pg *parsedpage) IsPost() bool { return true }

func (pg *parsedpage) AsPost(category, link string) *Post {
	return &Post{pg.title, category, link, pg.timing, pg.a}
}

func (pg *parsedpage) GenerateWithoutIndex(w io.Writer, pi PageInfo) error {
	return pi.Theme().ExecuteDefault(w, &theme.DefaultData{
		Title:   pg.title,
		Content: pg.doc,
		Date:    getdate(pg.timing),
		Authors: pg.a.getauthorsnoindex(),
		Head:    pi.Head(),
		Foot:    pi.Foot(),
	})
}

func getdate(t *timing) string {
	if t == nil || t.published.IsZero() {
		return ""
	}
	return t.published.Format("Jan 02, 2006")
}

func (pg *parsedpage) time() (time.Time, bool) {
	t := pg.timing
	if t == nil || t.published.IsZero() {
		return time.Time{}, false
	}
	return t.published, true
}

func (pg *parsedpage) ToResource(
	pagepath, emailhtmlpath, emailtextpath string,
) (sitefile.Resource, error) {
	if time, ok := pg.time(); ok {
		return sitefile.NewPostResource(
			pagepath,
			sitefile.NewTimedPost(
				pg.title, emailhtmlpath, emailtextpath, time,
			),
		), nil
	}
	return sitefile.NewPostResource(
		pagepath,
		sitefile.NewPost(pg.title, emailhtmlpath, emailtextpath),
	), nil
}

func (pg *parsedpage) GenerateEmailHtml(
	w io.Writer, pi PageInfo,
) error {
	if err := pi.Theme().ExecuteDefault(
		w, &theme.DefaultData{
			Title:   pg.title,
			Content: pg.doc,
			Date:    getdate(pg.timing),
			Authors: pg.a.getauthorsnoindex(),
			Head:    "",
			Foot:    "",
		},
	); err != nil {
		return fmt.Errorf("cannot execute: %w", err)
	}
	return nil
}

func (pg *parsedpage) GenerateEmailText(w io.Writer) error {
	return pandoc.ConvertPlaintext(pg.rawmd, w)
}

func (pg *parsedpage) Generate(w io.Writer, pi PageInfo, index Page) error {
	assert.Assert(index != nil)
	indexppg, ok := index.(*parsedpage)
	assert.Assert(ok)

	return pi.Theme().ExecuteDefault(w, &theme.DefaultData{
		Title:     pg.title,
		Content:   pg.doc,
		SiteTitle: indexppg.title,
		Date:      getdate(pg.timing),
		Authors:   pg.a.getauthors(&indexppg.a),
		Head:      pi.Head(),
		Foot:      pi.Foot(),
	})
}
