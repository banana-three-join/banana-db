package main

import (
	"fmt"
	"os"
)

// works as a cached layer from the file
// get requests and if the page looked after doesn't exists then it goes into the file to write
type Pager struct {
	Filename string
	Pages    [MaxPagesPerTable][]byte
}

func (p *Pager) SetRow(numRows uint32, row Row) error {
	//should write into db instead and implement lifo with the table
	if numRows >= MaxRowsPerTable {
		return fmt.Errorf("table is currently full")
	}

	pageNumber := numRows / RowsPerPage
	rowNumber := numRows % RowsPerPage
	rowOffset := RowSize * rowNumber

	if p.Pages[pageNumber] == nil {
		p.Pages[pageNumber] = make([]byte, PageSize)
	}

	dst := p.Pages[pageNumber][rowOffset : rowOffset+RowSize]
	Serialize(row, dst)

	fd, err := os.OpenFile(p.Filename, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	defer fd.Close()

	offset := pageNumber*PageSize + rowOffset
	_, err = fd.Seek(int64(offset), 0)
	if err != nil {
		return err
	}

	n, err := fd.Write(dst)
	if err != nil {
		return err
	}

	if n != RowSize {
		return fmt.Errorf("row wasn't fully written onto database")
	}

	return nil
}

func (p *Pager) GetPage(pageIndex int) ([]byte, error) {
	if pageIndex < 0 || pageIndex >= MaxPagesPerTable {
		return nil, fmt.Errorf("fetched page %d is out of bound", pageIndex)
	}

	if p.Pages[pageIndex] == nil {
		fp, err := os.Open(p.Filename)
		if err != nil {
			return nil, err
		}
		defer fp.Close()

		st, err := fp.Stat()
		if err != nil {
			return nil, err
		}
		fileLen := st.Size()

		if fileLen == 0 {
			return nil, fmt.Errorf("fetched page %d is empty", pageIndex)
		}

		//get number of full pages in the current file
		numPages := fileLen / PageSize

		//check if there's a partial page
		if (fileLen % PageSize) != 0 {
			numPages++
		}

		//if the asked index is lesser or equal to the amount of available pages
		//then we can go ahead and read until the available number of pages
		if pageIndex < int(numPages) {
			//read from pages from file to cache onto pager
			for i := range numPages {
				writtenOntoPage := len(p.Pages[i])

				//page isn't empty and is full
				if p.Pages[i] != nil && (writtenOntoPage >= PageSize) {
					continue
				}

				fileOffset := i * PageSize

				//page isn't empty and isn't full
				if p.Pages[i] != nil && (writtenOntoPage < PageSize) {
					pageSpaceRemaining := PageSize - writtenOntoPage
					buffer := make([]byte, pageSpaceRemaining)
					fp.ReadAt(buffer, (int64(writtenOntoPage) + fileOffset))
					p.Pages[i] = append(p.Pages[i], buffer...)
					continue
				}

				//page is empty
				page := make([]byte, PageSize)
				//check if the chunk of the file wasn't fully written
				_, err := fp.ReadAt(page, fileOffset)
				if err != nil {
					return nil, err
				}

				p.Pages[i] = page
			}
		}

		if p.Pages[pageIndex] == nil {
			return nil, fmt.Errorf("file doesn't contain enough pages for index %d", pageIndex)
		}
	}

	return p.Pages[pageIndex], nil
}

func NewPager(filename string) (*Pager, error) {
	var pages [MaxPagesPerTable][]byte
	return &Pager{
		Filename: filename,
		Pages:    pages,
	}, nil
}

func (p *Pager) FlushPages() {
	var pages [MaxPagesPerTable][]byte
	p.Pages = pages
}
