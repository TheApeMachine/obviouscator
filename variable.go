package main

// Variable : Handles the variable context of the parser.
type Variable struct{}

func (variable Variable) getNextState(parser *Parser) string {
	if parser.filer.stringInSlice(parser.currentCharacter, parser.delimiters) && !parser.filer.stringInSlice(parser.previousToken, parser.superGlobals) {
		// Found a delimiter, we now have the complete variable name to obfuscate.
		parser.newVariable = parser.setVariable(parser.currentVar)
		parser.context = parser.previousContext
		newVariable := parser.newVariable
		parser.newVariable = ""
		parser.currentVar = ""
		return newVariable + parser.currentCharacter
	}

	// Keep building up the variable name until we meet a delimiter, signaling the completion of the name.
	parser.currentVar += parser.currentCharacter
	return ""
}
