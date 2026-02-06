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

package main

import (
	"bufio"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := readCommand(scanner.Text()); err != nil {
			break
		}
	}
}
