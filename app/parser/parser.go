package parser

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type Parser struct {
	reader io.Reader
}

type PageHeader struct {
	CellsCount int
}

func New(file *os.File) *Parser {
	reader := bufio.NewReader(file)

	return &Parser{
		reader: reader,
	}
}

func (p *Parser) ParsePageHeader() PageHeader {
	header := p.GetBytes(8)

	return PageHeader{
		CellsCount: int(binary.BigEndian.Uint16(header[3:5])),
	}
}

func (p *Parser) GetBytes(n int) []byte {
	b := make([]byte, n)
	_, err := p.reader.Read(b)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	return b
}