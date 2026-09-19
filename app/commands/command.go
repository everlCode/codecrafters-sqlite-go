package commands

import "github.com/codecrafters-io/sqlite-starter-go/app/parser"

type Command interface {
	Run() []string
}

type Registry struct {
	Commands map[string]Command
}

func New(p *parser.Parser) *Registry {
	return &Registry{
		Commands: map[string]Command{
			"dbinfo": NewDbInfo(p),
			"tables": NewTables(p),
		},
	}
}
