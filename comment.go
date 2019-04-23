package main

// Comment : Handles the comment context of the parser.
type Comment struct{}

func (comment Comment) getNextState(parser *Parser) string {
	switch parser.context {
	case COMMENTSTART:
		if parser.currentCharacter == "/" {
			// We have previously seen a slash, second slash confirms we are in a line comment.
			parser.context = LINECOMMENT
		} else if parser.currentCharacter == "*" {
			// We have previously seen a slash, star confirms we are in a block comment.
			parser.context = BLOCKCOMMENT
		} else {
			// Handle false positive comment, whatever code may have slashes in it like paths.
			parser.context = START
			return parser.previousCharacter + parser.currentCharacter
		}

	case LINECOMMENT:
		if parser.currentCharacter == "\n" {
			parser.context = START
		} else if parser.previousCharacter == "?" && parser.currentCharacter == ">" {
			// Yeah this is getting hacky again, but there are cases with comments inside views that use php tags
			// to be able to use php comment style, instead of html comment style.
			parser.context = START
			return parser.previousCharacter + parser.currentCharacter
		}

	case BLOCKCOMMENT:
		if parser.currentCharacter == "*" {
			parser.context = COMMENTEND
		}

	case COMMENTEND:
		if parser.currentCharacter == "/" && parser.previousCharacter == "*" {
			parser.context = START
		}
	}

	return ""
}
