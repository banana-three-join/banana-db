package main

import (
	"fmt"
)

type Table struct {
	NumRows uint32
	Pager   *Pager
}

// struct to abstract away interacting directly with []byte
type Row struct {
	Id       uint32
	Username string
	Email    string
}

/*
	numrows provides tracking of the stage of the table
	if functions don't take into account the numrows then they could be out of sync and throw off the checks
	it would be better if they were part of the pages as a header to allocate all of the information about the rows in one place
*/

func (r *Row) String() string {
	return fmt.Sprintf(`%d %s %s`, r.Id, r.Username, r.Email)
}
