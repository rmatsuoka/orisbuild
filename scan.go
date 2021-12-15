package main

import (
	"bufio"
	"strings"
	"io"
)

type Scanner struct {
	s *bufio.Scanner
	words []string
	offset int
}

func NewListScanner(r io.Reader) *Scanner {
	return &Scanner{s: bufio.NewScanner(r), offset: -1}
}

func (s *Scanner) Scan() bool {
	s.offset++
	if len(s.words) == s.offset || "#" == s.words[s.offset] {
		for {
			if !s.s.Scan() {
				return false
			}
			s.offset = 0
			s.words = strings.Fields(s.s.Text())
			if len(s.words) > 0 && "#" != s.words[0] {
				break
			}
		}
	}
	return true
}

func (s *Scanner) Err() error {
	return s.s.Err()
}

func (s *Scanner) Text() string {
	return s.words[s.offset]
}
