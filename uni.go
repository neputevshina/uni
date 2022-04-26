package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strconv"
	"unicode/utf8"
)

//go:embed UnicodeData-15.0.0d5.txt
var unicodeData []byte

var lines [][]byte

var just = flag.Bool("just", false, "print just the first matched code point")
var cl = flag.String("cl", "", "show only symbols from the specified class")
var fuzzy = flag.Bool("f", false, "perform a fuzzy search and print the counterweight for sorting")
var glyph = flag.String("g", "", "give a description to the specified unicode glyph; incompatible with other flags")

func bye(msg string) {
	fmt.Fprintln(os.Stderr, os.Args[0]+": "+msg)
	os.Exit(1)
}

func main() {
	flag.Parse()
	query := bytes.ToUpper([]byte(flag.Arg(0)))
	lines = bytes.Split(unicodeData, []byte{'\n'})
	if (*fuzzy || *cl != "" || *just) && *glyph != "" {
		bye("-g is incompatible with other flags")
	}

	for _, line := range lines {
		if bytes.Contains(line, []byte("<control>")) {
			continue
		}
		if score := containsSeq(line, query); score >= 0 {
			if !*fuzzy && score != 0 {
				continue
			}
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
			if *fuzzy {
				fmt.Printf("%v\t%c\t%s\t%s\n", score, cp, cols[0], cols[1])
			} else {
				fmt.Printf("%c\t%s\t%s\n", cp, cols[0], cols[1])
			}
		}
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
