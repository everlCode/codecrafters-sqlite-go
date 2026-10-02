package parser

import (
	"log"
	"strings"

	"github.com/rqlite/sql"
)

type DBSchema struct {
	Tables  map[string]*TableSchema
	Indexes map[string]*IndexSchema
}

func NewSchema() *DBSchema {
	return &DBSchema{
		Tables:  make(map[string]*TableSchema),
		Indexes: make(map[string]*IndexSchema),
	}
}

type IndexSchema struct {
	Name     string
	Table    string
	Rootpage int
}
type TableSchema struct {
	Name     string
	Columns  []string
	Rootpage int
}

func (p *Parser) ReadSqliteSchema() *DBSchema {
	pg := p.ReadPage(1)
	page, ok := pg.(*LeafPage)
	if !ok {
		log.Fatal("Неверная структура!")
	}

	schema := NewSchema()
	for _, v := range page.Records {
		values := v.Values

		switch values[0].String {
		case TableType:
			tableSchema := ReadTableSchema(values)
			schema.Tables[tableSchema.Name] = tableSchema
		case IndexType:
			indexSchema := ReadIndexSchema(values)
			schema.Indexes[indexSchema.Name] = indexSchema
		}
	}

	return schema
}

func (ts *TableSchema) GetColumnIndex(name string) int {
	var index int = -1
	for i, column := range ts.Columns {
		if column == name {
			index = i
			break
		}
	}

	return index
}

func ReadTableSchema(values []*Value) *TableSchema {
	tableSchema := TableSchema{}
	tableSchema.Name = values[2].String
	tableSchema.Rootpage = int(values[3].Integer)

	parser := sql.NewParser(strings.NewReader(values[4].String))
	stm, err := parser.ParseStatement()
	if err != nil {
		log.Fatal(err)
	}
	createTable, ok := stm.(*sql.CreateTableStatement)
	if !ok {
		return nil
	}
	var columns []string
	for _, v := range createTable.Columns {
		columns = append(columns, v.Name.Name)
	}
	tableSchema.Columns = columns

	return &tableSchema
}

func ReadIndexSchema(values []*Value) *IndexSchema {
	tableSchema := IndexSchema{}
	tableSchema.Name = values[1].String
	tableSchema.Table = values[2].String
	tableSchema.Rootpage = int(values[3].Integer)

	parser := sql.NewParser(strings.NewReader(values[4].String))
	stm, err := parser.ParseStatement()
	if err != nil {
		log.Fatal(err)
	}
	_, ok := stm.(*sql.CreateIndexStatement)
	if !ok {
		return nil
	}
	// var columns []string
	// for _,v := range createTable.Columns {
	// 	columns = append(columns, v.X.String())
	// }
	// tableSchema.Columns = columns

	return &tableSchema
}
