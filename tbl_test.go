package main

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestGetPage(t *testing.T) {
	var firstFullPage [100][]byte
	var secondFullPage [100][]byte

	createExpectedPage := func(numRows int) []byte {
		page := make([]byte, RowSize)
		for i := range numRows {
			binary.LittleEndian.PutUint32(page[IdOffset:], uint32(i))
			copy(page[UsernameOffset:UsernameOffset+UsernameSize], "test")
			copy(page[EmailOffset:EmailOffset+EmailSize], "test@test.gmail")
		}
		return page
	}

	var fullRow [PageSize]byte
	for i := range RowsPerPage {
		rowStart := i * RowSize
		binary.LittleEndian.PutUint32(fullRow[rowStart:], uint32(i))
		copy(fullRow[rowStart+UsernameOffset:rowStart+UsernameOffset+UsernameSize], "test")
		copy(fullRow[rowStart+EmailOffset:rowStart+EmailOffset+EmailSize], "test@test.gmail")
	}

	firstFullPage[0] = fullRow[:]
	secondFullPage[0] = fullRow[:]
	secondFullPage[1] = fullRow[:]

	var firstPartialPage [100][]byte
	var secondPartialPage [100][]byte

	var partialRow [PageSize / 2]byte
	for i := range RowsPerPage / 2 {
		rowStart := i * RowSize
		binary.LittleEndian.PutUint32(partialRow[rowStart:], uint32(i))
		copy(partialRow[rowStart+UsernameOffset:rowStart+UsernameOffset+UsernameSize], "test")
		copy(partialRow[rowStart+EmailOffset:rowStart+EmailOffset+EmailSize], "test@test.gmail")
	}

	firstPartialPage[0] = partialRow[:]
	secondPartialPage[0] = fullRow[:]
	secondPartialPage[1] = partialRow[:]

	//read from table tests
	for _, tt := range []struct {
		name    string
		pageNum int
		tbl     Table
		want    []byte
	}{
		{
			"request first full page",
			1,
			Table{NumRows: uint32(RowsPerPage), Pager: &Pager{Pages: firstFullPage}},
			createExpectedPage(RowsPerPage),
		},
		{
			"request second full page",
			2,
			Table{NumRows: uint32(RowsPerPage * 2), Pager: &Pager{Pages: secondFullPage}},
			createExpectedPage(RowsPerPage * 2),
		},
		{
			"request first partial page",
			2,
			Table{NumRows: uint32(RowsPerPage / 2), Pager: &Pager{Pages: firstPartialPage}},
			createExpectedPage(RowsPerPage / 2),
		},
		{
			"request second partial page",
			2,
			Table{NumRows: uint32(RowsPerPage + (RowsPerPage / 2)), Pager: &Pager{Pages: secondPartialPage}},
			createExpectedPage(RowsPerPage + (RowsPerPage / 2)),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.tbl.Pager.GetPage(tt.pageNum)

			if !bytes.Equal(got, tt.want) {
				t.Errorf("test: %s failed. page number: %d, got: %v,  want: %v", tt.name, tt.pageNum, got, tt.want)
			}

		})
	}

	//fail tests
	for _, tt := range []struct {
		name  string
		input int
		tbl   Table
	}{
		{"fail on out of bounds page request [negative value]", -1, Table{}},
		{"fail on out of bounds page request [positive value]", 101, Table{}},
		{"fail on request of empty page", 30, Table{}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.tbl.Pager.GetPage(tt.input); err == nil {
				t.Errorf("test: %s failed", tt.name)
			}
		})
	}

}
