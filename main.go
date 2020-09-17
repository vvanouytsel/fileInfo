package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"text/tabwriter"
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

func listInfo(paths []string) (err error) {
	logDebug("Listing info of paths: %v\n", paths)

	success, filesNonExisting := verifyFilesExist(paths)
	if !success {
		// At least one of the passed files does not exist
		logError("The following files did not exist: %v", filesNonExisting)
	}

	err = writeFileInfo(paths)
	if err != nil {
		return err
	}

	err = writePermissions(paths)
	if err != nil {
		return err
	}

	return nil
}

func writeFileInfo(paths []string) (err error) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)
	fmt.Fprintln(w, "Path\tSize(bytes)\tInode\tModify")

	for _, path := range paths {
		logDebug("Listing file info of: %v\n", path)

		fi, err := os.Stat(path)
		if err != nil {
			return err
		}

		stat, ok := fi.Sys().(*syscall.Stat_t)
		if !ok {
			logError("Not a syscall.Stat_t")
		}

		fmt.Fprintf(w, "%s\t%v\t%v\t%v\n", path, fi.Size(), stat.Ino, fi.ModTime().Format("01/02/2006 - 15:04:05"))

	}
	if verbose {
		explanationSlice := []string{"The path to your file",
			"This is the size of the file in bytes",
			"This is the inode address of the file",
			"This is date since the last time the file was modified"}

		text := `
An inode is a data structure. It defines a file or a directory on the file system and is stored in the directory entry.
Inodes point to blocks that make up a file. The inode contains all the administrative data needed to read a file.
Every file’s metadata is stored in inodes in a table structure.
`

		err := explain(w, explanationSlice, text)
		if err != nil {
			return err
		}
	}

	w.Flush()
	return
}
func writePermissions(paths []string) (err error) {
	// Write the output of the file permissions
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)
	fmt.Fprintln(w, "Path\tPermissions(text)\tPermissions(binary)\tPermissions(octal)")

	for _, path := range paths {
		logDebug("Listing permission of: %v\n", path)

		fi, err := os.Stat(path)
		if err != nil {
			return err
		}

		fmt.Fprintf(w, "%s\t%s\t%b\t%o\n", path, fi.Mode().Perm(), fi.Mode().Perm(), fi.Mode().Perm())

	}
	fmt.Println()

	if verbose {
		explanationSlice := []string{"The path to your file",
			"This is the permission of the file written in text format",
			"This is the permission of the file written in binary format",
			"This is the permission of the file written in octal format"}

		text := `
On a Linux system, each file and directory is assigned access rights for the owner of the file, 
the members of a group of related users, and everybody else. Rights can be assigned to read a file, 
to write a file, and to execute a file (i.e., run the file as a program).

r(4) - Allows the contents of the directory to be listed if the x attribute is also set.
w(2) - Allows files within the directory to be created, deleted, or renamed if the x attribute is also set.
x(1) - Allows a directory to be entered (i.e. cd dir).
`

		err := explain(w, explanationSlice, text)
		if err != nil {
			return err
		}
	}

	w.Flush()
	return
}

// explain adds fancy "└>" syntax to explain the output that is listed above.
func explain(w io.Writer, stringSlice []string, text string) (err error) {
	logDebug("Explaining: %v\n", stringSlice)

	// The amount of values in the slice represents the amount of collumns
	col := len(stringSlice)

	for i := col; i > 0; i-- {
		// 1 2 3 4 ...
		fmt.Fprintf(w, strings.Repeat("|\t", i)+"\n")
		fmt.Fprintf(w, strings.Repeat("|\t", i-1))
		fmt.Fprintf(w, "└> %v\n", stringSlice[i-1])
	}

	fmt.Printf("\n%s\n", text)

	return
}

func main() {
	var args []string
	verbose, debug, args = handleFlags()
	logDebug("Arguments received from CLI: %v\n", args)
	err := listInfo(args)
	if err != nil {
		printError(err)
	}
}
