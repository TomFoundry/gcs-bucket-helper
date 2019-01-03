package tilde

import (
	"strings"

	tilde "gopkg.in/mattes/go-expand-tilde.v1"
)

// ExpandPath expands a path that has a tilde prefix.
// If there is no tilde prefix then it returns path unchanged.
func ExpandPath(s string) (string, error) {

	if !strings.HasPrefix(s, "~") { // No need to expand
		return s, nil
	}

	expanded, err := tilde.Expand(s)

	if err != nil {
		return s, err
	}

	return expanded, nil
}
