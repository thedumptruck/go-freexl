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
		t.Log(w.Sheets[0].MaxRow)
		t.Log(w.Sheets[0].MaxCol)
		// t.Log(w.Sheets[0].Values)
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
