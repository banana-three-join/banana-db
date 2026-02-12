package main

import (
	"fmt"
	"os"
	"strings"
)

func readCommand(t *Table, input string) error {
	if strings.HasPrefix(input, ".") {
		if err := doMetaCommand(input); err != nil {
			return err
		}

		return nil
	}

	return doStatement(t, input)
}

func doStatement(t *Table, input string) error {
	stmt := &Statement{}

	output, err := stmt.Prepare(input)
	fmt.Println(output)

	if err != nil {
		return err
	}

	output, err = stmt.Execute(t)
	fmt.Println(output)
	if err != nil {
		return err
	}

	return nil
}

func doMetaCommand(cmd string) error {
	switch cmd {
	case ".exit":
		os.Exit(0)

	case ".help":
		/*
			print a throughout list of all the features that the db implements
		*/
		fmt.Println(`[META COMMANDS HELP]
.help
	Throughout list of all available commands
.exit
	Exits app

[STATEMENT HELP]
SELECT: RETURNS PAGE IF AVAIlABLE
	FORMAT: SELECT [PAGE INDEX]
INSERT: INSERTS VALUES INTO A CURRENT TABLE
	FORMAT: INSERT [INDEX] [USERNAME] [EMAIL]`)
		return nil
	}
	return fmt.Errorf("error meta cmd not found")
}
