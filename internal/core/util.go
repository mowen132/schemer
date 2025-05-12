// Copyright (c) 2025 Mark Owen
// Licensed under the MIT License. See LICENSE file in the project root for details.

package core

import (
	"bufio"
	"io"
	"os"
)

func readFileLines(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()
	return readLines(file)
}

func readLines(r io.Reader) ([]string, error) {
	var lines []string
	s := bufio.NewScanner(r)

	for s.Scan() {
		lines = append(lines, s.Text())
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
