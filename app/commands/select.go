package commands

import (
	"strconv"

	"github.com/codecrafters-io/sqlite-starter-go/app/parser"
	"github.com/xwb1989/sqlparser"
)

type Select struct {
	parser *parser.Parser
}

func NewSelect(p *parser.Parser) *Select {
	return &Select{parser: p}
}

func (c *Select) Run(stm sqlparser.Statement) []string {
	sel, ok := stm.(*sqlparser.Select)
	if !ok {
		return []string{"Error!"}
	}
	
	schema := c.parser.ReadSqliteSchema()
	tableExpr, ok := sel.From[0].(*sqlparser.AliasedTableExpr)
	if !ok {
		return []string{"unsupported FROM expression"}
	}
	tableName := sqlparser.GetTableName(tableExpr.Expr).String()
	tableSchema, ok := schema[tableName]
	if !ok {
		return []string{"Table not found!"} 
	}

	page := c.parser.ReadPage(tableSchema.Rootpage)

	return []string{strconv.Itoa(page.PageHeader.CellsCount) + "\n"}
}
