package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
)

// bytes: Provides functions for byte slice manipulation, like checking for character presence.
// fmt: Handles formatted I/O, similar to printf in C/C++.
// io: Contains functions for reading and writing data, similar to file handling.
// os: Provides OS-level functionality, including command-line argument handling.
// unicode/utf8: Helps with working with UTF-8 encoded strings, ensuring that characters are handled correctly.

// Ensures gofmt doesn't remove the "bytes" import above (feel free to remove this!)
var _ = bytes.ContainsAny

// Usage: echo <input_text> | your_program.sh -E <pattern>

func main() {
	// os.Args holds command-line arguments. The program expects at least three arguments: the program name, -E flag, and a pattern.
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	// hit: interesting, a destruct very simillar to python and js, very concise execution of err
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}
	
	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

// mastermind function
func matchLine(line []byte, pattern string) (bool, error) {

	// Checks if the character is more than one character long, characters of 1 length not supported.
	// if utf8.RuneCountInString(pattern) != 1 {
	// 	return false, fmt.Errorf("unsupported pattern: %q", pattern)
	// }

	var flag bool = false
	// var ok bool 

	for _, runeValue := range pattern {
		if runeValue == '[' || runeValue == ']' {
			continue
		}

		var regPattern = regexp.MustCompile(string(runeValue))
		ok := regPattern.Match(line)

		fmt.Println(ok)

		if ok {
			flag = !flag
			break
		}
	}

	
	if flag{
		fmt.Println("Your word ", string(line), "contains the pattern", pattern)
		os.Exit(0)
	} else {
		fmt.Println("Your word ", string(line), " doesn't contains the pattern", pattern)
		os.Exit(1)
	}

	return flag, nil
}
