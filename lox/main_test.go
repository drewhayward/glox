package lox

import (
	"github.com/gkampitakis/go-snaps/snaps"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	v := m.Run()
	snaps.Clean(m)
	os.Exit(v)
}
