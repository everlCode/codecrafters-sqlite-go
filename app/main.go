package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codecrafters-io/sqlite-starter-go/app/commands"
	"github.com/codecrafters-io/sqlite-starter-go/app/dispatcher"
	"github.com/codecrafters-io/sqlite-starter-go/app/parser"
	// Available if you need it!
	// "github.com/xwb1989/sqlparser"
)

// Usage: your_program.sh sample.db .dbinfo
func main() {
	databaseFilePath := os.Args[1]
	commandName := os.Args[2]

	databaseFile, err := os.Open(databaseFilePath)
	if err != nil {
		log.Fatal(err)
	}

	registry := commands.New(parser.New(databaseFile))

	dispatcher := dispatcher.New(
		registry.Commands,
	)
	response := dispatcher.Dispatch(commandName)
	for _, v := range response {
		fmt.Println(v)
	}
}
