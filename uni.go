package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strconv"
)

//go:embed UnicodeData-15.0.0d5.txt
var unicodeData []byte

var lines [][]byte

var just *bool = flag.Bool("just", false, "print just the first matched code point")

func main() {
	flag.Parse()
	query := bytes.ToUpper([]byte(flag.Arg(0)))
	lines = bytes.Split(unicodeData, []byte{'\n'})
	for _, line := range lines {
		if bytes.Contains(line, []byte("<control>")) {
			continue
		}
		if bytes.Contains(line, query) {
			cols := bytes.SplitN(line, []byte{';'}, 2)
			cp, err := strconv.ParseUint(string(cols[0]), 16, 64)
			if err != nil {
				panic(err)
			}
			if *just {
				fmt.Printf("%c", cp)
				os.Exit(0)
			}
			fmt.Printf("%c\t%s\t%s\n", cp, cols[0], cols[1])
		}
	}
}
