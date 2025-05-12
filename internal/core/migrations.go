// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package core

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/mowen132/schemer/internal/input"
)

type Migrations []*Migration

type Migration struct {
	Name    string
	Path    string
	Version int64
	Up      string
	Down    string
}

type actionType int

const (
	noAction actionType = iota
	upAction
	downAction
)

func LoadMigrations(dir string) (Migrations, error) {
	files, err := os.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	migrations := make(Migrations, 0, len(files))
	versions := map[int64]struct{}{}

	for _, f := range files {
		name := f.Name()
		version, err := parseVersion(name)

		if err != nil {
			return nil, err
		}

		if _, ok := versions[version]; ok {
			return nil, fmt.Errorf("duplicate version %d found", version)
		}

		migrations = append(migrations, &Migration{
			Name:    name,
			Path:    filepath.Join(dir, name),
			Version: version,
		})

		versions[version] = struct{}{}
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func (ms Migrations) Run(cmd *input.Command, conn *Connection, version int64) error {
	plan, err := resolveExecutionPlan(cmd, version, ms)

	if err != nil {
		return err
	}

	return plan.execute(conn, ms)
}

func (ms Migrations) search(version int64) (int, bool) {
	if version == 0 {
		return -1, true
	}

	i, j := 0, len(ms)

	for j-i > 1 {
		k := (i + j) / 2

		if version < ms[k].Version {
			j = k
		} else {
			i = k
		}
	}

	return i, ms[i].Version == version
}

func (m *Migration) run(action actionType, conn *Connection, final int64) error {
	get := func(action actionType, up, down string) string {
		if action == upAction {
			return up
		}

		return down
	}

	fmt.Printf("%s %q\n", get(action, "UP", "DOWN"), m.Name)

	if err := m.load(); err != nil {
		return err
	}

	return conn.transact(func(tx *sql.Tx) error {
		if _, err := tx.Exec(get(action, m.Up, m.Down)); err != nil {
			return fmt.Errorf("failed to execute: %v", err)
		}

		if err := conn.saveVersion(tx, final); err != nil {
			return err
		}

		return nil
	})
}

func (m *Migration) load() error {
	lines, err := readFileLines(m.Path)

	if err != nil {
		return err
	}

	loader := newMigrationLoader(lines)

	if err := loader.skipHeader(); err != nil {
		return err
	}

	up, err := loader.loadUpSection()

	if err != nil {
		return err
	}

	down, err := loader.loadDownSection()

	if err != nil {
		return err
	}

	m.Up = up
	m.Down = down
	return nil
}

type executionPlan struct {
	action actionType
	beg    int
	end    int
}

func resolveExecutionPlan(cmd *input.Command, version int64, migrations Migrations) (*executionPlan, error) {
	cur, ok := migrations.search(version)

	if !ok {
		return nil, fmt.Errorf("current version %d not found", version)
	}

	plan := &executionPlan{}

	switch cmd.Type {
	case input.UpCommand:
		plan.action = upAction
		plan.beg = cur + 1
		last := len(migrations) - 1

		if cmd.HasOperand {
			plan.end = min(cur+int(cmd.Operand), last)
		} else {
			plan.end = last
		}

	case input.DownCommand:
		plan.action = downAction
		plan.beg = cur

		if cmd.HasOperand {
			plan.end = max(cur-int(cmd.Operand), -1)
		} else {
			plan.end = -1
		}

	case input.GotoCommand:
		tar, ok := migrations.search(cmd.Operand)

		if !ok {
			return nil, fmt.Errorf("target version %d not found", cmd.Operand)
		}

		if cur < tar {
			plan.action = upAction
			plan.beg = cur + 1
			plan.end = tar
		} else if cur > tar {
			plan.action = downAction
			plan.beg = cur
			plan.end = tar
		} else {
			plan.action = noAction
		}
	}

	return plan, nil
}

func (p *executionPlan) execute(conn *Connection, migrations Migrations) error {
	switch p.action {
	case upAction:
		for i := p.beg; i <= p.end; i++ {
			migration := migrations[i]

			if err := migration.run(upAction, conn, migration.Version); err != nil {
				return fmt.Errorf("%q: %v", migration.Name, err)
			}
		}

	case downAction:
		if p.beg == -1 {
			return nil
		}

		zero := p.end == -1

		if zero {
			p.end = 0
		}

		for i := p.beg; i > p.end; i-- {
			migration := migrations[i]

			if err := migration.run(downAction, conn, migrations[i-1].Version); err != nil {
				return fmt.Errorf("%q: %v", migration.Name, err)
			}
		}

		if zero {
			migration := migrations[0]

			if err := migration.run(downAction, conn, 0); err != nil {
				return fmt.Errorf("%q: %v", migration.Name, err)
			}
		}
	}

	return nil
}

type migrationLoader struct {
	lines []string
	i     int
	n     int
}

func newMigrationLoader(lines []string) *migrationLoader {
	return &migrationLoader{
		lines: lines,
		i:     0,
		n:     len(lines),
	}
}

func (l *migrationLoader) skipHeader() error {
	for {
		if l.eof() {
			return invalidFileErrorf("missing required '-- @up'")
		}

		if m, ok := matchDirective(l.line()); ok {
			if m == "up" {
				l.next()
				return nil
			}

			return invalidFileErrorf("found '-- @down' before '-- @up'")
		}

		l.next()
	}
}

func (l *migrationLoader) loadUpSection() (string, error) {
	var b strings.Builder

	for {
		if l.eof() {
			return "", invalidFileErrorf("missing required '-- @down'")
		}

		line := l.line()

		if m, ok := matchDirective(line); ok {
			if m == "down" {
				l.next()
				return b.String(), nil
			}

			return "", invalidFileErrorf("found duplicate '-- @up'")
		}

		fmt.Fprintln(&b, line)
		l.next()
	}
}

func (l *migrationLoader) loadDownSection() (string, error) {
	var b strings.Builder

	for !l.eof() {
		line := l.line()

		if m, ok := matchDirective(line); ok {
			return "", invalidFileErrorf("unexpected '-- @%s'", m)
		}

		fmt.Fprintln(&b, line)
		l.next()
	}

	return b.String(), nil
}

func (l *migrationLoader) eof() bool {
	return l.i == l.n
}

func (l *migrationLoader) line() string {
	return l.lines[l.i]
}

func (l *migrationLoader) next() {
	l.i++
}

var versionRe = regexp.MustCompile(`^\d+`)

func parseVersion(name string) (int64, error) {
	s := versionRe.FindString(name)

	if s == "" {
		return 0, fmt.Errorf("%q does not specify a version", name)
	}

	version, err := strconv.ParseInt(s, 10, 64)

	if err != nil || version == 0 {
		return 0, fmt.Errorf("%q does not have a valid version", name)
	}

	return version, nil
}

var directiveRe = regexp.MustCompile(`^\s*--\s*@(up|down)\s*$`)

func matchDirective(line string) (string, bool) {
	if m := directiveRe.FindStringSubmatch(line); m != nil {
		return m[1], true
	}

	return "", false
}

func invalidFileErrorf(format string, args ...any) error {
	return fmt.Errorf("invalid file: "+format, args...)
}
