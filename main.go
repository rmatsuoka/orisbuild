package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
)

var (
	modFileName  = flag.String("mod", "mod.list", "a module file")
	seedFileName = flag.String("seed", "seed.list", "a seed file")
	delay        = flag.Int("delay", 3, "delay")
	name         = flag.String("name", "oritatami_system", "name")
	count        = flag.Int("count", 1, "period count")
	highlightMod = flag.String("high", "", "specify a module to highlight")
)

type Oris struct {
	Name       string                    `json:"name"`
	Rule       [][2]string               `json:"rule"`
	Transcript []string                  `json:"compactTranscriptPeriod"`
	Seed       []string                  `json:"seedConformation"`
	Delay      int                       `json:"delay"`
	Count      int                       `json:"periodCount"`
	Colors     map[string]*CategoryColor `json:"categoryColors"`
}

type CategoryColor struct {
	Color string `json:"name"`
}

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)
	flag.Parse()

	oris := Oris{
		Name:  *name,
		Delay: *delay,
		Count: *count,
	}

	seedF, err := os.Open(*seedFileName)
	if err != nil {
		log.Fatal(err)
	}
	if oris.Seed, err = unmarshalSeed(seedF); err != nil {
		seedF.Close()
		log.Fatal(err)
	}
	seedF.Close()

	var modF *os.File
	modF, err = os.Open(*modFileName)
	if err != nil {
		log.Fatal(err)
	}
	oris.Rule, oris.Transcript, oris.Colors, err = readModFile(modF)
	if err != nil {
		modF.Close()
		log.Fatal(err)
	}
	modF.Close()

	var b []byte
	b, err = json.Marshal(oris)
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	if err = json.Indent(buf, b, "", "    "); err != nil {
		log.Fatal(err)
	}
	buf.WriteString("\n")
	buf.WriteTo(os.Stdout)
}

func readModFile(r io.Reader) (rule [][2]string, transcript []string, colors map[string]*CategoryColor, err error) {
	s := NewListScanner(r)
	seenMod := make(map[string]*Mod)
	colors = make(map[string]*CategoryColor)

	for s.Scan() {
		modName := s.Text()
		mod, ok := seenMod[modName]
		if !ok {
			mod, err = UnmarshalMod(modName)
			if err != nil {
				return nil, nil, nil, err
			}
			if c, ok := mod.Color(); ok {
				colors[modName] = c
			}
			seenMod[modName] = mod
		}
		transcript = append(transcript, mod.Transcript()...)
	}
	if err = s.Err(); err != nil {
		return nil, nil, nil, err
	}

	for _, mod := range seenMod {
		rule = append(rule, mod.Rule()...)
	}

	// highlight one module.
	// other modules are filled in Indigo500
	if *highlightMod != "" {
		for k := range seenMod {
			if k == *highlightMod {
				continue
			}
			colors[k] = &CategoryColor{"grey300"}
		}
	}

	return rule, transcript, colors, nil
}

func unmarshalSeed(r io.Reader) ([]string, error) {
	s := NewListScanner(r)
	var seed []string
	for s.Scan() {
		seed = append(seed, s.Text())
	}
	return seed, s.Err()
}
