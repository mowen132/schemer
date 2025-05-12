// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package main

const usage = `usage:

  schemer COMMAND [ARGUMENT]

commands:

  up [N]            Apply up to N migrations. If N is omitted, apply all remaining migrations.
  down [N]          Roll back up to N migrations. If N is omitted, roll back all applied migrations.
  goto VERSION      Migrate directly to the specified VERSION. Use 0 to roll back all migrations.
  help              Show this usage information.

examples:

  schemer up        # Apply all pending migrations
  schemer up 2      # Apply the next 2 migrations
  schemer down 1    # Roll back the most recent migration
  schemer goto 101  # Migrate to version 101

environment variables:

  DB_DRIVER         The database driver to use (e.g., postgres, mysql).
  DB_CONN           Optional. Full connection string (DSN). If set, overrides all other DB_* variables except DB_DRIVER.
                    Use this for advanced settings (e.g., SSL, timeouts) when basic DB_* values are not sufficient.
  DB_USER           Database username.
  DB_PASS           Database password.
  DB_HOST           Database hostname.
  DB_PORT           Database port. Optional — defaults to the standard port for the selected driver
                    (5432 for postgres, 3306 for mysql).
  DB_NAME           Database name.
  MIGRATIONS        Path to the directory containing migration files.

supported drivers:

  postgres          (PostgreSQL)
  mysql             (MySQL)

dsn documentation:

  PostgreSQL: https://pkg.go.dev/github.com/lib/pq#hdr-Connection_String_Parameters
  MySQL:      https://github.com/go-sql-driver/mysql#dsn-data-source-name

`
