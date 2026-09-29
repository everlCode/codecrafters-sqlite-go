package commands

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/sqlite-starter-go/app/parser"
	"github.com/rqlite/sql"
)

type Select struct {
	parser *parser.Parser
}

func NewSelect(p *parser.Parser) *Select {
	return &Select{parser: p}
}

func (c *Select) Run(stm sql.Statement) []string {
	sel, ok := stm.(*sql.SelectStatement)
	if !ok {
		return []string{"Error!"}
	}

	var columns []string
	var isCount bool
	for _, col := range sel.Columns {
		switch expr := col.Expr.(type) {
			case *sql.Ident:
				columns = append(columns, expr.Name)
			case *sql.Call:
				if strings.EqualFold(expr.Name.Name, "count") && expr.Star.IsValid() {
					isCount = true
				}
			}
	}

	whereExpr, ok := sel.WhereExpr.(*sql.BinaryExpr)
	var whereCol string
	var whereVal string
	if whereExpr != nil {
		whereCol = whereExpr.X.(*sql.Ident).Name
		whereVal = whereExpr.Y.(*sql.StringLit).Value
	}
	
	schema := c.parser.ReadSqliteSchema()
	table, ok := sel.Source.(*sql.QualifiedTableName)
	if !ok {
		return []string{"unsupported FROM"} // JOIN, подзапрос и т.п.
	}
	tableSchema, ok := schema[table.TableName()]
	if !ok {
		return []string{"Table not found!"}
	}

	page := c.parser.ReadPage(tableSchema.Rootpage)

	var result []string
	var count int

	for _, rec := range page.Records {

		if whereExpr != nil {
			whereIndex := tableSchema.GetColumnIndex(whereCol)
			if rec.Values[whereIndex].String != whereVal {
				continue
			}
		}

		if isCount {
			count++
			continue
		}

		var row string
		for colIndex, col := range columns {
			values := rec.Values
			index := tableSchema.GetColumnIndex(col)
			if  index > 0 {
				row += values[index].String
				if colIndex + 1 < len(columns) {
					row += "|"
				}
			}
		}
		row += "\n"
		result = append(result, row)
	}

	if isCount {
		return []string{strconv.Itoa(count) + "\n"}
	}

	return result
}
