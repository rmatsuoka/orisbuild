package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type settings struct {
	Length int
	Prefix string
	Macro  map[string][]string
	Color  string
}

type Mod struct {
	length int
	prefix string
	rule   [][2]string
	macro  map[string][]string
	color  string
}

var (
	numre = regexp.MustCompile(`^[0-9]+$`)
)

func UnmarshalMod(dir string) (*Mod, error) {
	if dir == "" || dir == "." || dir == ".." || numre.MatchString(dir) {
		return nil, fmt.Errorf("unmarshal module \"%s\": invalid Module Name", dir)
	}

	setb, err := os.ReadFile(filepath.Join(dir, "settings.json"))
	if err != nil {
		return nil, fmt.Errorf("unmarshal module \"%s\": %v", dir, err)
	}

	s := new(settings)
	if err := json.Unmarshal(setb, &s); err != nil {
		return nil, fmt.Errorf("unmarshal module \"%s\": %v", dir, err)
	}

	m := &Mod{length: s.Length, prefix: s.Prefix, macro: s.Macro, color: s.Color}
	ruleF, err := os.Open(filepath.Join(dir, "rule.pair"))
	if err != nil {
		return nil, fmt.Errorf("unmarshal module \"%s\": %v", dir, err)
	}
	defer ruleF.Close()

	if err = m.readRule(ruleF); err != nil {
		return nil, fmt.Errorf("unmarshal module \"%s\": %v", dir, err)
	}

	return m, nil
}

func (m *Mod) Color() (*CategoryColor, bool) {
	if m.color == "" {
		return nil, false
	}
	return &CategoryColor{m.color}, true
}

func (m *Mod) Transcript() []string {
	tran := make([]string, m.length)
	for i := 0; i < m.length; i++ {
		tran[i] = fmt.Sprintf("%s%d", m.prefix, i)
	}
	return tran
}

func (m *Mod) Rule() [][2]string {
	return m.rule
}

func (m *Mod) readRule(r io.Reader) error {
	s := bufio.NewScanner(r)

	f := func(p, n string) string {
		if numre.MatchString(n) {
			return p + n
		}
		return n
	}

	for s.Scan() {
		r := strings.Fields(s.Text())
		if len(r) < 2 {
			continue
		}
		if r[0] == "#" || r[1] == "#" {
			continue
		}

		if a, ok := m.macro[r[0]]; ok {
			c := f(m.prefix, r[1])
			for _, x := range a {
				m.rule = append(m.rule, [2]string{x, c})
			}
		} else if a, ok := m.macro[r[1]]; ok {
			c := f(m.prefix, r[0])
			for _, x := range a {
				m.rule = append(m.rule, [2]string{c, x})
			}
		} else {
			m.rule = append(m.rule, [2]string{f(m.prefix, r[0]), f(m.prefix, r[1])})
		}
	}
	return s.Err()
}
