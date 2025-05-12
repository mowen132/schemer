# schemer

**schemer** is a lightweight command-line tool for managing database schema changes. It supports both PostgreSQL and MySQL, and uses versioned migrations with clear `up` and `down` sections defined in a single SQL file.

---

## Features

- Simple CLI interface.
- Environment variable configuration.
- Sequential versioned migrations using numeric prefixes.
- Single-file migrations with `-- @up` and `-- @down` sections.
- Supports PostgreSQL and MySQL databases.
- Docker and Task integration.

---

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.24 or higher)
- [Task](https://taskfile.dev/#/installation)
- [Docker](https://www.docker.com/get-started)

### Build from Source

1. Clone the repository:

  ```bash
  git clone https://github.com/mowen132/schemer.git
  cd schemer
  ```

2. Build the binary (`./bin/schemer`):

  ```bash
  task build
  ```

3. (Optional) Install the binary system-wide (`/usr/bin/schemer`):

  ```bash
  task install
  ```

---

## Usage

### Command Syntax

  ```bash
  schemer COMMAND [ARGUMENT]
  ```

### Commands

- `up [N]` Apply up to N migrations. If N is omitted, apply all remaining migrations.
- `down [N]` Roll back up to N migrations. If N is omitted, roll back all applied migrations.
- `goto VERSION` Migrate directly to the specified VERSION. Use 0 to roll back all migrations.
- `help` Show usage information.

### Examples

  ```bash
  schemer up        # Apply all pending migrations
  schemer up 2      # Apply the next 2 migrations
  schemer down 1    # Roll back the most recent migration
  schemer goto 101  # Migrate to version 101
  ```

## Configuration

### Environment Variables:

Configure `schemer` using the following environment variables:

- `DB_DRIVER` The database driver to use (e.g., postgres, mysql).
- `DB_CONN` Optional. Full connection string (DSN). If set, overrides all other DB_* variables except DB_DRIVER. Use this for advanced settings (e.g., SSL, timeouts) when basic DB_* values are not sufficient.
- `DB_USER` Database username.
- `DB_PASS` Database password.
- `DB_HOST` Database hostname.
- `DB_PORT` Database port. Optional — defaults to the standard port for the selected driver (5432 for postgres, 3306 for mysql).
- `DB_NAME` Database name.
- `MIGRATIONS` Path to the directory containing migration files.

### Supported Drivers:

- postgres (PostgreSQL)
- mysql (MySQL)

### DSN Documentation:

- [PostgreSQL](https://pkg.go.dev/github.com/lib/pq#hdr-Connection_String_Parameters)
- [MySQL](https://github.com/go-sql-driver/mysql#dsn-data-source-name)

### Examples

  ```bash
  export DB_DRIVER=postgres
  export DB_USER=user
  export DB_PASS=password
  export DB_HOST=localhost
  export DB_PORT=5432
  export DB_NAME=test
  export MIGRATIONS=./migrations
  ```

## Migration Files

Place your SQL migration files in the directory specified by the `MIGRATIONS` environment variable. Each file should be named with a numeric prefix indicating its version, followed by an optionally descriptive name. Only the numeric prefix is used to determine the migration order.

### Example Directory Structure

  ```
  migrations/
  ├── 101_create_users_table.sql
  ├── 102_add_email_to_users.sql
  └── 103_create_orders_table.sql
  ```

### Migration File Format

Each migration file must include both `up` and `down` SQL statements, separated by special directives:

  ```sql
  -- @up

  CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL
  );

  -- @down

  DROP TABLE users;
  ```

- The `-- @up` and `-- @down` directives must appear **in order**, each on their **own line**, and with no additional content.
- Any content before the `-- @up` directive is ignored and can be used for comments or documentation.
- All SQL between `-- @up` and `-- @down` is treated as the migration "up" step; everything after `-- @down` is treated as the "down" step (rollback).

## Docker Support

Build a Docker image tagged with the latest Git tag (i.e., the version):

  ```
  task build:docker
  ```

This command extracts the latest Git tag, removes the leading 'v' if present, and builds the Docker image accordingly.

## Additional Task Commands

- `task clean` Remove the `bin` directory.
- `task fmt` Format the code using `gofmt`.
- `task uninstall` Remove the installed binary.

## Technical Details

### Migration Version Tracking

`schemer` stores the current migration version in a table named:

  ```
  schemer
  ```

This table contains a single `BIGINT` value representing the version number, which is derived from the numeric prefix of each migration file (e.g., `101_add_users_table.sql` → version `101`).

> **Note:** Version numbers must be positive and fit within the range of a 64-bit signed integer (`int64` in Go). Values outside this range will cause errors.

### Transactions and Rollbacks

Each migration (both `up` and `down`) runs inside its own database transaction. If a migration fails, any changes made during that version are rolled back.

- Migrations are **atomic per version**.
- Failed migrations do **not** alter the database or version table.
- There is no "dirty" state recorded — failures are surfaced to the user, but leave the database unchanged.

## License

This project is licensed under the [MIT License](LICENSE).

## Commit Message Convention

This project follows the [Conventional Commits](https://www.conventionalcommits.org/) specification for commit messages. Please refer to the official documentation for guidelines on how to format your commits.
