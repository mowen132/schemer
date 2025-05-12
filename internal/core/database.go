// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package core

import (
	"database/sql"
	"fmt"

	"github.com/mowen132/schemer/internal/drivers"
	"github.com/mowen132/schemer/internal/input"
)

type Connection struct {
	driver drivers.Driver
	db     *sql.DB
}

func ConnectDatabase(cfg *input.Config) (*Connection, error) {
	driver := drivers.Load(cfg.Driver)

	if driver == nil {
		return nil, fmt.Errorf("driver %q not found", cfg.Driver)
	}

	var dsn string

	if cfg.DSN != "" {
		dsn = cfg.DSN
	} else {
		dsn = driver.GetDSN(cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)
	}

	db, err := sql.Open(cfg.Driver, dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Connection{
		driver: driver,
		db:     db,
	}, nil
}

func (c *Connection) LoadVersion() (int64, error) {
	var version int64

	err := c.transact(func(tx *sql.Tx) error {
		val, err := c.driver.LoadVersion(tx)

		if err != nil {
			return err
		}

		version = val
		return nil
	})

	return version, err
}

func (c *Connection) saveVersion(tx *sql.Tx, version int64) error {
	return c.driver.SaveVersion(tx, version)
}

func (c *Connection) transact(fn func(tx *sql.Tx) error) error {
	tx, err := c.db.Begin()

	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (c *Connection) Close() error {
	return c.db.Close()
}
