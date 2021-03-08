package matcher

import (
	"regexp"
	"sort"
	"strings"
	"sync"
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
	sync.RWMutex
	entries entries
}

type entries []*matchEntry

func (a entries) Len() int           { return len(a) }
func (a entries) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a entries) Less(i, j int) bool { return a[i].length < a[j].length }

func (m *Matcher) Make() {
	m.entries = make(entries, 0)
}

func (m *Matcher) Register(pattern string, handler MatchFunc) {
	m.Lock()
	defer m.Unlock()
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

	m.entries = append(m.entries, &entry)
	sort.Sort(sort.Reverse(m.entries))
}

func (m *Matcher) Call(path string, msg interface{}) {
	m.RLock()
	defer m.RUnlock()
	for _, entry := range m.entries {
		matches := entry.exp.FindStringSubmatch(path)
		if matches != nil {
			vars := make(Vars)
			for i, item := range matches[1:] {
				vars[entry.vars[i]] = item
			}
			entry.handler(vars, msg)
			break
		}
	}
}
