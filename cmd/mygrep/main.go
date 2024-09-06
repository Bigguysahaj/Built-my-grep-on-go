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

  var ok bool = true
  var lineIndex = 0

  for i := 0; i < len(pattern); i++ {
    if lineIndex >= len(line) {
      ok = false
      break
    }
    
    if pattern[i] == '\\' && i+1 < len(pattern) {
      if pattern[i+1] == 'd' {
        if !bytes.ContainsAny([]byte{line[lineIndex]}, "0123456789") {
          ok = false
          break
        }
        i++ 
      } else if pattern[i+1] == 'w' {
        if !bytes.ContainsAny([]byte{line[lineIndex]}, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_") {
          ok = false
          break
        }
        i++
      } else {
        if line[lineIndex] != '\\' {
          ok = false
          break
        }
      }
    } else if pattern[i] == '[' {
      // Handle character class
      closeBracket := bytes.IndexByte([]byte(pattern[i:]), ']')
      if closeBracket == -1 {
        return false, fmt.Errorf("unclosed character class")
      }
      charClass := pattern[i+1 : i+closeBracket]
      if !bytes.ContainsAny([]byte{line[lineIndex]}, charClass) {
        ok = false
        break
      }
      i += closeBracket // Skip to after the closing bracket
    } else if pattern[i] != line[lineIndex] {
      ok = false
      break
    }
    lineIndex++
  }

  if ok && lineIndex < len(line) {
    ok = false // If we've matched the entire pattern but there are still characters in the line, it's not a full match
  }

	if ok {
		fmt.Println("Your word ", string(line), "contains the pattern", pattern)
	} else {
		fmt.Println("Your word ", string(line), " doesn't contains the pattern", pattern)		
	}
	

  return ok, nil
}
// func matchLine(line []byte, pattern string) (bool, error) {

// 	// Checks if the character is more than one character long, characters of 1 length not supported.
// 	// if utf8.RuneCountInString(pattern) != 1 {
// 	// 	return false, fmt.Errorf("unsupported pattern: %q", pattern)
// 	// }
		
// 	var ok bool 
// 	var lineIndex = 0

// 	for i, runeValue := range pattern {
// 		fmt.Println("values", i, runeValue)
// 		if runeValue == '\' {
// 			if pattern[i+1] == 'd' {
// 				if !bytes.ContainsAny(line[lineIndex], "0123456789") {
// 					ok = !ok
// 					break
// 				}
// 			} else if pattern[i+1] == 'w' {
// 				if !bytes.ContainsAny(line[lineIndex], "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_") {
// 					ok = !ok
// 					break
// 				}
// 			}
// 		} else if !(runeValue == line[lineIndex]){
// 			break
// 		}
// 		lineIndex += 1
// 		// else if runeValue == '[' {
// 		// 	positiveChars := strings.TrimSuffix(strings.TrimPrefix(pattern, "["), "]")
// 		// 	count = bytes.Count(positiveChars) + 1
// 		// 	ok = bytes.ContainsAny(line, positiveChars)
// 		// 	i = i + count

// 		// } 
// 	}

// 	// counter of \d and \w 
// 	// ok = bytes.ContainsAny(line, pattern)
	
// 	// if strings.HasPrefix(pattern, "[") && strings.HasSuffix(pattern, "]") {
		
// 	// 	positiveChars := strings.TrimSuffix(strings.TrimPrefix(pattern, "["), "]")
// 	// 	ok = bytes.ContainsAny(line, positiveChars)

// 	// 	if strings.HasPrefix(positiveChars, "^") {
// 	// 		ok = !ok
// 	// 	}
		
// 	// } else if strings.Contains(pattern, "\\d") {


// 	// 	ok = bytes.ContainsAny(line, "0123456789")

// 	// } else if strings.Contains(pattern, "\\w") {

// 	// 	ok = bytes.ContainsAny(line, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")
// 	// }

	
// 	if ok {
// 		fmt.Println("Your word ", string(line), "contains the pattern", pattern)
// 	} else {
// 		fmt.Println("Your word ", string(line), " doesn't contains the pattern", pattern)		
// 	}

// 	return ok, nil
// }
