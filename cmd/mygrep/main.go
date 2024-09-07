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

func matchLine(line []byte, pattern string) (bool, error) {
	parenExpressions := []string{""}
	lineIndex := 0
	patternIndex := 0

	fmt.Printf("Starting matchLine with line: %s and pattern: %s\n", string(line), pattern)

	for patternIndex < len(pattern) {
		fmt.Printf("Current state: lineIndex=%d, patternIndex=%d\n", lineIndex, patternIndex)

		if lineIndex >= len(line) && patternIndex < len(pattern) {
			fmt.Println("Reached end of line before end of pattern")
			return false, nil
		}

		if pattern[patternIndex] == '(' {
			closeParen := strings.IndexByte(pattern[patternIndex:], ')')
			if closeParen == -1 {
				fmt.Println("Error: unclosed parenthesis")
				return false, fmt.Errorf("unclosed parenthesis")
			}
			expression := pattern[patternIndex+1 : patternIndex+closeParen]
			fmt.Printf("Found capture group: %s\n", expression)
			if strings.HasPrefix(string(line[lineIndex:]), expression) {
				parenExpressions = append(parenExpressions, expression)
				fmt.Printf("Matched capture group. parenExpressions: %v\n", parenExpressions)
				lineIndex += len(expression)
				patternIndex += closeParen + 1
			} else {
				fmt.Printf("Failed to match capture group: %s\n", expression)
				return false, nil
			}
		} else if pattern[patternIndex] == '\\' && patternIndex+1 < len(pattern) {
			if pattern[patternIndex+1] >= '1' && pattern[patternIndex+1] <= '9' {
				groupNum := int(pattern[patternIndex+1] - '0')
				fmt.Printf("Found backreference: \\%d\n", groupNum)
				if groupNum < len(parenExpressions) {
					expression := parenExpressions[groupNum]
					fmt.Printf("Attempting to match backreference: %s\n", expression)
					if strings.HasPrefix(string(line[lineIndex:]), expression) {
						fmt.Println("Backreference matched successfully")
						lineIndex += len(expression)
						patternIndex += 2
					} else {
						fmt.Println("Failed to match backreference")
						return false, nil
					}
				} else {
					fmt.Printf("Error: invalid backreference \\%d\n", groupNum)
					return false, fmt.Errorf("invalid backreference")
				}
			} else if pattern[patternIndex+1] == 'w' {
				fmt.Printf("Matching \\w with character: %c\n", line[lineIndex])
				if !bytes.ContainsAny([]byte{line[lineIndex]}, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_") {
					fmt.Println("Failed to match \\w")
					return false, nil
				}
				patternIndex += 2
				lineIndex++
			} else if pattern[patternIndex+1] == 'd' {
				fmt.Printf("Matching \\d with character: %c\n", line[lineIndex])
				if !bytes.ContainsAny([]byte{line[lineIndex]}, "0123456789") {
					fmt.Println("Failed to match \\d")
					return false, nil
				}
				patternIndex += 2
				lineIndex++
			} else {
				fmt.Printf("Unhandled escape sequence: \\%c\n", pattern[patternIndex+1])
				patternIndex += 2
				lineIndex++
			}
		} else if pattern[patternIndex] == line[lineIndex] {
			fmt.Printf("Matched character: %c\n", pattern[patternIndex])
			patternIndex++
			lineIndex++
		} else {
			fmt.Printf("Failed to match character: expected %c, got %c\n", pattern[patternIndex], line[lineIndex])
			return false, nil
		}
	}

	matchResult := lineIndex == len(line)
	fmt.Printf("Finished matching. Result: %v\n", matchResult)
	return matchResult, nil
}

// mastermind function
// func matchLine(line []byte, pattern string) (bool, error) {
// 	parenExpressions := []string{"filler"}
// 	hasCircumflex := pattern[0] == '^'
// 	hasPlus := false
// 	lineIndex := 0 
// 	hasDollar := strings.HasSuffix(pattern,"$")
// 	i := 0
// 	if (hasCircumflex || hasDollar) {
// 		i = 1
// 	}
// 	if hasDollar {
// 		// line = line[:len(line)-1]
// 		line = []byte(ReverseString(string(line)))
// 		pattern = ReverseString(pattern)
// 		fmt.Println("lines and patterns", line, pattern)
// 	}
// 	for startIndex := 0; startIndex < len(line); startIndex++ {
// 		if lineIndex >= len(line) {
// 			fmt.Println("segmentation fault outer loop")
			
// 			break
// 		}
// 		ok := true
// 		lineIndex = startIndex
		

// 			for ; i < len(pattern); i++ {
// 					if lineIndex >= len(line) {
// 						fmt.Println("segmentation fault inner loop")
// 						ok = (i == len(pattern)-1)
// 						break
// 					}
					
