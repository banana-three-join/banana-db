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

func (stmt *Statement) Prepare(input string) (string, error) {
	getStmtErrorMsg := func(input string) string {
		return fmt.Sprintf(`"%s STATEMENT FAILED TO PREPARE"`, input)
	}

	getStmtSuccessMsg := func(input string) string {
		return fmt.Sprintf(`"%s STATEMENT SUCCESSFULLLY PREPARED"`, input)
	}

	tokens := strings.Split(input, " ")
	tokensLen := len(tokens)

	switch {
	case tokensLen <= 0:
		return "", fmt.Errorf("statement can't be empty")
	case tokensLen == 1:
		return "", fmt.Errorf("not enough values provided")
	}

	t, id := tokens[0], tokens[1]
	t = strings.ToUpper(t)

	uid, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}
	stmt.r.Id = uint32(uid)

	switch StmtType(t) {
	case Insert:
		if err := stmt.prepInsert(tokens); err != nil {
			return getStmtErrorMsg(input), err
		}
		return getStmtSuccessMsg(input), nil
	case Select:
		if err := stmt.prepSelect(); err != nil {
			return getStmtErrorMsg(input), err
		}
		return getStmtSuccessMsg(input), nil
	default:
		return getStmtErrorMsg(input), fmt.Errorf("statement's action doesn't exist")
	}
}

// keep this function like this because select is going to be added a wider set of functionalities later on
func (stmt *Statement) prepSelect() error {
	stmt.t = Select
	return nil
}

func (stmt *Statement) prepInsert(values []string) error {
	if len(values) < 3 {
		return fmt.Errorf("not enough statement values provided for insert")
	}

	username, email := values[2], values[3]

	stmt.t = Insert
	stmt.r.Username = username
	stmt.r.Email = email

	return nil
}

func (stmt *Statement) Execute(tbl *Table) (string, error) {
	getStmtErrorMsg := func(input string) string {
		return fmt.Sprintf("%s FAILED TO EXECUTE", input)
	}

	getStmtSuccessMsg := func(input string) string {
		return fmt.Sprintf("%s STATEMENT SUCCESSFULLY EXECUTED", input)
	}

	if tbl.Pager == nil {
		p, err := NewPager("temp")
		if err != nil {
			return "", err
		}

		tbl.Pager = p
	}

	switch stmt.t {
	case Insert:
		if err := stmt.execInsert(tbl); err != nil {
			return getStmtErrorMsg(stmt.String()), err
		}

		return getStmtSuccessMsg(stmt.String()), nil
	case Select:
		if err := stmt.execSelect(tbl); err != nil {
			return getStmtErrorMsg(stmt.String()), err
		}

		return getStmtSuccessMsg(stmt.String()), nil
	default:
		return getStmtErrorMsg(stmt.String()), fmt.Errorf("statement's action doesn't exist")
	}
}

func (stmt *Statement) execInsert(tbl *Table) error {
	if err := tbl.Pager.SetRow(tbl.NumRows, stmt.r); err != nil {
		return err
	}

	tbl.NumRows++
	return nil
}

func (stmt *Statement) execSelect(tbl *Table) error {
	rows, err := tbl.Pager.GetPage(int(stmt.r.Id)) //row id is currently page
	if err != nil {
		return err
	}

	var tempRow Row
	for _, r := range SplitRowsFromPage(rows) {
		if err := Deserialize(r, &tempRow); err != nil {
			return err
		}

		fmt.Printf("Index: %d, Username: %s, Email: %s", tempRow.Id, tempRow.Username, tempRow.Email)
	}

	return nil
}

func (stmt *Statement) String() string {
	return fmt.Sprintf(`STATEMENT %s %s`, stmt.t, stmt.r.String())
}
