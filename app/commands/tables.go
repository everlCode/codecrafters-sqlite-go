package commands

import (
	"sort"
	"strings"

	"github.com/codecrafters-io/sqlite-starter-go/app/parser"
)

type Tables struct {
	parser *parser.Parser
}

func NewTables(p *parser.Parser) *Tables {
	return &Tables{parser: p}
}

func (c *Tables) Run()[]string {
	schema := c.parser.ReadSqliteSchema()

	var names []string
	for _, table := range schema {
		names = append(names, table.Name)
	}
	sort.Strings(names)
	response := strings.Join(names, " ")

	// c.parser.ReadFileHeader()
	// pageHeader := c.parser.ReadPageHeader(100)
	// cellsCount := pageHeader.CellsCount
	// cells := c.parser.ReadCellPointerArray(cellsCount)
	
	// var tableNames []string
	// for _, v := range cells {
	// 	_, n1 := c.parser.GetVarint(int64(v))
	// 	_, n2 := c.parser.GetVarint(int64(v))
	// 	record := c.parser.ReadRecord(int64(v) + int64(n1) + int64(n2))
	// 	values := record.Values
	// 	if values[0].String == parser.TableType {
	// 		tableNames = append(tableNames, values[2].String)
	// 	}

	// }

	return []string{response}
}