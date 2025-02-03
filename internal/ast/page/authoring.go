package page

import (
	"github.com/hyloblog/hyloblog-ssg/internal/theme"
)

// An authoring represents all the information regarding the author(s) of a page
// that has been hitherto gathered
type authoring struct {
	_metaauthors []string
	_metadefs    map[string]authordef
	_gitauthor   string
}

func newAuthoring(authors []string, defs map[string]authordef) *authoring {
	return &authoring{authors, defs, ""}
}

func (a *authoring) addgitauthor(author string) {
	a._gitauthor = author
}

func (a *authoring) getauthors(index *authoring) []theme.Author {
	if len(a._metaauthors) == 0 {
		if len(index._metaauthors) > 0 {
			return index.getauthorsnoindex()
		}
		return a.getauthorsnoindex()
	}
	return getauthors(
		defineauthors(a._metaauthors, a._metadefs),
		index._metadefs,
	)
}

func (a *authoring) getauthorsnoindex() []theme.Author {
	if len(a._metaauthors) == 0 {
		if len(a._gitauthor) > 0 {
			return []theme.Author{theme.Author{Name: a._gitauthor}}
		}
		return []theme.Author{}
	}
	return getauthors(
		defineauthors(a._metaauthors, a._metadefs),
		map[string]authordef{},
	)
}

func getauthors(undef []authordef, defs map[string]authordef) []theme.Author {
	var authors []theme.Author
	for _, author := range defineauthors(tostrings(undef), defs) {
		authors = append(authors, theme.Author(author))
	}
	return authors
}

func defineauthors(undef []string, defs map[string]authordef) []authordef {
	var defined []authordef
	for _, author := range undef {
		if def, ok := defs[author]; ok {
			defined = append(defined, def)
		} else {
			defined = append(defined, authordef{Name: author})
		}
	}
	return defined
}
