package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// Global variables
var (
	verbose, debug bool
)

// log will output text to STDOUT.
// It is used for text which always needs to be shown to the user.
func logF(text string) {
	log.Printf("%v\n", text)
}

// logDebug will output debug text to STDOUT.
// It is used for extra logging useful to debug the program.
func logDebug(text string) {
	if debug {
		log.Printf("DEBUG: %v\n", text)
	}
}

// logVerbose will output verbose text to STDOUT.
// It is used to explain the user how the program achieved its result.
func logVerbose(text string) {
	if verbose {
		log.Printf("V: %v\n", text)
	}
}

// logError will output an ERROR text to STDOUT.
func logError(text string) {
	log.Fatalf("ERROR: %v\n", text)
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

	// Check if arguments is only one
	if len(flag.Args()) != 1 {
		flag.Usage()
		logError("Please specify one argument")
	}
	return
}

func main() {
	var args []string
	verbose, debug, args = handleFlags()

	logDebug("Arguments passed: " + strings.Join(args, ","))
	logF("This is always shown")
	logVerbose("This is shown with -v flag")
	logDebug("This is shown with -d flag")

}
