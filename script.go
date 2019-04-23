package main

// Script : Handles the script context of the parser.
type Script struct{}

func (script Script) getNextState(parser *Parser) string {
	return parser.currentCharacter
}
