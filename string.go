package main

// String : Handles the string context of the parser.
type String struct{}

func (string String) getNextState(parser *Parser) string {
	switch parser.context {
	case SINGLESTRING:
		if parser.currentCharacter == "'" {
			parser.context = START
		} else if parser.currentCharacter == "$" {
			parser.currentVar = ""
			parser.previousContext = START
			parser.context = VARIABLE
		}

	case DOUBLESTRING:
		// Make sure not to close the string if there is an escaped double quote in it.
		if parser.currentCharacter == "\"" && parser.previousCharacter != "\\" {
			parser.context = START
		} else if parser.currentCharacter == "$" {
			parser.currentVar = ""
			parser.previousContext = START
			parser.context = VARIABLE
		}
	}

	return parser.currentCharacter
}
