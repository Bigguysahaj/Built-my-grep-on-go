package main

import "fmt"

func matchFromIndex(state *matchState, startIndex int) bool {
	state.lineIndex = startIndex
	state.patternIndex = 0
	if state.hasCircumflex || state.hasDollar {
		state.patternIndex = 1
	}

	for ; state.patternIndex < len(state.pattern); state.patternIndex++ {
		if state.lineIndex >= len(state.line) {
			fmt.Println("segmentation fault inner loop")
			return state.patternIndex == len(state.pattern)-1
		}

		if !matchCharacter(state) {
			return false
		}
	}

	return true
}