// 					if pattern[i] == '\\' && i+1 < len(pattern) {
// 							if pattern[i+1] == 'd' {
// 									if !bytes.ContainsAny([]byte{line[lineIndex]}, "0123456789") {
// 											fmt.Println("Error on \\d")
// 											ok = false
// 											break
// 									}
// 									i++
// 									lineIndex++
// 							} else if pattern[i+1] == 'w' {
// 									if !bytes.ContainsAny([]byte{line[lineIndex]}, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_") {
// 										fmt.Println("Error on \\w")	
// 										ok = false
// 										break
// 									}
// 									i++
// 									lineIndex++
// 							} else if pattern[i+1] >= '1' && pattern[i+1] <= '9' {
// 								currExpression := parenExpressions[int(pattern[i+1] - '0')]
// 								for lineIndex < len(line) {
// 									fmt.Println("we have in repeat", string(line[lineIndex:]), currExpression)
// 									if strings.HasSuffix(string(line[lineIndex:]) , currExpression) {
// 										ok = true
// 										lineIndex += len(currExpression)
// 										break
// 									} else {
// 										ok = false
// 										lineIndex++
// 									}
// 								}
// 							} else {
// 								if line[lineIndex] != '\\' {

// 									ok = false
// 									break
// 								}
// 								lineIndex++
// 							}
// 					} else if pattern[i] == ' ' {
// 							if line[lineIndex] != ' ' {
// 								fmt.Println("Error on space")
// 									ok = false
// 									break
// 							}
// 							lineIndex++
// 					} else if pattern[i] == '[' {
// 						closeBracket := bytes.IndexByte([]byte(pattern[i:]), ']')
// 						if closeBracket == -1 {
// 							return false, fmt.Errorf("unclosed character class")
// 						}
// 						positiveCharacters := pattern[i+1 : i+closeBracket]
// 						if !bytes.ContainsAny(line, positiveCharacters) {
// 							fmt.Println("No positive character error")
// 							ok = false
// 						}
// 						if bytes.ContainsAny([]byte(positiveCharacters), "^") {
// 							ok = !ok
// 						}

// 						if !ok {
// 							fmt.Println("Lost to brackets")
// 							return false, nil
// 						}
// 						i += closeBracket
// 					} else if pattern[i] == '(' {
// 						closeParen := bytes.IndexByte([]byte(pattern[i:]), ')')
						
// 						expression := pattern[i+1 : i+closeParen]
// 						alternatives := strings.Split(expression, "|")  
// 						if len(alternatives) > 1 {
// 							alternationMatch := false
							
// 							for _, alt := range alternatives {
// 								if strings.HasPrefix(string(line[lineIndex:]), alt) {
// 									alternationMatch = true
// 									lineIndex = len(alt)
// 									break
// 								}
// 							}
							
// 							if !alternationMatch {
// 								fmt.Println("got faltered at alternation")
// 								ok = false
// 								break
// 							}
							
// 							i += closeParen
// 						} else {
// 							// Handles cases of "(cat) and \1"
// 							parenExpressions = append(parenExpressions, expression)
// 							for lineIndex < len(line) {
// 								fmt.Println("we have ", string(line[lineIndex:]), expression)
// 								if strings.HasSuffix(string(line[lineIndex:]) , expression) {
// 									ok = true
// 									lineIndex += len(expression)
// 									i+= len(expression) + 1
// 									break
// 								} else {
// 									ok = false
// 									lineIndex++
// 								}
// 							}
// 						}
							
// 					} else if pattern[i] == '+' {
// 						hasPlus = true
// 						for (line[lineIndex] == pattern[i-1]){
// 							lineIndex++
// 						}
// 					} else if pattern[i] == '.'{
// 						fmt.Println("We are doing nothing")
// 					} else {
// 						fmt.Printf("Let's see all %c, %c \n", line[lineIndex], pattern[i])
// 						if line[lineIndex] != pattern[i] {
// 							// condition for '?'
// 							if pattern[i] == '?' {
// 								lineIndex--
// 							} else {
// 								if (i+1 < len(pattern) && pattern[i+1] == '?') {
// 									i += 2
// 									ok = pattern[i] == line[lineIndex]
// 								} else{
// 									fmt.Printf("Basic mismatch error %c, %c \n", line[lineIndex], pattern[i])
// 									ok = false
// 									break
// 								}
							
// 							}
// 						} 
// 						lineIndex++
// 					}
// 			}
			
// 			if !ok && (hasCircumflex || hasDollar || hasPlus ) {
// 				break 
// 			}

// 			if ok {
// 				fmt.Println("Your word ", string(line), "contains the pattern", pattern)
// 				return true, nil
// 			}
// 		}

// 	fmt.Println("Your word ", string(line), " doesn't contains the pattern", pattern)		
// 	return false, nil
// }

func ReverseString(s string) string {
	runes := []rune(s)
	size := len(runes)
	for i, j := 0, size-1; i < size>>1; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

