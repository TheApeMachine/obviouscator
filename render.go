package main

// Render : Handles the render context of the parser.
type Render struct{}

func (render Render) getNextState(parser *Parser) string {
	if parser.currentToken == "Array" || parser.currentToken == "array" {
		parser.context = ARRAY
	}

	return parser.currentCharacter
}
