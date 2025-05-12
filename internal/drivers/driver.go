// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package drivers

import (
	"database/sql"
)

type Driver interface {
	GetDSN(user string, pass string, host string, port int, name string) string
	LoadVersion(tx *sql.Tx) (int64, error)
	SaveVersion(tx *sql.Tx, version int64) error
}

func Load(name string) Driver {
	switch name {
	case "postgres":
		return &postgresDriver{}

	case "mysql":
		return &mysqlDriver{}
	}

	return nil
}
