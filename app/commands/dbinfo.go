package commands

import (
	"fmt"

	"github.com/codecrafters-io/sqlite-starter-go/app/parser"
)

type DbInfo struct {
	parser *parser.Parser
}

func NewDbInfo(p *parser.Parser) *DbInfo {
	return &DbInfo{parser: p}
}

func (c *DbInfo) Run() []string {
	fileHeader := c.parser.ReadFileHeader()
	pageHeader := c.parser.ReadPageHeader()

	return []string{
		fmt.Sprintf("database page size: %v\n", fileHeader.PageSize),
		fmt.Sprintf("number of tables: %v\n", pageHeader.CellsCount),
	}
}
