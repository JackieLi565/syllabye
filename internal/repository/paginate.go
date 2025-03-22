package repository

import "strconv"

type Paginate struct {
	Page uint
	Size uint
}

// Default page 1
// Default page size 20
func NewPaginate(pageStr, sizeStr string) Paginate {
	page, err := strconv.ParseUint(pageStr, 10, 32)
	if err != nil || page == 0 {
		page = 1
	}

	size, err := strconv.ParseUint(sizeStr, 10, 32)
	if err != nil || size == 0 {
		size = 25
	}

	return Paginate{
		Page: uint(page),
		Size: uint(size),
	}
}
