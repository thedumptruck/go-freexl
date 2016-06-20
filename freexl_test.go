package freexl

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestOpenFile(t *testing.T) {
	glob, err := filepath.Glob("./_xls/*.xls")
	if err != nil {
		t.Error(err)
	}
	for _, name := range glob {
		_, err := Open(name)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestOpenBinary(t *testing.T) {
	glob, err := filepath.Glob("./_xls/*.xls")
	if err != nil {
		t.Error(err)
	}
	for _, name := range glob {
		b, err := ioutil.ReadFile(name)
		if err != nil {
			t.Error(err)
		}
		_, err = OpenBinary(b)
		if err != nil {
			t.Error(err)
		}
	}
}
