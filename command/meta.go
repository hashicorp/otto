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
	"github.com/mitchellh/go-homedir"
)

const (
	// DefaultAppfile is the default filename for the Appfile
	DefaultAppfile = "Appfile"

	// DefaultLocalDataDir is the default path to the local data
	// directory.
	DefaultLocalDataDir         = "~/.otto.d"
	DefaultLocalDataDetectorDir = "detect"

	// DefaultOutputDir is the default filename for the output directory
	DefaultOutputDir                = ".otto"
	DefaultOutputDirCompiledAppfile = "appfile"
	DefaultOutputDirCompiledData    = "compiled"
	DefaultOutputDirLocalData       = "data"

	// DefaultDataDir is the default directory for the directory
	// data if a directory in the Appfile isn't specified.
	DefaultDataDir = "otto-data"
)

var (
	// AltAppfiles is the list of alternative names for an Appfile that Otto can
	// detect and load automatically
	AltAppfiles = []string{"appfile.hcl"}
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
	startDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	rootDir, err := m.RootDir(startDir)
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
	if f.File == nil || f.File.Path == "" {
		return nil, fmt.Errorf("Could not determine Appfile dir")
	}

	rootDir, err := m.RootDir(filepath.Dir(f.File.Path))
	if err != nil {
		return nil, err
	}

	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}

	dataDir, err := m.DataDir()
	if err != nil {
		return nil, err
	}

	config := *m.CoreConfig
	config.Appfile = f
	config.DataDir = dataDir
	config.LocalDir = filepath.Join(
		rootDir, DefaultOutputDir, DefaultOutputDirLocalData)
	config.CompileDir = filepath.Join(
		rootDir, DefaultOutputDir, DefaultOutputDirCompiledData)
	config.Ui = m.OttoUi()

	config.Directory, err = m.Directory(&config)
	if err != nil {
		return nil, err
	}

	return otto.NewCore(&config)
}

// DataDir returns the user-local data directory for Otto.
func (m *Meta) DataDir() (string, error) {
	return homedir.Expand(DefaultLocalDataDir)
}

// RootDir finds the "root" directory. This is the working directory of
// the Appfile and Otto itself. To find the root directory, we traverse
// upwards until we find the ".otto" directory and assume that is where
// it is.
func (m *Meta) RootDir(startDir string) (string, error) {
	current := startDir

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
		i++
	}

	return "", fmt.Errorf(
		"Otto doesn't appear to have compiled your Appfile yet!\n\n" +
			"Run `otto compile` in the directory with the Appfile or\n" +
			"with the `-appfile` flag in order to compile the files for\n" +
			"developing, building, and deploying your application.\n\n" +
			"Once the Appfile is compiled, you can run `otto` in any\n" +
			"subdirectory.")
}

// Directory returns the Otto directory backend for the given
// Appfile. If no directory backend is specified, a local folder
// will be used.
func (m *Meta) Directory(config *otto.CoreConfig) (directory.Backend, error) {
	return &directory.BoltBackend{
		Dir: filepath.Join(config.DataDir, "directory"),
	}, nil
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

// confirmDestroy is a little helper that will ask the user to confirm a
// destroy action using the provided msg, unless -force is included in args it
// returns true if the destroy should be considered confirmed, and false if
// the destroy should be aborted.
func (m *Meta) confirmDestroy(msg string, args []string) bool {
	destroyForce := false
	for _, arg := range args {
		if arg == "-force" {
			destroyForce = true
		}
	}

	if !destroyForce {
		v, err := m.OttoUi().Input(&ui.InputOpts{
			Id:    "destroy",
			Query: "Do you really want to destroy?",
			Description: fmt.Sprintf("%s\n"+
				"There is no undo. Only 'yes' will be accepted to confirm.", msg),
		})
		if err != nil {
			m.Ui.Error(fmt.Sprintf("Error asking for confirmation: %s", err))
			return false
		}
		if v != "yes" {
			m.Ui.Output("Destroy cancelled.")
			return false
		}
	}

	return true
}
