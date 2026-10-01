package parser

type LeafPage struct {
	PageStart int64
	PageHeader   *PageHeader
	CellPointerArray []int16
	Cells []*LeafCell
	Records []*Record
}

type InteriorPage struct {
	PageHeader   *PageHeader
	Cells []*InteriorCell
}

type Page interface {

}

type LeafCell struct {
	PayloadSize uint64
	RowId int
	PayloadStart int32
	FirstOverflowPage int32
}

type InteriorCell struct {
	LeftChildPageNumber uint32
	RowId uint64
}