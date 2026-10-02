package parser

type LeafPage struct {
	PageStart        int64
	PageHeader       *PageHeader
	CellPointerArray []int16
	Cells            []Cell
	Records          []*Record
}

type InteriorPage struct {
	PageHeader *PageHeader
	Cells      []Cell
}

type Page interface {
}

type Cell interface {
}

type LeafTableCell struct {
	PayloadSize       uint64
	RowId             int
	PayloadStart      int32
	FirstOverflowPage int32
}

type InteriorTableCell struct {
	LeftChildPageNumber uint32
	RowId               uint64
}

type InteriorIndexCell struct {
	LeftChildPageNumber uint32
	PayloadSize         uint64
	Payload             *Record
	OverflowPage        int32
}

type LeaIndexfCell struct {
	PayloadSize         uint64
	Payload             *Record
	OverflowPage        int32
}
