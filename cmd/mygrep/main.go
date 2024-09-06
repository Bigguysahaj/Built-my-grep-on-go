package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
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
	hasCircumflex := pattern[0] == '^'
	hasDollar := strings.HasSuffix(pattern,"$")
	i := 0
	if (hasCircumflex || hasDollar) {
		i = 1
	}
	if hasDollar {
		// line = line[:len(line)-1]
		line = []byte(ReverseString(string(line)))
		pattern = ReverseString(pattern)
		fmt.Println("lines and patterns", line, pattern)
	}
	for startIndex := 0; startIndex < len(line); startIndex++ {
			ok := true
			lineIndex := startIndex

			for ; i < len(pattern); i++ {
					if lineIndex >= len(line) {
							fmt.Println("segmentation fault")
							ok = false
							break
					}
					
					if pattern[i] == '\\' && i+1 < len(pattern) {
							if pattern[i+1] == 'd' {
									if !bytes.ContainsAny([]byte{line[lineIndex]}, "0123456789") {
											fmt.Println("Error on \\d")
											ok = false
											break
									}
									i++
									lineIndex++
							} else if pattern[i+1] == 'w' {
									if !bytes.ContainsAny([]byte{line[lineIndex]}, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_") {
										fmt.Println("Error on \\w")	
										ok = false
										break
									}
									i++
									lineIndex++
							} else {
									if line[lineIndex] != '\\' {

										ok = false
										break
									}
									lineIndex++
							}
					} else if pattern[i] == ' ' {
							if line[lineIndex] != ' ' {
								fmt.Println("Error on space")
									ok = false
									break
							}
							lineIndex++
					} else if pattern[i] == '[' {
						closeBracket := bytes.IndexByte([]byte(pattern[i:]), ']')
						if closeBracket == -1 {
							return false, fmt.Errorf("unclosed character class")
						}
						positiveCharacters := pattern[i+1 : i+closeBracket]
						if !bytes.ContainsAny(line, positiveCharacters) {
							fmt.Println("No positive character error")
							ok = false
						}
						if bytes.ContainsAny([]byte(positiveCharacters), "^") {
							ok = !ok
						}

						if !ok {
							fmt.Println("Lost to brackets")
							return false, nil
						}
						i += closeBracket
					} else {
						if line[lineIndex] != pattern[i] {
							fmt.Printf("Basic mismatch error %c, %c \n", line[lineIndex], pattern[i])
							ok = false
							break
						}
						lineIndex++
					}
			}
			
			if !ok && (hasCircumflex || hasDollar) {
				break 
			}

			if ok {
				fmt.Println("Your word ", string(line), "contains the pattern", pattern)
				return true, nil
			}
		}

	fmt.Println("Your word ", string(line), " doesn't contains the pattern", pattern)		
	return false, nil
}

func ReverseString(s string) string {
	runes := []rune(s)
	size := len(runes)
	for i, j := 0, size-1; i < size>>1; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

