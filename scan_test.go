package main

import (
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	text := `0 1    2	3 # comment
	4 5 	6 7

 # comment 1 2 3
8 9 10`
	expect := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	i := 0
	s := NewListScanner(strings.NewReader(text))
	for s.Scan() {
		if txt := s.Text(); txt != expect[i] {
			t.Fatalf("scanned: %s, expect: %s", txt, expect[i])
		}
		i++
	}
	if err := s.Err(); err != nil {
		t.Fatal(err)
	}
}
