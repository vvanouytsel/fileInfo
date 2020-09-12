package main

import (
	"flag"
	"fmt"
	"os"
)

// Global variables
var (
	verbose, debug bool
)

// logText outputs debug info to STDOUT.
// It is used for default output to be shown to the user.
// It is basically a wrapper around fmt.Printf
func logText(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}

// logDebug outputs debug info to STDOUT.
// It is used for extra logging useful to debug the program.
func logDebug(format string, a ...interface{}) {
	if debug {
		fmt.Fprintf(os.Stdout, format, a...)
	}
}

// logVerbose outputs verbose info to STDOUT.
// It is used to explain the user how the program achieved its result.
func logVerbose(format string, a ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stdout, format, a...)
	}
}

// logError outputs an ERROR info to STDERR and terminate the program.
func logError(format string, a ...interface{}) {
	fmt.Printf("ERROR: ")
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Printf("\n")
	os.Exit(1)
}

// printError prints the received ERROR to STDOUT and terminates the program.
// It takes an error type as input.
func printError(err error) {
	fmt.Printf("ERROR: ")
	fmt.Println(err)
	os.Exit(1)
}

// handleFlags will handle all the flags passed to the CLI.
func handleFlags() (v bool, d bool, s []string) {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "%s pathToFile \n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("Example:")
		fmt.Fprintf(flag.CommandLine.Output(), "%s /mnt \n", os.Args[0])
	}

	flag.BoolVar(&v, "v", false, "enable verbose mode")
	flag.BoolVar(&d, "d", false, "enable debug mode")
	flag.Parse()
	s = flag.Args()

	// Check if at least 1 arguments is passed
	if len(flag.Args()) < 1 {
		flag.Usage()
		logError("Please specify at least one path to a file")
	}
	return
}

// verifyFilesExist checks if all passed files/directories exist.
// A boolean and the list of files that did not exist is returned.
func verifyFilesExist(files []string) (success bool, nonExistentFiles []string) {
	success = true
	for _, f := range files {
		_, err := os.Stat(f)
		if os.IsNotExist(err) {
			logDebug("%v does not exist!\n", f)
			nonExistentFiles = append(nonExistentFiles, f)
			success = false
		}
	}
	return
}

func listPermissions(paths []string) error {
	logDebug("Listing permissions of paths: %v\n", paths)

	success, filesNonExisting := verifyFilesExist(paths)
	if !success {
		// At least one of the passed files does not exist
		logError("The following files did not exist: %v", filesNonExisting)
	}

	// TODO
	// Need to fix formatting with tabs, perhaps use tabwriter?
	// Path    Perm(text)      Perm(octat)
	// /tmp/config-err-dvmlfH  -rw-------      600
	// ↑       ↑       ↑
	// The path to your file   The permissions in text format  The permissions in octat format
	logText("Path\tPerm(text)\tPerm(octat)\n")
	for _, path := range paths {
		logDebug("Listing permission of: %v", path)
		fi, err := os.Stat(path)
		if err != nil {
			return err
		}
		logText("%s\t%s\t%o\n", path, fi.Mode(), fi.Mode())
		logVerbose("↑\t↑\t↑\n")
		logVerbose("The path to your file\tThe permissions in text format\tThe permissions in octat format\n")
	}

	return nil
}

func main() {
	var args []string
	verbose, debug, args = handleFlags()
	logDebug("Arguments received from CLI: %v\n", args)

	if err := listPermissions(args); err != nil {
		printError(err)
		return
	}
}
