package main

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
