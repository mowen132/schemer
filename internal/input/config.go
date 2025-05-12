// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package input

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Driver     string
	DSN        string
	User       string
	Pass       string
	Host       string
	Port       int
	Name       string
	Migrations string
}

func LoadConfig() (*Config, error) {
	driver, err := mustGet("DB_DRIVER")

	if err != nil {
		return nil, err
	}

	dsn := os.Getenv("DB_CONN")
	var user string
	var pass string
	var host string
	var port int
	var name string

	if dsn == "" {
		user, err = mustGet("DB_USER")

		if err != nil {
			return nil, err
		}

		pass, err = mustGet("DB_PASS")

		if err != nil {
			return nil, err
		}

		host, err = mustGet("DB_HOST")

		if err != nil {
			return nil, err
		}

		port, err = getPort("DB_PORT")

		if err != nil {
			return nil, err
		}

		name, err = mustGet("DB_NAME")

		if err != nil {
			return nil, err
		}
	}

	migrations, err := mustGet("MIGRATIONS")

	if err != nil {
		return nil, err
	}

	return &Config{
		Driver:     driver,
		DSN:        dsn,
		User:       user,
		Pass:       pass,
		Host:       host,
		Port:       port,
		Name:       name,
		Migrations: migrations,
	}, nil
}

func mustGet(key string) (string, error) {
	val := os.Getenv(key)

	if val == "" {
		return "", fmt.Errorf("%s is not set or is empty", key)
	}

	return val, nil
}

func getPort(key string) (int, error) {
	if s := os.Getenv(key); s != "" {
		i, err := strconv.ParseInt(s, 10, 64)

		if err != nil || i < 1 || i > 65535 {
			return 0, fmt.Errorf("unable to parse %s: %q is not a valid port", key, s)
		}

		return int(i), nil
	}

	return -1, nil
}
