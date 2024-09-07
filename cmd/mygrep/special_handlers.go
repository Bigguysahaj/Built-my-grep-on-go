package main

import (
	"bytes"
	"fmt"
	"strings"
)

func HandleBackslash(state *MatchState) bool {
	if state.patternIndex+1 >= len(state.pattern) {
		return false
	}
	switch state.pattern[state.patternIndex+1] {
	case 'd':
		if !bytes.ContainsAny([]byte{state.line[state.lineIndex]}, "0123456789") {
			fmt.Println("Error on \\d")
			return false
		}
	case 'w':
		if !bytes.ContainsAny([]byte{state.line[state.lineIndex]}, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_") {
			fmt.Println("Error on \\w")
			return false
		}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return handleBackReference(state)
	default:
		if state.line[state.lineIndex] != '\\' {
			return false
		}
	}
	state.patternIndex++
	state.lineIndex++
	return true
}

func handleBackReference(state *MatchState) bool {
	currExpression := state.parenExpressions[int(state.pattern[state.patternIndex+1]-'0')]
	for state.lineIndex < len(state.line) {
		fmt.Println("we have in repeat", string(state.line[state.lineIndex:]), currExpression)
		if strings.HasSuffix(string(state.line[state.lineIndex:]), currExpression) {
			state.lineIndex += len(currExpression)
			return true
		}
		state.lineIndex++
	}
	return false
}

func handleCharacterClass(state *MatchState) bool {
	closeBracket := bytes.IndexByte([]byte(state.pattern[state.patternIndex:]), ']')
	if closeBracket == -1 {
		return false
	}
	positiveCharacters := state.pattern[state.patternIndex+1 : state.patternIndex+closeBracket]
	ok := bytes.ContainsAny(state.line, positiveCharacters)
	if bytes.ContainsAny([]byte(positiveCharacters), "^") {
		ok = !ok
	}
	if !ok {
		fmt.Println("Lost to brackets")
		return false
	}
	state.patternIndex += closeBracket
	state.lineIndex++
	return true
}

func handleParentheses(state *MatchState) bool {
	closeParen := bytes.IndexByte([]byte(state.pattern[state.patternIndex:]), ')')
	expression := state.pattern[state.patternIndex+1 : state.patternIndex+closeParen]
	alternatives := strings.Split(expression, "|")
	if len(alternatives) > 1 {
		return handleAlternation(state, alternatives)
	}
	return handleSimpleParentheses(state, expression)
}

func handleAlternation(state *MatchState, alternatives []string) bool {
	for _, alt := range alternatives {
		if strings.HasPrefix(string(state.line[state.lineIndex:]), alt) {
			state.lineIndex += len(alt)
			state.patternIndex += len(alt) + 2 // +2 for the parentheses
			return true
		}
	}
	fmt.Println("got faltered at alternation")
	return false
}

func handleSimpleParentheses(state *MatchState, expression string) bool {
	state.parenExpressions = append(state.parenExpressions, expression)
	for state.lineIndex < len(state.line) {
		fmt.Println("we have ", string(state.line[state.lineIndex:]), expression)
		if strings.HasSuffix(string(state.line[state.lineIndex:]), expression) {
			state.lineIndex += len(expression)
			state.patternIndex += len(expression) + 1
			return true
		}
		state.lineIndex++
	}
	return false
}

func handlePlus(state *MatchState) bool {
	state.hasPlus = true
	for state.line[state.lineIndex] == state.pattern[state.patternIndex-1] {
		state.lineIndex++
	}
	return true
}