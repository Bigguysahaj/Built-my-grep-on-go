package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
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

}

// mastermind function
func matchLine(line []byte, pattern string) (bool, error) {

	for startIndex :=0 ; startIndex < len(line) ; startIndex++ {
		ok := true
		lineIndex := startIndex

		for i:=0 ; i < len(pattern) ; i++ {
			if lineIndex >= len(line) {
				ok = false
			}

			if pattern[i] == '\\' && i+1 < len(pattern) {
				if (pattern[i+1] == 'd') {
					if !bytes.ContainsAny([]byte{line[lineIndex]}, "0123456789"){
						ok = false
						break
					}
					i++
				} else if (pattern[i+1] == 'd') {
					if !bytes.ContainsAny([]byte{line[lineIndex]},  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_") {
						ok = false
						break
					}
					i++
				} else {
					// Handle literal backslash
					if line[lineIndex] != '\\' {
							ok = false
							break
					}
				}
			} else if pattern[i] != line[lineIndex] {
				ok = false
				break
			}
			lineIndex++
		}
		if ok {
			fmt.Println("Your word ", string(line), "contains the pattern", pattern)
			return true, nil
		}
	}
	fmt.Println("Your word ", string(line), " doesn't contains the pattern", pattern)		
	return false, nil
}

