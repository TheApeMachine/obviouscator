package main

// Class : Handles the class context of the parser.
type Class struct{}

func (class Class) getNextState(parser *Parser) string {
	parser.context = START
	return parser.currentCharacter
}
