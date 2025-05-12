// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package main

import (
	"fmt"
	"os"

	"github.com/mowen132/schemer/internal/core"
	"github.com/mowen132/schemer/internal/input"
)

func main() {
	e := &errorHandler{usage: true}

	cmd, err := input.ParseCommand(os.Args[1:])
	e.handle(err, "parsing command-line arguments")

	if cmd.Type == input.HelpCommand {
		fmt.Print(usage)
		return
	}

	cfg, err := input.LoadConfig()
	e.handle(err, "loading environment variables")

	e.usage = false

	conn, err := core.ConnectDatabase(cfg)
	e.handle(err, "connecting to database")
	defer conn.Close()

	e.conn = conn

	version, err := conn.LoadVersion()
	e.handle(err, "loading version")

	migrations, err := core.LoadMigrations(cfg.Migrations)
	e.handle(err, "loading migrations")

	if len(migrations) == 0 {
		fmt.Println("no migrations found")
		return
	}

	e.handle(migrations.Run(cmd, conn, version), "running migrations")
}

type errorHandler struct {
	usage bool
	conn  *core.Connection
}

func (c *errorHandler) handle(err error, context string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s: %v\n", context, err)

		if c.usage {
			fmt.Fprint(os.Stderr, usage)
		}

		if c.conn != nil {
			c.conn.Close()
		}

		os.Exit(1)
	}
}
