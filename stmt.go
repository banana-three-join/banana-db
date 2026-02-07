package main

import (
	"fmt"
	"strconv"
	"strings"
)

type StmtType string

type Statement struct {
	t StmtType
	r Row
}

func (stmt *Statement) Prepare(input string) error {
	values := strings.Fields(input)

	if len(values) < 1 {
		return fmt.Errorf("Statement can't be empty")
	}

	t, id := values[0], values[1]
	t = strings.ToUpper(t)

	uid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	stmt.r.Id = uint32(uid)

	switch StmtType(t) {
	case Insert:
		return stmt.prepInsert(values)
	case Select:
		return stmt.prepSelect()
	default:
		return fmt.Errorf("invalid statement provided")
	}
}

func (stmt *Statement) prepSelect() error {
	stmt.t = Select
	return nil
}

func (stmt *Statement) prepInsert(values []string) error {
	if len(values) < 3 {
		return fmt.Errorf("not enough statement values provided")
	}

	username, email := values[2], values[3]

	stmt.t = Insert
	stmt.r.Username = username
	stmt.r.Email = email

	return nil
}

func (stmt *Statement) Execute(tbl *Table) error {
	switch stmt.t {
	case Insert:
		return stmt.execInsert(tbl)
	case Select:
		return stmt.execSelect(tbl)
	default:
		return fmt.Errorf("invalid statement action")
	}
}

func (stmt *Statement) execInsert(tbl *Table) error {
	if tbl.NumRows >= MaxRowsPerTable {
		return fmt.Errorf("table is full")
	}

	pageNumber := tbl.NumRows / RowsPerPage
	rowNumber := tbl.NumRows % RowsPerPage
	rowOffset := RowSize * rowNumber

	if tbl.Pages[pageNumber] == nil {
		tbl.Pages[pageNumber] = make([]byte, PageSize)
	}

	dst := tbl.Pages[pageNumber][rowOffset : rowOffset+RowSize]
	Serialize(stmt.r, dst)

	tbl.NumRows++

	return nil
}

func (stmt *Statement) execSelect(tbl *Table) error {
	rows, err := tbl.GetRowsByPage(int(stmt.r.Id)) //row id is currently page
	if err != nil {
		return err
	}
	var r Row

	for i := range rows {
		if err := Deserialize(rows[i], &r); err != nil {
			return err
		}

		fmt.Printf("Index: %d, Username: %s, Email: %s", r.Id, r.Username, r.Email)
	}

	return nil
}
