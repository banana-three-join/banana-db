package main

import (
	"fmt"
)

type Table struct {
	NumRows uint32
	Pages   [TableMaxPages][]byte
}

type Row struct {
	Id       uint32
	Username string
	Email    string
}

func (t *Table) GetRowsByPage(pageNum int) ([][]byte, error) {

	if pageNum <= 0 || pageNum >= TableMaxPages {
		return nil, fmt.Errorf("Page number is out of bounds")
	}

	if t.NumRows == 0 {
		return nil, fmt.Errorf("Table is currently empty")
	}

	pageIndex := pageNum - 1

	if t.Pages[pageIndex] == nil {
		return nil, fmt.Errorf("Page %d is empty", pageNum)
	}

	startRowIndex := pageIndex * RowsPerPage
	numRowsInPage := t.NumRows - uint32(startRowIndex)

	rows := make([][]byte, numRowsInPage)

	for i := range numRowsInPage {
		byteOffset := i * RowSize
		rows[i] = t.Pages[pageIndex][byteOffset : byteOffset+RowSize]
	}

	return rows, nil
}

func (t *Table) GetRowByNum(rowNum int) ([]byte, error) {
	pageNumber := rowNum / RowsPerPage

	if pageNumber >= TableMaxPages {
		return nil, fmt.Errorf("Row number is out of bounds")
	}

	rowNumToPage := rowNum % RowsPerPage
	byteOffset := RowSize * rowNumToPage

	return t.Pages[pageNumber][byteOffset : byteOffset+RowSize], nil
}
