package main

import (
	"os"
	"testing"
)

func TestUnmarshalMod(t *testing.T) {
	os.Chdir("./sample")
	m, err := UnmarshalMod("L")
	os.Chdir("..")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", m)
}