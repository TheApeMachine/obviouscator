package main

// Start : Handles the starting context of the parser.
type Start struct{}

// getNextSate : Implementation of the common interface of a context, returns a boolean value
// that determines if the current context should print, which concatenates to the new lines being
// generated.
func (start Start) getNextState(parser *Parser) string {
	if parser.filer.stringInSlice(parser.currentToken, parser.keywords) {
		if parser.newline == "<script" || parser.newline == "</script>" {
			if parser.context == JAVASCRIPT {
				parser.context = START
			} else {
				parser.context = JAVASCRIPT
			}
		} else if parser.currentToken == "render" {
			// We're possibly in a render statement of a controller, we should convert the obfuscate the injected
			// variable dictionary.
			parser.context = RENDER
		}
	} else if parser.currentCharacter == "/" {
		parser.context = COMMENTSTART
		return ""
	} else if parser.currentCharacter == "$" {
		parser.currentVar = ""
		parser.context = VARIABLE
	} else if parser.currentCharacter == "{" {
		parser.currentLevel++
	} else if parser.currentCharacter == "}" {
		parser.currentLevel--
	} else if parser.currentCharacter == "'" {
		parser.context = SINGLESTRING
	} else if parser.currentCharacter == "\"" {
		parser.context = DOUBLESTRING
	}

	parser.previousContext = START
	return parser.currentCharacter
}
