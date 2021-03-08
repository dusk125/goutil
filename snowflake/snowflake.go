package snowflake

import (
	"fmt"
	"strings"

	"github.com/twinj/uuid"
)

func New(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, strings.ReplaceAll(uuid.NewV4().String(), "-", ""))
}
