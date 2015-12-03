package appfile

import (
	"path/filepath"

	"github.com/hashicorp/otto/appfile/detect"
)

// Default will generate a default Appfile for the given directory.
//
// The path to the directory must be absolute, since the path is used
// as a way to determine the name of the application.
func Default(dir string, det *detect.Config) (*File, error) {
	var appType string
	appName := filepath.Base(dir)
	if det != nil {
		t, err := detect.App(dir, det)
		if err != nil {
			return nil, err
		}

		appType = t
	}

	return &File{
		Path: filepath.Join(dir, "Appfile"),

		Application: &Application{
			Name:   appName,
			Type:   appType,
			Detect: true,
		},

		Project: &Project{
			Name:           appName,
			Infrastructure: appName,
		},

		Infrastructure: []*Infrastructure{
			&Infrastructure{
				Name:   appName,
				Type:   "aws",
				Flavor: "simple",

				Foundations: []*Foundation{
					&Foundation{
						Name: "consul",
					},
				},
			},
		},
	}, nil
}
