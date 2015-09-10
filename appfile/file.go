package appfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/helper/oneline"
	"github.com/hashicorp/otto/helper/uuid"
)

const (
	IDFile = ".ottoid"
)

// File is the structure of a single Appfile.
type File struct {
	// ID is a unique UUID that represents this file. It is generated the
	// first time on compile. This will be blank until the Appfile is
	// compiled with Compile.
	ID string

	// Path is the path to the root file that was loaded. This might be
	// empty if the appfile was parsed from an io.Reader.
	Path string

	Application    *Application
	Project        *Project
	Infrastructure []*Infrastructure
	Customization  *CustomizationSet
}

// Application is the structure of an application definition.
type Application struct {
	Name         string
	Type         string
	Dependencies []*Dependency `mapstructure:"dependency"`
}

// Customization is the structure of customization stanzas within
// the Appfile.
type Customization struct {
	Type   string
	Config map[string]interface{}
}

// Dependency is another Appfile that an App depends on
type Dependency struct {
	Source string
}

// Project is the structure of a project that many applications
// can belong to.
type Project struct {
	Name           string
	Infrastructure string
}

// Infrastructure is the structure of defining the infrastructure
// that an application must run on.
type Infrastructure struct {
	Name   string
	Type   string
	Flavor string
}

//-------------------------------------------------------------------
// Merging
//-------------------------------------------------------------------

// Merge will merge the other File onto this one, modifying this
// File with the merged contents.
func (f *File) Merge(other *File) error {
	if other.ID != "" {
		f.ID = other.ID
	}
	if other.Path != "" {
		f.Path = other.Path
	}

	// Application
	if f.Application == nil {
		f.Application = other.Application
	} else if other.Application != nil {
		// Note this won't copy dependencies properly
		*f.Application = *other.Application
	}

	// Project
	if f.Project == nil {
		f.Project = other.Project
	} else if other.Project != nil {
		// Note this won't copy dependencies properly
		*f.Project = *other.Project
	}

	// Infrastructure
	infraMap := make(map[string]int)
	for i, infra := range f.Infrastructure {
		infraMap[infra.Name] = i
	}
	for _, i := range other.Infrastructure {
		if idx, ok := infraMap[i.Name]; !ok {
			f.Infrastructure = append(f.Infrastructure, i)
		} else {
			f.Infrastructure[idx] = i
		}
	}

	// TODO: customizations
	f.Customization = other.Customization

	return nil
}

//-------------------------------------------------------------------
// Helper Methods
//-------------------------------------------------------------------

// ActiveInfrastructure returns the Infrastructure that is being
// used for this Appfile.
func (f *File) ActiveInfrastructure() *Infrastructure {
	for _, i := range f.Infrastructure {
		if i.Name == f.Project.Infrastructure {
			return i
		}
	}

	return nil
}

// resetID deletes the ID associated with this file.
func (f *File) resetID() error {
	return os.Remove(filepath.Join(filepath.Dir(f.Path), IDFile))
}

// hasID checks whether we have an ID file. This can return an error
// for filesystem errors.
func (f *File) hasID() (bool, error) {
	path := filepath.Join(filepath.Dir(f.Path), IDFile)
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	return err == nil, nil
}

// initID creates a new UUID and writes the file. This will overwrite
// any prior ID file.
func (f *File) initID() error {
	path := filepath.Join(filepath.Dir(f.Path), IDFile)
	uuid := uuid.GenerateUUID()
	data := strings.TrimSpace(fmt.Sprintf(idFileTemplate, uuid)) + "\n"
	return ioutil.WriteFile(path, []byte(data), 0644)
}

// loadID loads the ID for this File.
func (appF *File) loadID() error {
	hasID, err := appF.hasID()
	if err != nil {
		return err
	}
	if !hasID {
		appF.ID = ""
		return nil
	}

	path := filepath.Join(filepath.Dir(appF.Path), IDFile)
	uuid, err := oneline.Read(path)
	if err != nil {
		return err
	}

	appF.ID = uuid
	return nil
}

//-------------------------------------------------------------------
// GoStringer
//-------------------------------------------------------------------

func (v *Customization) GoString() string {
	return fmt.Sprintf("*%#v", *v)
}

func (v *Infrastructure) GoString() string {
	return fmt.Sprintf("*%#v", *v)
}

func (v *Project) GoString() string {
	return fmt.Sprintf("*%#v", *v)
}

const idFileTemplate = `
%s

DO NOT MODIFY OR DELETE THIS FILE!

This file should be checked in to version control. Do not ignore this file.

The first line is a unique UUID that represents the Appfile in this directory.
This UUID is used globally across your projects to identify this specific
Appfile. This UUID allows you to modify the name of an application, or have
duplicate application names without conflicting.

If you delete this file, then deploys may duplicate this application since
Otto will be unable to tell that the application is deployed.
`
