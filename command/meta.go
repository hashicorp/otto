package command

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/otto"
	"github.com/hashicorp/otto/ui"
	"github.com/mitchellh/cli"
)

const (
	// DefaultAppfile is the default filename for the Appfile
	DefaultAppfile = "Appfile"

	// DefaultOutputDir is the default filename for the output directory
	DefaultOutputDir                = ".otto"
	DefaultOutputDirCompiledAppfile = "appfile"
	DefaultOutputDirCompiledData    = "compiled"

	// DefaultDataDir is the default directory for the directory
	// data if a directory in the Appfile isn't specified.
	DefaultDataDir = "otto-data"

	// EnvAppFile is the environment variable to point to an appfile.
	EnvAppFile = "OTTO_APPFILE"
)

// FlagSetFlags is an enum to define what flags are present in the
// default FlagSet returned by Meta.FlagSet
type FlagSetFlags uint

const (
	FlagSetNone FlagSetFlags = 0
)

// Meta are the meta-options that are available on all or most commands.
type Meta struct {
	CoreConfig *otto.CoreConfig
	Ui         cli.Ui
}

// Appfile loads the compiled Appfile. If the Appfile isn't compiled yet,
// then an error will be returned.
func (m *Meta) Appfile() (*appfile.Compiled, error) {
	// Find the root directory
	rootDir, err := m.RootDir()
	if err != nil {
		return nil, err
	}

	return appfile.LoadCompiled(filepath.Join(
		rootDir, DefaultOutputDir, DefaultOutputDirCompiledAppfile))
}

// Core returns the core for the given Appfile. The file where the
// Appfile was loaded from should be set in appfile.File.Path. This
// root appfile path will be used as the default output directory
// for Otto.
func (m *Meta) Core(f *appfile.Compiled) (*otto.Core, error) {
	dir, err := m.Directory(f.File)
	if err != nil {
		return nil, err
	}

	rootDir, err := m.RootDir()
	if err != nil {
		return nil, err
	}

	config := *m.CoreConfig
	config.Appfile = f
	config.Directory = dir
	config.OutputDir = filepath.Join(
		rootDir, DefaultOutputDir, DefaultOutputDirCompiledData)
	config.Ui = m.OttoUi()

	return otto.NewCore(&config)
}

// RootDir finds the "root" directory. This is the working directory of
// the Appfile and Otto itself. To find the root directory, we traverse
// upwards until we find the ".otto" directory and assume that is where
// it is.
func (m *Meta) RootDir() (string, error) {
	// First, get our current directory
	current, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Traverse upwards until we find the directory. We also protect this
	// loop with a basic infinite loop guard.
	i := 0
	prev := ""
	for prev != current && i < 1000 {
		if _, err := os.Stat(filepath.Join(current, DefaultOutputDir)); err == nil {
			// Found it
			return current, nil
		}

		prev = current
		current = filepath.Dir(current)
	}

	return "", fmt.Errorf(
		"Otto doesn't appear to have compiled!\n\n" +
			"Run `otto compile` in the directory with the Appfile or\n" +
			"with the `-appfile` flag in order to compile the files for\n" +
			"developing, building, and deploying your application.")
}

// Directory returns the Otto directory backend for the given
// Appfile. If no directory backend is specified, a local folder
// will be used.
func (m *Meta) Directory(f *appfile.File) (directory.Backend, error) {
	// TODO: Appfile can't specify directory configuration

	return &directory.FolderBackend{Dir: DefaultDataDir}, nil
}

// FlagSet returns a FlagSet with the common flags that every
// command implements. The exact behavior of FlagSet can be configured
// using the flags as the second parameter.
func (m *Meta) FlagSet(n string, fs FlagSetFlags) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)

	// Create an io.Writer that writes to our Ui properly for errors.
	// This is kind of a hack, but it does the job. Basically: create
	// a pipe, use a scanner to break it into lines, and output each line
	// to the UI. Do this forever.
	errR, errW := io.Pipe()
	errScanner := bufio.NewScanner(errR)
	go func() {
		for errScanner.Scan() {
			m.Ui.Error(errScanner.Text())
		}
	}()
	f.SetOutput(errW)

	return f
}

// OttoUi returns the ui.Ui object.
func (m *Meta) OttoUi() ui.Ui {
	return NewUi(m.Ui)
}
