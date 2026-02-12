package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func Serialize(src Row, dst []byte) {
	binary.LittleEndian.PutUint32(dst[:IdSize], src.Id)
	copy(dst[UsernameOffset:UsernameOffset+UsernameSize], src.Username)
	copy(dst[EmailOffset:EmailOffset+EmailSize], src.Email)
}

// deserializes only one row. maybe implement a deserialize for a full page?
func Deserialize(src []byte, dst *Row) error {

	if len(src) < RowSize {
		return fmt.Errorf("buffer too small to contain a Row")
	}

	dst.Id = binary.LittleEndian.Uint32(src[:IdSize])
	dst.Username = string(bytes.Trim(src[UsernameOffset:UsernameOffset+UsernameSize], "\x00"))
	dst.Email = string(bytes.Trim(src[EmailOffset:EmailOffset+EmailSize], "\x00"))

	return nil
}

func SplitRowsFromPage(p []byte) [][]byte {
	pageLen := len(p)
	if pageLen == 0 {
		return nil
	}

	numRows := pageLen / RowSize
	if pageLen%RowSize != 0 {
		numRows++
	}

	rows := make([][]byte, numRows)

	for i := range numRows {
		offset := RowSize * i
		row := make([]byte, RowSize)
		copy(row, p[offset:offset+RowSize])
		rows[i] = row
	}

	return rows
}
