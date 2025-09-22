// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package input

import (
	"fmt"
	"strconv"
)

type Command struct {
	Type       CommandType
	HasOperand bool
	Operand    int64
}

type CommandType int

const (
	CommandUp CommandType = iota
	CommandDown
	CommandGoto
	CommandHelp
)

func ParseCommand(args []string) (*Command, error) {
	n := len(args)

	if n == 0 {
		return nil, fmt.Errorf("no command was specified")
	}

	name := args[0]

	switch name {
	case "up", "down":
		var t CommandType

		if name == "up" {
			t = CommandUp
		} else {
			t = CommandDown
		}

		if n == 1 {
			return &Command{
				Type: t,
			}, nil
		}

		num, err := parseNumber(args[1])

		if err != nil {
			return nil, err
		}

		return &Command{
			Type:       t,
			HasOperand: true,
			Operand:    num,
		}, nil

	case "goto":
		if n < 2 {
			return nil, fmt.Errorf("no target was specified")
		}

		tar, err := parseTarget(args[1])

		if err != nil {
			return nil, err
		}

		return &Command{
			Type:       CommandGoto,
			HasOperand: true,
			Operand:    tar,
		}, nil

	case "help":
		return &Command{
			Type: CommandHelp,
		}, nil
	}

	return nil, fmt.Errorf("%q is not a valid command", name)
}

func parseNumber(arg string) (int64, error) {
	num, err := strconv.ParseInt(arg, 10, 64)

	if err != nil {
		return 0, fmt.Errorf("%q is not a valid number of migrations", arg)
	}

	if num <= 0 {
		return 0, fmt.Errorf("%q is not a postive number", arg)
	}

	return num, nil
}

func parseTarget(arg string) (int64, error) {
	tar, err := strconv.ParseInt(arg, 10, 64)

	if err != nil || tar < 0 {
		return 0, fmt.Errorf("%q is not a valid target", arg)
	}

	return tar, nil
}
