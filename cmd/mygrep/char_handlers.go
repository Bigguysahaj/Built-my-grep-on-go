package main

import "fmt"

func MatchCharacter(state *MatchState) bool {
	switch state.pattern[state.patternIndex] {
	case '\\':
		return HandleBackslash(state)
	case ' ':
		return HandleSpace(state)
	case '[':
		return handleCharacterClass(state)
	case '(':
		return handleParentheses(state)
	case '+':
		return handlePlus(state)
	case '.':
		fmt.Println("We are doing nothing")
		state.lineIndex++
		return true
	default:
		return HandleDefault(state)
	}
}

func HandleSpace(state *MatchState) bool {
	if state.line[state.lineIndex] != ' ' {
		fmt.Println("Error on space")
		return false
	}
	state.lineIndex++
	return true
}

func HandleDefault(state *MatchState) bool {
	fmt.Printf("Let's see all %c, %c \n", state.line[state.lineIndex], state.pattern[state.patternIndex])
	if state.line[state.lineIndex] != state.pattern[state.patternIndex] {
		if state.pattern[state.patternIndex] == '?' {
			state.lineIndex--
		} else if state.patternIndex+1 < len(state.pattern) && state.pattern[state.patternIndex+1] == '?' {
			state.patternIndex += 2
			return state.pattern[state.patternIndex] == state.line[state.lineIndex]
		} else {
			fmt.Printf("Basic mismatch error %c, %c \n", state.line[state.lineIndex], state.pattern[state.patternIndex])
			return false
		}
	}
	state.lineIndex++
	return true
}