package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type Parser struct {
	file *os.File
}

type Page struct {
	PageHeader   *PageHeader
	Records []*Record
}

type FileHeader struct {
	PageSize int
}

type PageHeader struct {
	HeaderSize int
	CellsCount int
}

type TableSchema struct {
	Name     string
	Rootpage int
}

func New(file *os.File) *Parser {
	return &Parser{
		file: file,
	}
}

func (p *Parser) ReadPage(number int) *Page {
	fileHeader := p.ReadFileHeader()
	var offset int64
	var pageStart int64
	if number == 1 {
		offset = 100
		pageStart = 0
	} else {
		offset = int64(fileHeader.PageSize * (number - 1))
		pageStart = offset
	}
	page := Page{}
	page.PageHeader = p.ReadPageHeader(offset)
	cellsCount := page.PageHeader.CellsCount
	cells := p.ReadCellPointerArray(cellsCount, offset + int64(page.PageHeader.HeaderSize))

	var records []*Record
	for _, v := range cells {
		_, n1 := p.GetVarint(int64(v) + pageStart)
		_, n2 := p.GetVarint(int64(v) + pageStart)
		record := p.ReadRecord(int64(v) + int64(n1) + int64(n2) + pageStart)
		records = append(records, &record)
	}
	page.Records = records

	return &page
}

func (p *Parser) ReadSqliteSchema() map[string]*TableSchema {
	page := p.ReadPage(1)

	schema := map[string]*TableSchema{}
	for _,v := range page.Records {
		values := v.Values
		if values[0].String == TableType {
			tableSchema := TableSchema{}
			tableSchema.Name = values[2].String
			tableSchema.Rootpage = values[3].Integer
			schema[tableSchema.Name] = &tableSchema
		}
	}

	return schema
}

func (p *Parser) ReadFileHeader() *FileHeader {
	header, err := p.GetBytes(100, 0)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	var pageSize uint16
	if err := binary.Read(bytes.NewReader(header[16:18]), binary.BigEndian, &pageSize); err != nil {
		fmt.Println("Failed to read integer:", err)
		return nil
	}

	return &FileHeader{
		PageSize: int(pageSize),
	}
}

func (p *Parser) ReadPageHeader(offset int64) *PageHeader {
	pageHeader := PageHeader{}
	header, err := p.GetBytes(8, offset)
	b := header[0]
	if b == 2 || b == 10 {
		pageHeader.HeaderSize = 12
	} else {
		pageHeader.HeaderSize = 8
	}
	if err != nil {
		fmt.Errorf(err.Error())
	}

	pageHeader.CellsCount = int(binary.BigEndian.Uint16(header[3:5]))

	return &pageHeader
}

func (p *Parser) ReadCellPointerArray(size int, offset int64) []uint16 {
	bytess, err := p.GetBytes(size*2, offset)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	var cellPointer uint16
	var cellPointerArray []uint16
	for i := 0; i < len(bytess); i += 2 {
		binary.Read(bytes.NewReader(bytess[i:i+2]), binary.BigEndian, &cellPointer)
		cellPointerArray = append(cellPointerArray, cellPointer)
	}

	return cellPointerArray
}

func (p *Parser) GetBytes(n int, offset int64) ([]byte, error) {
	buf := make([]byte, n)
	if _, err := p.file.ReadAt(buf, offset); err != nil {
		return []byte{}, err
	}

	return buf, nil
}

func (p *Parser) GetVarint(offset int64) (uint64, int) {
	var result uint64
	var length int

	for i := 0; i < 9; i++ {
		length++
		b := make([]byte, 1)
		if _, err := p.file.ReadAt(b, offset); err != nil {
			break
		}
		if i == 8 {
			result = result<<8 | uint64(b[0])
		} else {
			result = result<<7 | uint64(b[0]&0x7F)
		}

		if (b[0] >> 7) == 0 {
			break
		}
		offset += 1
	}

	return result, length
}
