package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// BUFFERSIZE : Controls the size of the file copying buffer.
var BUFFERSIZE int64

// Filer : Implements a type to abstract recursive file operations.
type Filer struct {
	sourcePath   string
	targetPath   string
	isView       bool
	ignoredFiles []string
	visitedFiles []string
}

// visit : Visit the directory and add all found files to the files array.
func (filer *Filer) visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		// Found a directory, do not append to filenames list.
		if info.IsDir() {
			return nil
		}

		// Append the file to the filenames list.
		*files = append(*files, path)

		return nil
	}
}

func (filer *Filer) copyBuffered(src, dst string, BUFFERSIZE int64) error {
	sourceFileStat, err := os.Stat(src)

	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)

	if err != nil {
		return err
	}

	defer source.Close()

	_, err = os.Stat(dst)

	if err == nil {
		return fmt.Errorf("File %s already exists", dst)
	}

	destination, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer destination.Close()

	if err != nil {
		panic(err)
	}

	buf := make([]byte, BUFFERSIZE)

	for {
		n, err := source.Read(buf)

		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

// getDestination : Extract the destination path from the filename coming in, and take
// care of creating any subdirectories when needed.
func (filer *Filer) getDestination(filename string) string {
	saveFullPath := filer.targetPath + strings.Replace(filename, filer.sourcePath, "", -1)
	pathArray := strings.Split(saveFullPath, "\\")
	savePath := strings.Join(pathArray[:len(pathArray)-1], "\\")

	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		filer.logOutput("Creating new directory: " + savePath)
		os.MkdirAll(savePath, os.ModePerm)
	}

	return saveFullPath
}

func (filer *Filer) getIgnoredFiles() {
	// Read the .obfsignore file to take a list of files to ignore for obfuscation.
	file, err := os.Open(".obfsignore")

	if err != nil {
		log.Fatal(err)
	}

	// Make sure to close the file once we're done with it.
	defer file.Close()

	// We use a scanner here, because we really don't care to have an actual new line character in our
	// output.
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		filer.ignoredFiles = append(filer.ignoredFiles, scanner.Text())
	}

	// Do some error handling, because we like it.
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// stringInSlice : Checks if an instance exists in an array. This behavior really does not belong
// to the Filer, but it's better than duplicating it all over the place. Should be extracted to
// some helper object.
func (filer *Filer) stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}

	return false
}

func (filer *Filer) logOutput(str string) {
	logline := time.Now().Format(time.RFC1123) + " : " + str

	if os.Args[3] != "--debug" {
		fmt.Println(logline)
	} else {
		f, err := os.OpenFile("logfile.txt", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if _, err = f.WriteString(logline + "\n"); err != nil {
			panic(err)
		}
	}
}

func (filer *Filer) shouldObfuscate(filename string) bool {
	extension := string(filename[len(filename)-3:])

	// We're only obfuscating php files, so check for the extension.
	if extension != "php" {
		return false
	}

	// Not the greatest way to deal with this, but nevertheless a workable solition to check for the
	// files that should be ignored for obfuscation.
	for _, ignorePath := range filer.ignoredFiles {
		if strings.Contains(filename, ignorePath) {
			return false
		}
	}

	// All good, let us obfuscate.
	return true
}

func (filer *Filer) doObfuscate(filename string, file *os.File) {
	filer.logOutput("Obfuscating file: " + filename)
	parser := new(Parser)
	parser.parse(file, filer)
	newfile := strings.Join(parser.obfuscatedLines, "\n")

	savelocation := filer.getDestination(filename)

	filer.logOutput("Writing new file: " + savelocation)
	fout, err := os.Create(savelocation)

	if err != nil {
		panic(err)
	}

	fout.Sync()
	w := bufio.NewWriter(fout)
	w.WriteString(newfile)
	w.Flush()
}

func (filer *Filer) doCopy(filename string) {
	// Set up the buffersize for copying files.
	BUFFERSIZE, err := strconv.ParseInt("1000000", 10, 64)

	if err != nil {
		log.Fatal("Invalid buffer size!")
		return
	}

	savelocation := filer.getDestination(filename)

	filer.logOutput("Copying " + filename + " to " + savelocation)

	err = filer.copyBuffered(filename, savelocation, BUFFERSIZE)

	if err != nil {
		log.Fatal("File copy failed!")
	}
}

// walkWithMe : Loop over the sourcePath passed into the instace, instantiate parser, and write outpupt file.
func (filer *Filer) walkWithMe() {
	// Let's fill the ignored files array so we have something to keep track whether to obfuscate or just copy a file.
	filer.getIgnoredFiles()

	// List of filenames that we will obfuscate.
	var filenames []string

	// Fill up the filenames array with all the files found in the root path, which we'll walk recursively.
	err := filepath.Walk(filer.sourcePath, filer.visit(&filenames))

	// Fire in the disco.
	if err != nil {
		panic(err)
	}

	// Looping over each found file.
	for _, filename := range filenames {
		file, err := os.Open(filename)

		// Fire in the, Taco-Bell.
		if err != nil {
			log.Fatal(err)
		}

		// Make sure to close the file once we're done with it.
		defer file.Close()

		if filer.shouldObfuscate(filename) {
			// This file should be obfuscated, because it is not in the ignored list.
			filer.doObfuscate(filename, file)
		} else {
			// No obfuscation needed, just copy the file.
			filer.doCopy(filename)
		}
	}
}
