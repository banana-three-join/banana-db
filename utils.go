package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

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
