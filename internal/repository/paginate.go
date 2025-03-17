package repository

import "strconv"

type Paginate struct {
	Page string
	Size string

	page uint
	size uint
}

// Default page 1
// Default page size 20
func (p *Paginate) parsePaginate() {
	page, err := strconv.ParseUint(p.Page, 10, 32)
	if err != nil {
		p.page = 1
	} else {
		p.page = uint(page)
	}

	size, err := strconv.ParseUint(p.Size, 10, 32)
	if err != nil {
		p.size = 25
	} else {
		p.size = uint(size)
	}
}
