package hook

import (
	"fmt"
	"strings"
)

func formatLabels(m map[string]string) string {
	var s string
	for k, v := range m {
		s += fmt.Sprintf("%s=%s,", k, v)
	}
	return strings.TrimRight(s, ",")
}
