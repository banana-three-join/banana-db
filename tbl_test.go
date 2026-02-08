package main

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestGetRowsByPage(t *testing.T) {
	var firstFullPage [100][]byte
	var secondFullPage [100][]byte

	createExpectedRows := func(numRows int) [][]byte {
		rows := make([][]byte, numRows)
		for i := range numRows {
			row := make([]byte, RowSize)
			binary.LittleEndian.PutUint32(row[IdOffset:], uint32(i))
			copy(row[UsernameOffset:UsernameOffset+UsernameSize], "test")
			copy(row[EmailOffset:EmailOffset+EmailSize], "test@test.gmail")

			rows[i] = row
		}
		return rows
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
		want    [][]byte
	}{
		{
			"request first full page",
			1,
			Table{NumRows: uint32(RowsPerPage), Pages: firstFullPage},
			createExpectedRows(RowsPerPage),
		},
		{
			"request second full page",
			2,
			Table{NumRows: uint32(RowsPerPage * 2), Pages: secondFullPage},
			createExpectedRows(RowsPerPage * 2),
		},
		{
			"request first partial page",
			2,
			Table{NumRows: uint32(RowsPerPage / 2), Pages: firstPartialPage},
			createExpectedRows(RowsPerPage / 2),
		},
		{
			"request second partial page",
			2,
			Table{NumRows: uint32(RowsPerPage + (RowsPerPage / 2)), Pages: secondPartialPage},
			createExpectedRows(RowsPerPage + (RowsPerPage / 2)),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.tbl.GetRowsByPage(tt.pageNum)
			for i, r := range got {
				if !bytes.Equal(r, tt.want[i]) {
					t.Errorf("test: %s failed. at index: %d, got: %v is different from want: %v", tt.name, i, r, tt.want[i])
				}
			}
		})
	}

	//fail tests
	for _, tt := range []struct {
		name  string
		input int
		tbl   Table
	}{
		{"fail on out of bounds page request [negative value]", 0, Table{}},
		{"fail on out of bounds page request [positive value]", 101, Table{}},
		{"fail on request of empty page", 30, Table{}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.tbl.GetRowsByPage(tt.input); err == nil {
				t.Errorf("test: %s failed", tt.name)
			}
		})
	}

}
