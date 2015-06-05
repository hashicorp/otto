package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/mitchellh/cli"
	"github.com/mitchellh/panicwrap"
	"github.com/mitchellh/prefixedio"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	// Set a custom panicwrap cookie key and value. Since we're executing
	// other panicwrap executables, the default cookie key/value cause
	// weird errors if we don't change them.
	var wrapConfig panicwrap.WrapConfig
	wrapConfig.CookieKey = "OTTO_PANICWRAP_COOKIE"
	wrapConfig.CookieValue = fmt.Sprintf(
		"otto-%s-%s-%s", Version, VersionPrerelease, GitCommit)

	if !panicwrap.Wrapped(&wrapConfig) {
		// Determine where logs should go in general (requested by the user)
		logWriter, err := logOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't setup log output: %s", err)
			return 1
		}
		if logWriter == nil {
			logWriter = ioutil.Discard
		}

		// We always send logs to a temporary file that we use in case
		// there is a panic. Otherwise, we delete it.
		logTempFile, err := ioutil.TempFile("", "otto-log")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't setup logging tempfile: %s", err)
			return 1
		}
		defer os.Remove(logTempFile.Name())
		defer logTempFile.Close()

		// Tell the logger to log to this file
		os.Setenv(EnvLog, "")
		os.Setenv(EnvLogFile, "")

		// Setup the prefixed readers that send data properly to
		// stdout/stderr.
		doneCh := make(chan struct{})
		outR, outW := io.Pipe()
		go copyOutput(outR, doneCh)

		// Create the configuration for panicwrap and wrap our executable
		wrapConfig.Handler = panicHandler(logTempFile)
		wrapConfig.Writer = io.MultiWriter(logTempFile, logWriter)
		wrapConfig.Stdout = outW
		exitStatus, err := panicwrap.Wrap(&wrapConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't start Otto: %s", err)
			return 1
		}

		// If >= 0, we're the parent, so just exit
		if exitStatus >= 0 {
			// Close the stdout writer so that our copy process can finish
			outW.Close()

			// Wait for the output copying to finish
			<-doneCh

			return exitStatus
		}

		// We're the child, so just close the tempfile we made in order to
		// save file handles since the tempfile is only used by the parent.
		logTempFile.Close()
	}

	// Call the real main
	return wrappedMain()
}

func wrappedMain() int {
	log.SetOutput(os.Stderr)
	log.Printf(
		"[INFO] Otto version: %s %s %s",
		Version, VersionPrerelease, GitCommit)

	// Load the configuration
	config := BuiltinConfig

	// Run checkpoint
	go runCheckpoint(&config)

	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "-version" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	cli := &cli.CLI{
		Args:       args,
		Commands:   Commands,
		HelpFunc:   cli.BasicHelpFunc("otto"),
		HelpWriter: os.Stdout,
	}

	exitCode, err := cli.Run()
	if err != nil {
		Ui.Error(fmt.Sprintf("Error executing CLI: %s", err.Error()))
		return 1
	}

	return exitCode
}

// copyOutput uses output prefixes to determine whether data on stdout
// should go to stdout or stderr. This is due to panicwrap using stderr
// as the log and error channel.
func copyOutput(r io.Reader, doneCh chan<- struct{}) {
	defer close(doneCh)

	pr, err := prefixedio.NewReader(r)
	if err != nil {
		panic(err)
	}

	stderrR, err := pr.Prefix(ErrorPrefix)
	if err != nil {
		panic(err)
	}
	stdoutR, err := pr.Prefix(OutputPrefix)
	if err != nil {
		panic(err)
	}
	defaultR, err := pr.Prefix("")
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		io.Copy(os.Stderr, stderrR)
	}()
	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, stdoutR)
	}()
	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, defaultR)
	}()

	wg.Wait()
}
