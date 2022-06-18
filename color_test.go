package main

import (
	"encoding/json"
	"testing"
)

func TestColor(t *testing.T) {
	colors := []*Color{
		{[]byte(`{"name": "red500"}`)},
		{[]byte(`{"rgba": [1, 2, 3, 4]}`)},
	}
	for _, c := range colors {
		j, err := json.Marshal(c)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%v -> %s", c, j)
	}
}