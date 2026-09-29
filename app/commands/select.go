package commands

import (
	"strconv"
	"strings"

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
	
	var columns []string
	for _, selectExpr := range sel.SelectExprs {
		s := sqlparser.String(selectExpr)
		columns = append(columns, s)
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

	var result []string
	
	for _, col := range columns {
		for _, rec := range page.Records {
			values := rec.Values
			if strings.ToLower(col) == "count(*)" {
				result = append(result, strconv.Itoa(page.PageHeader.CellsCount) + "\n")
				break
			}
			index := tableSchema.GetColumnIndex(col)
			if  index > 0 {
				result = append(result, values[index].String + "\n")
			} 
		}
	} 

	return result
}
