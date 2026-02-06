package main

import "fmt"

type Table struct {
	NumRows uint32
	Pages   [TableMaxPages][]byte
}

type Row struct {
	Id       uint32
	Username string
	Email    string
}

func (t *Table) GetRowByNum(rowNum int) ([]byte, error) {
	pageNumber := rowNum / RowsPerPage

	if pageNumber >= TableMaxPages {
		return nil, fmt.Errorf("row number out of bounds")
	}

	rowNumToPage := rowNum % RowsPerPage
	byteOffset := RowSize * rowNumToPage

	return t.Pages[pageNumber][byteOffset : byteOffset+RowSize], nil
}
