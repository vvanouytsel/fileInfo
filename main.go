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
	if debug {
		fmt.Fprintf(os.Stdout, format, a...)
	}
}

// logError outputs an ERROR info to STDERR and terminate the program.
func logError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
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

func listPermissions(paths []string) {
	logDebug("Listing permissions of paths: %v\n", paths)
	for _, path := range paths {
		logDebug("Listing permission of: %v", path)
		fmt.Printf(path + "\n")
	}
}

func main() {
	var args []string
	verbose, debug, args = handleFlags()
	logDebug("Arguments received from CLI: %v\n", args)

	listPermissions(args)
}
