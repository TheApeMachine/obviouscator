package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io"
	"strings"
)

// ContextState : State machine that holds the current context state of our parser.
type ContextState int

const (
	// START : Currently no defined state.
	START ContextState = 0 + iota
	// COMMENTSTART : Possibly found a comment.
	COMMENTSTART
	// LINECOMMENT : Found a single line comment.
	LINECOMMENT
	// BLOCKCOMMENT : Found a block comment.
	BLOCKCOMMENT
	// COMMENTEND : Possibly found the end of a comment.
	COMMENTEND
	// VARIABLE : Found a variable, to be obfuscated.
	VARIABLE
	// SINGLESTRING : Found a string, do not try to parse comments.
	SINGLESTRING
	// DOUBLESTRING : Found a string, do not try to parse comments.
	DOUBLESTRING
	// RENDER : In a render statement, need to obfuscate the view.
	RENDER
	// CLASS : Found a class, need to check if we're in a controller and get the view details.
	CLASS
	// JAVASCRIPT : Found inline javascript, do not recognize $ as a variable.
	JAVASCRIPT
	// ARRAY : Inside a render array to set exposed variables to the view.
	ARRAY
	// DICTKEY : Dealing with a dictionary key
	DICTKEY
	// DICTVALUE : Dealing with a dictionary value
	DICTVALUE
)

// ObfuscatedVariable : Container to store variable references.
type ObfuscatedVariable struct {
	level          int
	originalName   string
	obfuscatedName string
}

// Parser : Holds all data and methods related to parsing the file.
type Parser struct {
	context           ContextState
	previousContext   ContextState
	currentContext    Context
	filer             *Filer
	currentCharacter  string
	previousCharacter string
	currentVar        string
	currentLevel      int
	currentToken      string
	previousToken     string
	newVariable       string
	newline           string
	delimiters        []string
	keywords          []string
	variables         []ObfuscatedVariable
	obfuscatedLines   []string
	superGlobals      []string
}

func (parser *Parser) obfuscate(varName string) string {
	algorithm := fnv.New32a()
	algorithm.Write([]byte(varName))
	return "AX" + fmt.Sprint(algorithm.Sum32())
}

// setVariable : Figures out what to set the new variable to.
func (parser *Parser) setVariable(newVar string) string {
	// Do not obfuscate PHP super-globals.
	if parser.filer.stringInSlice(newVar, parser.superGlobals) {
		return newVar
	}

	// If we already have a record of this variable, use that instance's name.
	for _, v := range parser.variables {
		if v.originalName == newVar {
			return v.obfuscatedName
		}
	}

	// Finally, we want to obfuscate the variable, and store it in the known instances.
	var newVariable = ObfuscatedVariable{parser.currentLevel, newVar, parser.obfuscate(newVar)}
	parser.variables = append(parser.variables, newVariable)

	return newVariable.obfuscatedName
}

func (parser *Parser) parse(file io.Reader, filer *Filer) {
	reader := bufio.NewReader(file)
	parser.filer = filer
	parser.context = START
	parser.keywords = []string{"script", "render"}
	parser.delimiters = []string{" ", ",", ".", "=", "(", ")", "[", "]", "{", "}", "-", ";", "+", "'", "\"", "\r", "\n", "\\", "<", ">"}
	parser.superGlobals = []string{"_COOKIE", "this", "_POST", "_GET", "_SERVER", "_REQUEST", "_FILES", "GLOBALS", "_ENV", "_SESSION", "layout", "funct_param", "err_msg_delete_funct", "is_only_view", "HTTP_POST_VARS", "HTTP_GET_VARS", "HTTP_COOKIE", "HTTP_SESSION", "HTTP_RAW_POST_DATA", "HTTP_RESPONSE_HEADER", "argc", "argv", "php_errormsg", "content", "HTTP_SERVER_VARS", "breadcrumbs", "global_settings_list", "yii", "config", "blank", "menu", "is_partial_view", "ws_client", "soap_params", "client", "response"}
	parser.currentLevel = 0
	parser.previousContext = parser.context

	for {
		line, err := reader.ReadString('\n')

		// Holds the newly (obfuscated) generated line of code.
		parser.newline = ""

		// Looping over each character in the file.
		for _, r := range line {
			// Store the range pointer in a single character string.
			parser.currentCharacter = string(r)

			// The current context determines the behavior of the parser.
			switch parser.context {
			case START:
				parser.currentContext = Context{Start{}}

			case SINGLESTRING, DOUBLESTRING:
				parser.currentContext = Context{String{}}

			case VARIABLE:
				parser.currentContext = Context{Variable{}}

			case COMMENTSTART, LINECOMMENT, BLOCKCOMMENT, COMMENTEND:
				parser.currentContext = Context{Comment{}}

			case CLASS:
				parser.currentContext = Context{Class{}}

			case JAVASCRIPT:
				parser.currentContext = Context{Script{}}

			case RENDER:
				parser.currentContext = Context{Render{}}

			case ARRAY, DICTKEY, DICTVALUE:
				parser.currentContext = Context{Array{}}
			}

			// Call the commonly implemented method on the current context to get the print response
			// from the current ContextState.
			parser.newline += parser.currentContext.getNextState(parser)

			// Let's break stuff up into tokens.
			if filer.stringInSlice(parser.currentCharacter, parser.delimiters) {
				if parser.currentToken != "" {
					parser.previousToken = parser.currentToken
				}

				parser.currentToken = ""
			} else {
				parser.currentToken += parser.currentCharacter
			}

			// Store the previous character so we have a reference for it later when we need it.
			parser.previousCharacter = parser.currentCharacter
		}

		// Since we don't concatenate newline, we'll end up with a lot of empty lines, so let's
		// check for these and not add them to the output array.
		if len(strings.Fields(parser.newline)) > 0 {
			parser.obfuscatedLines = append(parser.obfuscatedLines, strings.TrimRight(parser.newline, " "))
		}

		// End of file reached, break out of the loop.
		if err == io.EOF {
			break
		}
	}

}
