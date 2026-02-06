package main

import (
	"fmt"
	"os"
	"strings"
)

func readCommand(input string) error {
	if err := performCommand(input); err != nil {
		return err
	}
	return nil
}

func performCommand(input string) error {
	if strings.HasPrefix(input, ".") {
		if err := doMetaCommand(input); err != nil {
			return err
		}

		return nil
	}

	stmt := &Statement{}

	if err := stmt.Prepare(input); err != nil {
		return err
	}

	t := &Table{}
	if err := stmt.Execute(t); err != nil {
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
		return nil
	}
	return fmt.Errorf("error meta cmd not found")
}
