package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strconv"
	"unicode/utf8"

	"slices"
)

//go:embed UnicodeData-15.0.0d5.txt
var unicodeData []byte

var lines [][]byte

var just = flag.Bool("just", false, "print just the first matched code point")
var cl = flag.String("cl", "", "show only symbols from the specified class")

var pscore = flag.Bool("s", false, "print fuzzy search score")
var glyph = flag.String("g", "", "give a description to the specified unicode glyph; incompatible with other flags")

func bye(msg string) {
	fmt.Fprintln(os.Stderr, os.Args[0]+": "+msg)
	os.Exit(1)
}

func main() {
	flag.Parse()
	query := bytes.ToUpper([]byte(flag.Arg(0)))
	lines = bytes.Split(unicodeData, []byte{'\n'})
	if (*cl != "" || *just) && *glyph != "" {
		bye("-g is incompatible with other flags")
	}
	strs := []string{}
	nonzero := []string{}

	for _, line := range lines {
		if bytes.Contains(line, []byte("<control>")) {
			continue
		}
		if score := containsSeq(line, query); score >= 0 {
			cols := bytes.SplitN(line, []byte{';'}, 2)
			if *cl != "" {
				cols := bytes.SplitN(cols[1], []byte{';'}, 3)
				if string(cols[1]) != *cl {
					continue
				}
			}
			cp, _ := strconv.ParseUint(string(cols[0]), 16, 64)
			if *glyph != "" {
				r, sz := utf8.DecodeRuneInString(*glyph)
				_, sz = utf8.DecodeRuneInString((*glyph)[sz:])
				if sz != 0 {
					bye("-g flag requires only one code point to be present")
				}

				if cp != uint64(r) {
					continue
				}
			}
			if *just {
				fmt.Printf("%c", cp)
				os.Exit(0)
			}
			s := ""
			if *pscore {
				s = fmt.Sprintf("%3d\t", score)
			}
			o := fmt.Sprintf("%s%c\t%s\t%s", s, cp, cols[0], cols[1])

			if score == 0 {
				strs = append(strs, o)
			} else {
				nonzero = append(nonzero, o)
			}
		}
	}
	slices.Sort(nonzero)
	for _, v := range strs {
		fmt.Println(v)
	}
	for _, v := range nonzero {
		fmt.Println(v)
	}
}

// not a real substring search
func containsSeq(s, sep []byte) (skip int) {
	if len(s) < len(sep) {
		return -1
	}
	if bytes.Contains(s, sep) {
		return 0
	}
	first := true
	for _, b := range s {
		n := len(sep)
		if n == 0 {
			return
		}
		if b == sep[0] {
			sep = sep[1:]
			if len(sep) == 0 {
				return
			}
			first = false
		} else if !first {
			skip++
		}
	}
	return -1
}
