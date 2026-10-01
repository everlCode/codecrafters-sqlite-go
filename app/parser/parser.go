package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rqlite/sql"
)

type Parser struct {
	file *os.File
}

type FileHeader struct {
	PageSize int
}

type PageHeader struct {
	HeaderSize int
	CellsCount int
	PageType string
	RightMostPointer uint32
}

const (
	LEAF_PAGE_TYPE = "leaf"
	INTERIOR_PAGE_TYPE = "interior"
)


func New(file *os.File) *Parser {
	return &Parser{
		file: file,
	}
}


func (p *Parser) GetRecords(number int) []*Record {
	page := p.ReadPage(number)
	ip, ok := page.(*InteriorPage)
	var records []*Record
	if ok {
		for _, cell := range ip.Cells {
			leftRecords := p.GetRecords(int(cell.LeftChildPageNumber))
			records = append(records, leftRecords...)
		}

		rightRecords := p.GetRecords(int(ip.PageHeader.RightMostPointer))
		

		return append(records, rightRecords...)
	}

	lp, ok := page.(*LeafPage)
	if !ok {
		log.Fatal("Некоректный тип таблицы")
	}
	
	return lp.Records
}


func (p *Parser) ReadPage(number int) Page {
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
	pageHeader := p.ReadPageHeader(offset)
	cellsPointerArray := p.ReadCellPointerArray(pageHeader.CellsCount, offset + int64(pageHeader.HeaderSize))

	if pageHeader.PageType == LEAF_PAGE_TYPE {
		return p.ReadLeafPageContent(pageStart, pageHeader, cellsPointerArray)
	}

	return p.ReadInteriorPageContent(pageStart, pageHeader, cellsPointerArray)
}

func (p *Parser) ReadLeafPageContent(pageStart int64, pageHeader *PageHeader, cellPonterArr []uint16) *LeafPage {
	page := LeafPage{PageStart: pageStart, PageHeader: pageHeader}
	var records []*Record
	var cells []*LeafCell
	for _, offset := range cellPonterArr {
		cell := LeafCell{}
		payloadSize, n1 := p.GetVarint(int64(offset) + pageStart)
		cell.PayloadSize = payloadSize
		rowId, n2 := p.GetVarint(int64(offset) + int64(n1) + pageStart)
		cell.RowId = int(rowId)
		record := p.ReadRecord(int64(offset) + int64(n1) + int64(n2) + pageStart)
		record.RowId = int(rowId)
		cells = append(cells, &cell)
		records = append(records, &record)
	}
	page.Cells = cells
	page.Records = records

	return &page
}

func (p *Parser) ReadInteriorPageContent(pageStart int64, pageHeader *PageHeader, cellPonterArr []uint16) *InteriorPage {
	var cells []*InteriorCell
	for _, offset := range cellPonterArr {
		cell := InteriorCell{}
		bytes, err := p.GetBytes(4, int64(offset) + pageStart)
		if err != nil {
			log.Fatal(err.Error())
		}
		cell.LeftChildPageNumber = binary.BigEndian.Uint32(bytes)
		rowID, _ := p.GetVarint(pageStart + int64(offset) + 4)
		cell.RowId = rowID
		cells = append(cells, &cell)
	}
	page := InteriorPage{PageHeader: pageHeader}
	page.Cells = cells

	return &page
}

func (p *Parser) ReadSqliteSchema() map[string]*TableSchema {
	pg := p.ReadPage(1)
	page, ok := pg.(*LeafPage)
	if !ok {
		log.Fatal("Неверная структура!")
	}

	schema := map[string]*TableSchema{}
	for _,v := range page.Records {
		values := v.Values
		if values[0].String == TableType {
			tableSchema := TableSchema{}
			tableSchema.Name = values[2].String
			tableSchema.Rootpage = values[3].Integer
			

			parser := sql.NewParser(strings.NewReader(values[4].String))
			stm, err := parser.ParseStatement()
			if err != nil {
				log.Fatal(err)
			}
			createTable, ok := stm.(*sql.CreateTableStatement)
			var columns []string
			for _,v := range createTable.Columns {
				columns = append(columns, v.Name.Name)
			}
			tableSchema.Columns = columns
			if !ok {
				continue // это не CREATE TABLE
			}
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
	if err != nil {
		fmt.Errorf(err.Error())
	}
	b := header[0]
	if b == 2 || b == 5 {
		header, err := p.GetBytes(12, offset)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		pageHeader.HeaderSize = 12
		pageHeader.PageType = INTERIOR_PAGE_TYPE
		pageHeader.RightMostPointer = binary.BigEndian.Uint32(header[8:12])
	} else {
		pageHeader.HeaderSize = 8
		pageHeader.PageType = LEAF_PAGE_TYPE
	}
	

	pageHeader.CellsCount = int(binary.BigEndian.Uint16(header[3:5]))

	return &pageHeader
}

func (p *Parser) ReadCellPointerArray(cellsCount int, offset int64) []uint16 {
	bytess, err := p.GetBytes(cellsCount*2, offset)
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
