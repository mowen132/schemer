// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package drivers

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlDriver struct{}

func (d *mysqlDriver) GetDSN(user string, pass string, host string, port int, name string) string {
	if port == -1 {
		port = 3306
	}

	q := url.Values{}
	q.Set("parseTime", "true")
	q.Set("charset", "utf8mb4")

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		url.QueryEscape(user),
		url.QueryEscape(pass),
		host,
		port,
		url.QueryEscape(name),
		q.Encode(),
	)
}

func (d *mysqlDriver) LoadVersion(tx *sql.Tx) (int64, error) {
	if _, err := tx.Exec("CREATE TABLE IF NOT EXISTS schemer (version BIGINT NOT NULL)"); err != nil {
		return 0, fmt.Errorf("failed to ensure schemer table exists: %v", err)
	}

	rows, err := tx.Query("SELECT version FROM schemer")

	if err != nil {
		return 0, fmt.Errorf("failed to query schemer version: %v", err)
	}

	defer rows.Close()

	if !rows.Next() {
		if _, err := tx.Exec("INSERT INTO schemer (version) VALUES (0)"); err != nil {
			return 0, fmt.Errorf("failed to initialize schemer version to 0: %v", err)
		}

		return 0, nil
	}

	var version int64

	if err := rows.Scan(&version); err != nil {
		return 0, fmt.Errorf("failed to scan schemer version: %v", err)
	}

	return version, nil
}

func (d *mysqlDriver) SaveVersion(tx *sql.Tx, version int64) error {
	if _, err := tx.Exec("UPDATE schemer SET version = ?", version); err != nil {
		return fmt.Errorf("failed to update schemer version to %d: %v", version, err)
	}

	return nil
}
