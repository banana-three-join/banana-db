/*
	DB Theory:
		B+tree - used to store key-values that require fast indexing through a quick "merge" approach
		everything is stored through concatenation. how can we tell were something starts and ends?
		through offsets and size

		node - fixed-sized pages - sorted - 4kb to not waste any space to read/write from disk

		node - types
			internal
				pointers to other nodes
			leaf
				end node

	Stored in a concatenated manner
		Sequential approach for function call because it's stored in a concatenated manner

	Struct on disk:
		node
			[type, nkeys, pointers, offsets, key-values, unused]
			[header]

		key-values
			[key-size, value-size, key, value]

	Predetermined sizes:
		Header = 4 bytes
			type = 2 bytes
			nkeys = 2 bytes
		Pointer = 8 bytes - x64 system size
		Offset = 2 bytes
		Key-values
			...

	Rules:
		BNODE_HEADER + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE <= BTREE_PAGE_SIZE
*/

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*
	cli input
		->	turn into statement
			->	perform database action
					select
					insert
			->	table struct gets updated

	stmt
		type
		row
	table
		stores n pages
			stores n rows
				rows are stored through a concatenated byte array
				row[i] can be positioned through the offset
	row
		id
		username
		email
*/

/*
	TODO:
		Add logs to execution
		Implement B+tree
		Add transactions
		Save data to disk [durable]
			No in-place-updates
		Make it resistant to crashes [atomic]
			Incremental updates with logs
			Check with checksums
		Utilize vm to swap around indexes
*/

type StmtType string

const (
	Insert StmtType = "INSERT"
	Select StmtType = "SELECT"

	IdSize       = 4
	UsernameSize = 32
	EmailSize    = 255
	RowSize      = IdSize + UsernameSize + EmailSize

	IdOffset       = 0
	UsernameOffset = 4
	EmailOffset    = 32 + 4

	/*
		db shouldn't interact directly with disk
			db interacts with pages, they store regularly accessed data from disk
				page is 4kb
					vm is the one responsible of this layered transactions
	*/

	PageSize        = 4096 //4kb page size in most vm systems
	TableMaxPages   = 100
	RowsPerPage     = PageSize / RowSize
	RowsMaxPerTable = RowsPerPage * TableMaxPages
)

type Row struct {
	Id       uint32
	Username string
	Email    string
}

type Table struct {
	NumRows uint32
	Pages   [TableMaxPages][]byte
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

type Statement struct {
	t StmtType
	r Row
}

func Serialize(src Row, dst []byte) {
	binary.LittleEndian.PutUint32(dst[0:4], src.Id)
	copy(dst[4:36], src.Username)
	copy(dst[36:291], src.Email)
}

func Deserialize(src []byte, dst *Row) error {

	if len(src) < RowSize {
		return fmt.Errorf("buffer too small to contain a Row")
	}

	dst.Id = binary.LittleEndian.Uint32(src[0:4])
	dst.Username = string(bytes.Trim(src[4:36], "\x00"))
	dst.Email = string(bytes.Trim(src[36:291], "\x00"))

	return nil
}

func readCommand(r io.Reader) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		input := scanner.Text()
		if err := performCommand(input); err != nil {
			break
		}
	}
}

func performCommand(input string) error {
	if strings.HasPrefix(input, ".") {
		if err := doMetaCommand(input); err != nil {
			return err
		}

		return nil
	}

	stmt := &Statement{}

	if err := stmt.Prepare(input); err != nil {
		return err
	}

	var t *Table
	if err := stmt.Execute(t); err != nil {
		return err
	}

	return nil
}

func doMetaCommand(cmd string) error {
	switch cmd {
	case ".exit":
		os.Exit(0)

	case ".help":
		/*
			print a throughout list of all the features that the db implements
		*/
		return nil
	}
	return fmt.Errorf("error meta cmd not found")
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
		return stmt.execSelect(tbl)
	case Select:
		return stmt.execInsert(tbl)
	default:
		return fmt.Errorf("invalid statement action")
	}
}

func (stmt *Statement) execSelect(tbl *Table) error {
	if tbl.NumRows >= RowsMaxPerTable {
		return fmt.Errorf("table is full")
	}

	pageNumber := tbl.NumRows / RowsPerPage
	offset := RowSize * (tbl.NumRows % RowsPerPage)

	if tbl.Pages[pageNumber] == nil {
		tbl.Pages[pageNumber] = make([]byte, PageSize)
	}

	dst := tbl.Pages[pageNumber][offset : offset+RowSize]
	Serialize(stmt.r, dst)

	tbl.NumRows++

	return nil
}

func (stmt *Statement) execInsert(tbl *Table) error {
	src, err := tbl.GetRowByNum(int(stmt.r.Id))
	if err != nil {
		return err
	}
	var r Row
	if err := Deserialize(src, &r); err != nil {
		return err
	}

	fmt.Printf("Index: %d, Username: %s, Email: %s", r.Id, r.Username, r.Email)
	return nil
}
