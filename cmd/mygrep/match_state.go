package main

import "strings"

type matchState struct {
	line             []byte
	pattern          string
	lineIndex        int
	patternIndex     int
	parenExpressions []string
	hasCircumflex    bool
	hasDollar        bool
	hasPlus          bool
}

func NewMatchState(line []byte, pattern string) *matchState {
	return &matchState{
		line:             line,
		pattern:          pattern,
		parenExpressions: []string{"filler"},
		hasCircumflex:    pattern[0] == '^',
		hasDollar:        strings.HasSuffix(pattern, "$"),
	}
}