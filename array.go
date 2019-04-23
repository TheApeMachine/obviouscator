package main

// Array : Handles the render array context of the parser.
type Array struct{}

func (array Array) isKeyBoundry(parser *Parser) bool {
	if parser.currentCharacter == "'" || parser.currentCharacter == "\"" {
		if parser.previousCharacter == " " || parser.previousCharacter == "," {
			return true
		}
	}

	return false
}

func (array Array) getNextState(parser *Parser) string {
	// TODO: We need to keep track of open/close levels of brackets, otherwise we cannot do nested arrays.

	switch parser.context {
	case ARRAY:
		if parser.currentCharacter == "'" || parser.currentCharacter == "\"" {
			parser.previousContext = ARRAY
			parser.context = VARIABLE
		} else if parser.currentCharacter == ")" && parser.previousCharacter == ")" {
			parser.context = START
		}

	case DICTVALUE:
		if array.isKeyBoundry(parser) {
			parser.context = ARRAY
		} else if parser.currentCharacter == "$" {
			parser.currentVar = ""
			parser.previousContext = ARRAY
			parser.context = VARIABLE
		}
	}

	return parser.currentCharacter
}
