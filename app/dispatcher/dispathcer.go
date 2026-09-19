package dispatcher

import "github.com/codecrafters-io/sqlite-starter-go/app/commands"

type dispatcher struct {
	commands map[string]commands.Command
}

func New(commands map[string]commands.Command) *dispatcher {
	return &dispatcher{
		commands: commands,
	}
}

func (dispatcher *dispatcher) Dispatch(command string) []string {
	if (command[0] == '.') {
		command = command[1:]
	}
	cmd, ok := dispatcher.commands[command]
	if !ok {
		return []string{"Command doesnt exist!"}
	}

	return cmd.Run()
}
