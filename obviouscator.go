package main

import "os"

func main() {
	// Get the root directories for source and target from the command line arguments.
	sourceRootDir := os.Args[1]
	targetRootDir := os.Args[2]

	// Let's make sure our taget directory is completely empty, so we don't get any existing file errors later.
	os.RemoveAll(targetRootDir)

	// Setup a new Filer, which will handle traversing directories, copying files, and calling the parser when needed.
	filer := Filer{sourceRootDir, targetRootDir, false, make([]string, 0), make([]string, 0)}

	// Start the top-level parser, which will make calls to a new instance of Filer to obfuscate or copy
	// related (nested) documents.
	filer.walkWithMe()
}
