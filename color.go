package main

import (
	"encoding/json"
	"fmt"
)

type Color struct {
	data []byte
}

func (c *Color) UnmarshalJSON(b []byte) error {
	if !json.Valid(b) {
		return fmt.Errorf("invalid json")
	}
	c.data = make([]byte, len(b))
	copy(c.data, b)
	return nil
}

func (c *Color) MarshalJSON() ([]byte, error) {
	if json.Valid(c.data) {
		return c.data, nil
	}
	return nil, fmt.Errorf("marshal %s: strange value", c.data)
}
