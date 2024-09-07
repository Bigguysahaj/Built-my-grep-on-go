package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

// bytes: Provides functions for byte slice manipulation, like checking for character presence.
// fmt: Handles formatted I/O, similar to printf in C/C++.
// io: Contains functions for reading and writing data, similar to file handling.
// os: Provides OS-level functionality, including command-line argument handling.
// unicode/utf8: Helps with working with UTF-8 encoded strings, ensuring that characters are handled correctly.

// Usage: echo <input_text> | your_program.sh -E <pattern>

// type MatchState struct {
// 	line             []byte
// 	pattern          string
// 	lineIndex        int
// 	patternIndex     int
// 	parenExpressions []string
// 	hasCircumflex    bool
// 	hasDollar        bool
// 	hasPlus          bool
// }

// func NewMatchState(line []byte, pattern string) *MatchState {
// 	return &MatchState{
// 		line:             line,
// 		pattern:          pattern,
// 		parenExpressions: []string{"filler"},
// 		hasCircumflex:    pattern[0] == '^',
// 		hasDollar:        strings.HasSuffix(pattern, "$"),
// 	}
// }

func main() {
	// os.Args holds command-line arguments. The program expects at least three arguments: the program name, -E flag, and a pattern.
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
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


// func matchLine(line []byte, pattern string) (bool, error) {
// 	state := NewMatchState(line, pattern)

// 	if state.hasDollar {
// 		state.line = []byte(ReverseString(string(state.line)))
// 		state.pattern = ReverseString(state.pattern)
// 		fmt.Println("lines and patterns", state.line, state.pattern)
// 	}

// 	for startIndex := 0; startIndex < len(state.line); startIndex++ {
// 		if state.lineIndex >= len(state.line) {
// 			fmt.Println("segmentation fault outer loop")
// 			break
// 		}
// 		ok := MatchFromIndex(state, startIndex)
// 		if ok {
// 			fmt.Println("Your word ", string(state.line), "contains the pattern", state.pattern)
// 			return true, nil
// 		}
// 		if !ok && (state.hasCircumflex || state.hasDollar || state.hasPlus) {
// 			break
// 		}
// 	}

// 	fmt.Println("Your word ", string(state.line), " doesn't contains the pattern", state.pattern)
// 	return false, nil
// }

// // Char handlers starts from here.

// func MatchCharacter(state *MatchState) bool {
// 	switch state.pattern[state.patternIndex] {
// 	case '\\':
// 		return HandleBackslash(state)
// 	case ' ':
// 		return HandleSpace(state)
// 	case '[':
// 		return handleCharacterClass(state)
// 	case '(':
// 		return handleParentheses(state)
// 	case '+':
// 		return handlePlus(state)
// 	case '.':
// 		fmt.Println("We are doing nothing")
// 		state.lineIndex++
// 		return true
// 	default:
// 		return HandleDefault(state)
// 	}
// }

// func HandleSpace(state *MatchState) bool {
// 	if state.line[state.lineIndex] != ' ' {
// 		fmt.Println("Error on space")
// 		return false
// 	}
// 	state.lineIndex++
// 	return true
// }

// func HandleDefault(state *MatchState) bool {
// 	fmt.Printf("Let's see all %c, %c \n", state.line[state.lineIndex], state.pattern[state.patternIndex])
// 	if state.line[state.lineIndex] != state.pattern[state.patternIndex] {
// 		if state.pattern[state.patternIndex] == '?' {
// 			state.lineIndex--
// 		} else if state.patternIndex+1 < len(state.pattern) && state.pattern[state.patternIndex+1] == '?' {
// 			state.patternIndex += 2
// 			return state.pattern[state.patternIndex] == state.line[state.lineIndex]
// 		} else {
// 			fmt.Printf("Basic mismatch error %c, %c \n", state.line[state.lineIndex], state.pattern[state.patternIndex])
// 			return false
// 		}
// 	}
// 	state.lineIndex++
// 	return true
// }

// // Matching utilites starts from here.

// func MatchFromIndex(state *MatchState, startIndex int) bool {
// 	state.lineIndex = startIndex
// 	state.patternIndex = 0
// 	if state.hasCircumflex || state.hasDollar {
// 		state.patternIndex = 1
// 	}

// 	for ; state.patternIndex < len(state.pattern); state.patternIndex++ {
// 		if state.lineIndex >= len(state.line) {
// 			fmt.Println("segmentation fault inner loop")
// 			return state.patternIndex == len(state.pattern)-1
// 		}

// 		if !MatchCharacter(state) {
// 			return false
// 		}
// 	}

// 	return true
// }

// // Special Handlers start from here

// func HandleBackslash(state *MatchState) bool {
// 	if state.patternIndex+1 >= len(state.pattern) {
// 		return false
// 	}
// 	switch state.pattern[state.patternIndex+1] {
// 	case 'd':
// 		if !bytes.ContainsAny([]byte{state.line[state.lineIndex]}, "0123456789") {
// 			fmt.Println("Error on \\d")
// 			return false
// 		}
// 	case 'w':
// 		if !bytes.ContainsAny([]byte{state.line[state.lineIndex]}, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_") {
// 			fmt.Println("Error on \\w")
// 			return false
// 		}
// 	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
// 		return handleBackReference(state)
// 	default:
// 		if state.line[state.lineIndex] != '\\' {
// 			return false
// 		}
// 	}
// 	state.patternIndex++
// 	state.lineIndex++
// 	return true
// }

// func handleBackReference(state *MatchState) bool {
// 	currExpression := state.parenExpressions[int(state.pattern[state.patternIndex+1]-'0')]
// 	for state.lineIndex < len(state.line) {
// 		fmt.Println("we have in repeat", string(state.line[state.lineIndex:]), currExpression)
// 		if strings.HasSuffix(string(state.line[state.lineIndex:]), currExpression) {
// 			state.lineIndex += len(currExpression)
// 			return true
// 		}
// 		state.lineIndex++
// 	}
// 	return false
// }

// func handleCharacterClass(state *MatchState) bool {
// 	closeBracket := bytes.IndexByte([]byte(state.pattern[state.patternIndex:]), ']')
// 	if closeBracket == -1 {
// 		return false
// 	}
// 	positiveCharacters := state.pattern[state.patternIndex+1 : state.patternIndex+closeBracket]
// 	ok := bytes.ContainsAny(state.line, positiveCharacters)
// 	if bytes.ContainsAny([]byte(positiveCharacters), "^") {
// 		ok = !ok
// 	}
// 	if !ok {
// 		fmt.Println("Lost to brackets")
// 		return false
// 	}
// 	state.patternIndex += closeBracket
// 	state.lineIndex++
// 	return true
// }

// func handleParentheses(state *MatchState) bool {
// 	closeParen := bytes.IndexByte([]byte(state.pattern[state.patternIndex:]), ')')
// 	expression := state.pattern[state.patternIndex+1 : state.patternIndex+closeParen]
// 	alternatives := strings.Split(expression, "|")
// 	if len(alternatives) > 1 {
// 		return handleAlternation(state, alternatives)
// 	}
// 	return handleSimpleParentheses(state, expression)
// }

// func handleAlternation(state *MatchState, alternatives []string) bool {
// 	for _, alt := range alternatives {
// 		if strings.HasPrefix(string(state.line[state.lineIndex:]), alt) {
// 			state.lineIndex += len(alt)
// 			state.patternIndex += len(alt) + 2 // +2 for the parentheses
// 			return true
// 		}
// 	}
// 	fmt.Println("got faltered at alternation")
// 	return false
// }

// func handleSimpleParentheses(state *MatchState, expression string) bool {
// 	state.parenExpressions = append(state.parenExpressions, expression)
// 	for state.lineIndex < len(state.line) {
// 		fmt.Println("we have ", string(state.line[state.lineIndex:]), expression)
// 		if strings.HasSuffix(string(state.line[state.lineIndex:]), expression) {
// 			state.lineIndex += len(expression)
// 			state.patternIndex += len(expression) + 1
// 			return true
// 		}
// 		state.lineIndex++
// 	}
// 	return false
// }

// func handlePlus(state *MatchState) bool {
// 	state.hasPlus = true
// 	for state.line[state.lineIndex] == state.pattern[state.patternIndex-1] {
// 		state.lineIndex++
// 	}
// 	return true
// }


// func ReverseString(s string) string {
// 	runes := []rune(s)
// 	size := len(runes)
// 	for i, j := 0, size-1; i < size>>1; i, j = i+1, j-1 {
// 			runes[i], runes[j] = runes[j], runes[i]
// 	}
// 	return string(runes)
// }


func matchLine(line []byte, pattern string) (bool, error) {

	if utf8.RuneCountInString(pattern) == 0 {

		return false, fmt.Errorf("unsupported pattern: %q", pattern)

	}

	var ok bool

	var err error

	var tmp []string

	// You can use print statements as follows for debugging, they'll be visible when running tests.

	fmt.Println("Logs from your program will appear here!")

	pattern = strings.ReplaceAll(pattern, "\\d", "[0-9]")

	pattern = strings.ReplaceAll(pattern, "\\w", "[a-zA-Z0-9_]")

	fmt.Println(pattern)

	ree := regexp.MustCompile(`\(([^)]*)\)`)

	matches := ree.FindAllStringSubmatch(pattern, -1)

	if len(matches) != 0 {

		for _, match := range matches {

			fmt.Println(match)

			tmp = append(tmp, match[0])

		}

		pattern = strings.ReplaceAll(pattern, "\\1", tmp[0])

	}

	fmt.Println(tmp)

	fmt.Println(pattern)

	re, err := regexp.Compile(pattern)

	if err != nil {

		return false, err

	}

	ok = re.Match(line)

	fmt.Println(ok)

	return ok, nil

}