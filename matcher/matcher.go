package matcher

import (
	"cmp"
	"regexp"
	"strings"

	"github.com/dusk125/goutil/v4/lockable"
)

type MatchFunc func(vars Vars, msg interface{})

type matchEntry struct {
	vars    []string
	exp     *regexp.Regexp
	length  int
	handler MatchFunc
}

type Vars map[string]string

type Matcher struct {
	entries lockable.List[matchEntry]
}

func (m *Matcher) Register(pattern string, handler MatchFunc) {
	pieces := strings.Split(pattern, ".")
	entry := matchEntry{
		vars:    make([]string, 0),
		length:  len(pieces),
		handler: handler,
	}

	for i, piece := range pieces {
		if strings.HasPrefix(piece, "{") && strings.HasSuffix(piece, "}") {
			entry.vars = append(entry.vars, piece[1:len(piece)-1])
			pieces[i] = `(\w{1,})`
		}
	}

	entry.exp = regexp.MustCompile(strings.Join(pieces, `\.`))

	m.entries.Safe(true, func() {
		if m.entries.UnsafeNil() {
			m.entries.Make(0)
		}
		m.entries.UnsafeAppend(entry)
		m.entries.UnsafeSort(func(i, j matchEntry) int {
			return cmp.Compare(i.length, j.length)
		})
	})
}

func (m *Matcher) Call(path string, msg interface{}) {
	m.entries.Foreach(func(index int, val matchEntry) (shouldBreak bool) {
		matches := val.exp.FindStringSubmatch(path)
		if matches != nil {
			vars := make(Vars)
			for i, item := range matches[1:] {
				vars[val.vars[i]] = item
			}
			val.handler(vars, msg)
			return true
		}
		return
	})
}
