package freexl

/*
#cgo linux pkg-config: libfreexl
#cgo darwin LDFLAGS: -lfreexl
#cgo windows LDFLAGS: -lfreexl
#include <stdlib.h>
#include <stdio.h>
#include <freexl.h>
// End Import

static int freexl_cell_get_type(FreeXL_CellValue cell) {
	return cell.type;
}

static int freexl_cell_get_int_value(FreeXL_CellValue cell) {
  return cell.value.int_value;
}

static double freexl_cell_get_double_value(FreeXL_CellValue cell) {
  return cell.value.double_value;
}

static const char * freexl_cell_get_text_value(FreeXL_CellValue cell) {
  return cell.value.text_value;
}
*/
import "C"

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unsafe"
)

// Freexl handle
type Freexl struct {
	handle   unsafe.Pointer
	MaxSheet int
	Sheets   []Sheet
}

// Sheet handle
type Sheet struct {
	MaxRow int
	MaxCol int
	Name   string
	Values [][]string
}

// Open filepath
func Open(path string) (f *Freexl, err error) {
	f = &Freexl{}
	err = f.Open(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// OpenReader with temporary file
func OpenReader(r io.Reader) (f *Freexl, err error) {
	f = &Freexl{}
	tmp, err := ioutil.TempFile(os.TempDir(), "freexls")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())

	_, err = io.Copy(tmp, r)
	if err != nil {
		return nil, err
	}

	err = f.Open(tmp.Name())
	if err != nil {
		return nil, err
	}

	return f, nil
}

// OpenBinary with temporary file
func OpenBinary(b []byte) (f *Freexl, err error) {
	return OpenReader(bytes.NewReader(b))
}

// Open file with defined handle
func (f *Freexl) Open(path string) error {
	pName := C.CString(path)
	defer C.free(unsafe.Pointer(pName))
	defer C.freexl_close(f.handle)

	ret := C.freexl_open(pName, &f.handle)
	if ret != OK {
		return fmt.Errorf("Freexl: error opening %s : %d", path, ret)
	}

	info, err := f.getInfo(BIFF_SHEET_COUNT)
	if err != nil {
		return err
	}

	f.MaxSheet = info
	f.Sheets = make([]Sheet, f.MaxSheet)

	for i, sheet := range f.Sheets {
		if err = sheet.getInfo(i, f.handle); err != nil {
			return err
		}
		f.Sheets[i] = sheet
	}

	return nil
}

func (f *Freexl) getInfo(t int) (int, error) {
	var info C.uint
	ret := C.freexl_get_info(f.handle, C.ushort(t), &info)
	if ret != OK {
		return 0, fmt.Errorf("Freexl: error get info %d: %d", t, ret)
	}

	return int(info), nil
}

func (s *Sheet) getInfo(i int, handle unsafe.Pointer) error {
	name := C.CString("")
	defer C.free(unsafe.Pointer(name))

	var numRow C.uint
	var numCol C.ushort

	ret := C.freexl_get_worksheet_name(handle, C.ushort(i), &name)
	if ret != OK {
		return fmt.Errorf("Freexl: error get sheet name :%d", ret)
	}

	s.Name = C.GoString(name)

	ret = C.freexl_select_active_worksheet(handle, C.ushort(i))
	if ret != OK {
		return fmt.Errorf("Freexl: failed to selct worksheet: %s", ret)
	}

	ret = C.freexl_worksheet_dimensions(handle, &numRow, &numCol)
	if ret != OK {
		return fmt.Errorf("Freexl: failed to get dimensions: %s", ret)
	}

	s.MaxRow = int(numRow)
	s.MaxCol = int(numCol)

	s.Values = make([][]string, s.MaxRow)

	for i, row := range s.Values {
		row = make([]string, s.MaxCol)
		var cell C.FreeXL_CellValue
		for j := range row {
			ret = C.freexl_get_cell_value(handle, C.uint(i), C.ushort(j), &cell)
			if ret != OK {
				continue
			}

			switch C.freexl_cell_get_type(cell) {
			case CELL_DOUBLE:
				row[j] = fmt.Sprintf("%1.12f", float64(C.freexl_cell_get_double_value(cell)))
			case CELL_INT:
				row[j] = fmt.Sprintf("%d", int(C.freexl_cell_get_int_value(cell)))
			case CELL_TEXT, CELL_SST_TEXT, CELL_DATE, CELL_DATETIME, CELL_TIME:
				row[j] = C.GoString(C.freexl_cell_get_text_value(cell))
			default:
				fallthrough
			case CELL_NULL:
				row[j] = ""
			}
		}

		s.Values[i] = row
	}

	return nil
}
