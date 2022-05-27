package snowflake

import (
	"strings"

	"github.com/twinj/uuid"
)

func New(prefixes ...string) (s string) {
	id := strings.ReplaceAll(uuid.NewV4().String(), "-", "")

	for _, p := range prefixes {
		s += p + "_"
	}
	return s + id
}
