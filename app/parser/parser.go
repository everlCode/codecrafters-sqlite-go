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

type FileHeader struct {
	PageSize int
}

type PageHeader struct {
	CellsCount int
}


func New(file *os.File) *Parser {
	return &Parser{
		file: file,
	}
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

func (p *Parser) ReadPageHeader() *PageHeader {
	header, err := p.GetBytes(8, 100)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	return &PageHeader{
		CellsCount: int(binary.BigEndian.Uint16(header[3:5])),
	}
}

func (p *Parser) ReadCellPointerArray(size int) []uint16 {
	bytess, err := p.GetBytes(size * 2, 108)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	var cellPointer uint16
	var cellPointerArray []uint16
	for i := 0;i < len(bytess);i += 2 {
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
			break;
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
