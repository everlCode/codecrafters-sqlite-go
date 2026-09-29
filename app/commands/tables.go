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

	return []string{response}
}