package snowflake

import (
	"strings"

	"github.com/twinj/uuid"
)

// Creates a compressed ('-' removed) UUID with any given prefix(es)
func New(prefixes ...string) (s string) {
	id := strings.ReplaceAll(uuid.NewV4().String(), "-", "")

	for _, p := range prefixes {
		s += p + "_"
	}
	return s + id
}
