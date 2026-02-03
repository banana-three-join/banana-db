/*
	DB Theory:
		B+tree - used to store key-values that require fast indexing through a quick "merge" approach
		everything is stored through concatenation. how can we tell were something starts and ends?
		through offsets and size

		node - fixed-sized pages - sorted - 4kb to not waste any space to read/write from disk

		node - types
			internal - pointers to other nodes
			leaf - end node

	Stored in a concatenated manner
		Sequential approach for function call because it's stored in a concatenated manner

	Struct on disk:
		node
			[type, nkeys, pointers, offsets, key-values, unused]

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
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"os"
)

const (
	BNODE_NODE = 1
	BNODE_LEAF = 2

	BTREE_PAGE_SIZE    = 4096 //node
	BTREE_MAX_KEY_SIZE = 1000
	BTREE_MAX_VAL_SIZE = 3000
)

/*tmp files assure that the only files that will be created, are the files that complete their execution*/
func SaveData(path string, data []byte) error {
	tmp := fmt.Sprintf("%s.tmp.%d", path, rand.Uint())
	fp, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
	if err != nil {
		return err
	}
	defer func() {
		fp.Close()
		if err != nil {
			os.Remove(tmp)
		}
	}()

	if _, err := fp.Write(data); err != nil {
		return err
	}

	if err := fp.Sync(); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}

type BNode []byte

// type is represented to uint16 - 2 bytes
func (bn BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(bn[0:2])
}

func (bn BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(bn[2:4])
}

func (bn BNode) setHeader(t, ks uint16) {
	binary.LittleEndian.PutUint16(bn[0:2], t)
	binary.LittleEndian.PutUint16(bn[2:4], ks)
}

// ptr uint64 - 8bytes
func (bn BNode) getPtr(idx uint16) uint64 {
	pos := 4 + 8*idx
	return binary.LittleEndian.Uint64(bn[pos:])
}

func (bn BNode) setPtr(idx uint16, ptr uint64) {
	pos := 4 + 8*idx
	binary.LittleEndian.PutUint64(bn[pos:], ptr)
}

func (bn BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	pos := 4 + 8*bn.nkeys() + 2*(idx-1)
	return binary.LittleEndian.Uint16(bn[pos:])
}

func (bn BNode) kvPos(idx uint16) uint16 {
	return 4 + 8*bn.nkeys() + 2*bn.nkeys() + bn.getOffset(idx)
}

func (bn BNode) getKey(idx uint16) []byte {
	pos := bn.kvPos(idx)
	klen := binary.LittleEndian.Uint16(bn[pos:])
	return bn[pos+4:][:klen]
}

func (node BNode) getVal(idx uint16) []byte {
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos+0:])
	vlen := binary.LittleEndian.Uint16(node[pos+2:])
	return node[pos+4+klen:][:vlen]
}
