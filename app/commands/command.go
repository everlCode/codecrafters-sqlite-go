package commands

import (
	"github.com/codecrafters-io/sqlite-starter-go/app/parser"
	"github.com/xwb1989/sqlparser"
)

type DotCommand interface {
	Run() []string
}

type SqlCommand interface {
	Run(sqlparser.Statement) []string
}

type Registry struct {
	DotCommands map[string]DotCommand
	SqlCommands map[string]SqlCommand
}

func New(p *parser.Parser) *Registry {
	return &Registry{
		DotCommands: map[string]DotCommand{
			"dbinfo": NewDbInfo(p),
			"tables": NewTables(p),
		},
		SqlCommands: map[string]SqlCommand{
			"select": NewSelect(p),
		},
	}
}
