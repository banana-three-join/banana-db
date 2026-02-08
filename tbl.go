package main

import (
	"fmt"
)

type Table struct {
	NumRows uint32
	Pages   [MaxPagesPerTable][]byte
}

/*
	numrows provides tracking of the stage of the table
	if functions don't take into account the numrows then they could be out of sync and throw off the checks
	it would be better if they were part of the pages as a header to allocate all of the information about the rows in one place
*/

type Row struct {
	Id       uint32
	Username string
	Email    string
}

func (t *Table) GetRowsByPage(pageNum int) ([][]byte, error) {
	pageIndex := pageNum - 1

	if pageIndex < 0 || pageIndex >= MaxPagesPerTable {
		return nil, fmt.Errorf("requested page number is out of bound")
	}

	if t.Pages[pageIndex] == nil {
		return nil, fmt.Errorf("requested page is empty")
	}

	upperLimit := pageIndex*RowsPerPage + RowsPerPage
	var numRows int

	if t.NumRows >= uint32(upperLimit) {
		numRows = RowsPerPage
	} else {
		numRows = RowsPerPage - (upperLimit - int(t.NumRows))
	}

	rows := make([][]byte, numRows)
	for i := range numRows {
		offset := RowSize * i
		row := make([]byte, RowSize)
		copy(row, t.Pages[pageIndex][offset:offset+RowSize])
		rows[i] = row
	}

	return rows, nil
}
