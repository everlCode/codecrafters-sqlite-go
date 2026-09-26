package dispatcher

import (
	"github.com/codecrafters-io/sqlite-starter-go/app/commands"
	"github.com/xwb1989/sqlparser"
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
	stm, _ := sqlparser.Parse(input)
	switch stm.(type) {
	case *sqlparser.Select:
		cmd, ok := dispatcher.Registry.SqlCommands["select"]
		if !ok {
			return []string{"Command doesnt exist!"}
		}
		return cmd.Run(stm)
	}

	return []string{}
}
