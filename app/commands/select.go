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
	tableSchema, ok := schema.Tables[table.TableName()]
	if !ok {
		return []string{"Table not found!"}
	}

	index:=  getIndex(table.TableName(), schema.Indexes) 
	if index != nil {
		recs := c.GetRecordsByIndex(index.Rootpage, whereVal)

		var records []*parser.Record
		for _, rec := range recs {
			record := c.getRecordByRowId(int(rec.Values[1].Integer), tableSchema.Rootpage)
			records = append(records, record)
		}


		//return []string{strconv.Itoa(len(records))}

		return formatRecords(records, tableSchema, columns, "", "", isCount)
	}

	records := c.parser.GetRecords(tableSchema.Rootpage)

	return formatRecords(records, tableSchema, columns, whereCol, whereVal, isCount)
}

func formatRecords(records []*parser.Record, tableSchema *parser.TableSchema, columns []string, whereCol string, whereVal string, isCount bool) []string {
	whereIndex := -1
	if whereCol != "" {
		whereIndex = tableSchema.GetColumnIndex(whereCol)
	}

	var result []string
	var count int
	for _, rec := range records {
		if whereIndex >= 0 && rec.Values[whereIndex].String != whereVal {
			continue
		}

		if isCount {
			count++
			continue
		}

		var fields []string
		for _, col := range columns {
			index := tableSchema.GetColumnIndex(col)
			if index < 0 {
				continue
			}
			if col == "id" {
				fields = append(fields, strconv.Itoa(rec.RowId))
			} else {
				fields = append(fields, rec.Values[index].String)
			}
		}
		result = append(result, strings.Join(fields, "|")+"\n")
	}

	if isCount {
		return []string{strconv.Itoa(count) + "\n"}
	}

	return result
}

func getIndex(table string, indexes map[string]*parser.IndexSchema) *parser.IndexSchema {
	for _, index := range indexes {
		if index.Table == table {
			return index
		}
	}
	return nil
}

func (c *Select) getRecordByRowId(rowId int, rootPage int) *parser.Record {
	page := c.parser.ReadPage(rootPage)
	switch p := page.(type) {
	case *parser.InteriorPage:
		for _, cell := range p.Cells {
			ce, ok := cell.(*parser.InteriorTableCell)
			if !ok {
				continue
			}

			if rowId <= int(ce.RowId) {
				return c.getRecordByRowId(rowId, int(ce.LeftChildPageNumber))
			}
		}

	
		return c.getRecordByRowId(rowId, int(p.PageHeader.RightMostPointer))
	case *parser.LeafPage:
		for _, rec := range p.Records {
			if rec.RowId == rowId {
				return rec
			}
		}

	}

	return &parser.Record{}
}

func (c *Select) GetRecordsByIndex(rootapgeNumber int, whereVal string) []*parser.Record {
	records := []*parser.Record{}
	var compareResult int
	switch p := c.parser.ReadPage(rootapgeNumber).(type) {
	case *parser.LeafPage:
		for _, ce := range p.Cells {
			cell, ok := ce.(*parser.LeaIndexfCell)
			if !ok {
				continue
			}
			compareResult = comareValue(cell.Payload.Values[0].String, whereVal)
			if compareResult == 0 {
				records = append(records, cell.Payload)
			}

			if compareResult < 0 {
				break
			}
			
		}
	case *parser.InteriorPage:
		for _, ce := range p.Cells {
			cell, ok := ce.(*parser.InteriorIndexCell)
			if !ok {
				continue
			}
			compareResult = comareValue(cell.Payload.Values[0].String, whereVal)
			if compareResult == 0 {
				records = append(records, cell.Payload)
			}

			if compareResult < 0 {
				rs := c.GetRecordsByIndex(int(cell.LeftChildPageNumber), whereVal)
				records = append(records, rs...)
				break
			}
		}

		if compareResult > 0 {
			rs := c.GetRecordsByIndex(int(p.PageHeader.RightMostPointer), whereVal)
			records = append(records, rs...)
		}
		
		
	}

	return records
}

func comareValue(value string, condition string) int {
	if value == condition {
		return 0
	}

	if condition < value {
		return -1
	}

	return 1
}
