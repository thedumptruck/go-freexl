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
		w, err := Open(name)
		if err != nil {
			t.Error(err)
		}
		t.Logf("MaxRow: %d, MaxCol: %d, Name: %s, Count: %d",
			w.Sheets[0].MaxRow,
			w.Sheets[0].MaxCol,
			w.Sheets[0].Name,
			len(w.Sheets[0].Values))
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
