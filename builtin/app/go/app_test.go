package goapp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/otto"
)

func TestApp_impl(t *testing.T) {
	var _ app.App = new(App)
}

func TestApp_importPath_noGOPATH(t *testing.T) {
	compile.AppTest(true)
	defer compile.AppTest(false)

	// No GOPATH can be set
	defer os.Setenv("GOPATH", os.Getenv("GOPATH"))
	os.Setenv("GOPATH", "")

	otto.Test(t, otto.TestCase{
		Unit: true,
		Core: otto.TestCore(t, &otto.TestCoreOpts{
			Path: filepath.Join("./test-fixtures", "basic", "Appfile"),
			App:  new(App),
		}),

		Steps: []otto.TestStep{
			&compile.AppTestStepContext{
				Key:   "import_path",
				Value: "",
			},

			&compile.AppTestStepContext{
				Key:   "shared_folder_path",
				Value: "/vagrant",
			},
		},
	})
}

func TestApp_importPath(t *testing.T) {
	gopath := filepath.Join("./test-fixtures", "gopath")

	compile.AppTest(true)
	defer compile.AppTest(false)

	// No GOPATH can be set
	defer os.Setenv("GOPATH", os.Getenv("GOPATH"))
	os.Setenv("GOPATH", gopath)

	otto.Test(t, otto.TestCase{
		Unit: true,
		Core: otto.TestCore(t, &otto.TestCoreOpts{
			Path: filepath.Join(gopath, "src", "example.com", "Appfile"),
			App:  new(App),
		}),

		Steps: []otto.TestStep{
			&compile.AppTestStepContext{
				Key:   "import_path",
				Value: "example.com",
			},

			&compile.AppTestStepContext{
				Key:   "shared_folder_path",
				Value: "/opt/gopath/src/example.com",
			},
		},
	})
}
