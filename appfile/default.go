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
	appName := filepath.Base(dir)
	appType, err := detect.App(dir, det)
	if err != nil {
		return nil, err
	}

	return &File{
		Path: filepath.Join(dir, "Appfile"),

		Application: &Application{
			Name: appName,
			Type: appType,
		},

		Project: &Project{
			Name:           appName,
			Infrastructure: "aws",
		},

		Infrastructure: []*Infrastructure{
			&Infrastructure{
				Name:   appName,
				Type:   "aws",
				Flavor: "vpc-public-private",

				Foundations: []*Foundation{
					&Foundation{
						Name: "consul",
					},
				},
			},
		},
	}, nil
}
