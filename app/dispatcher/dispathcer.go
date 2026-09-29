package dispatcher

import (
	"strings"

	"github.com/codecrafters-io/sqlite-starter-go/app/commands"
	"github.com/rqlite/sql"
)

type dispatcher struct {
	Registry *commands.Registry
}

func New(registry *commands.Registry) *dispatcher {
	return &dispatcher{
		Registry: registry,
	}
}

func (dispatcher *dispatcher) Dispatch(input string) []string {
	if input[0] == '.' {
		input = input[1:]
		cmd, ok := dispatcher.Registry.DotCommands[input]
		if !ok {
			return []string{"Command doesnt exist!"}
		}
		return cmd.Run()
	}
	parser := sql.NewParser(strings.NewReader(input))
	stm, err := parser.ParseStatement()
	if err != nil {
		return []string{"Неверный запрос: " + err.Error()}
	}
	switch stm.(type) {
	case *sql.SelectStatement:
		cmd, ok := dispatcher.Registry.SqlCommands["select"]
		if !ok {
			return []string{"Command doesnt exist!"}
		}
		return cmd.Run(stm)
	}

	return []string{}
}
