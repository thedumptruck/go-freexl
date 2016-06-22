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
		w, err := OpenBinary(b)
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

func BenchmarkOpen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Open("./_xls/simple2003_4WB.xlw")
	}
}

func BenchmarkOpenBinary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buff, _ := ioutil.ReadFile("./_xls/simple2003_4WB.xlw")
		OpenBinary(buff)
	}
}
