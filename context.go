package main

// Context : Holds methods global to each context.
type Context struct {
	handler interface {
		getNextState(parser *Parser) string
	}
}

func (context *Context) getNextState(parser *Parser) string {
	newline := context.handler.getNextState(parser)

	if newline != "\n" && newline != "\r" {
		return newline
	}

	return ""
}
