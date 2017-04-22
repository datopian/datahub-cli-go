package cmd

import (
	"testing"
	"strings"
)

func TestVersion(t *testing.T) {
	s := GetVersion()
	if ! strings.Contains(s, Version) {
		t.Error("Wrong version", s)
	}
}

