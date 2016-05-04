package command

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type InitCommand struct {
	Meta
}

type AppFileConfig struct {
	Language string
	Version  string
}

var appFile = `#Application Configurations
application {
  name = "otto"
  type = "%v"
 
  #You can Add dependencies to your application enviroment
  #For more information, check this link: https://ottoproject.io/docs/appfile/dep-sources.html 
  #dependency { source = "github.com/hashicorp/otto/examples/mongodb" }
}

#You can bind your application to a project and an infrastructure
#For more information about it, check this link: https://ottoproject.io/docs/appfile/project.html
#project {
#  name = "otto"
#  infrastructure = "production"
#}

#You can describe what infrastructure your application will be deployed to
#For more information about it, check this link: https://ottoproject.io/docs/appfile/infra.html
#infrastructure "production" {
#  type = "aws"
#  flavor = "vpc-public-private"
#}

#Customization blocks change the behavior of Otto for a specific application type or infrastructure. 
#For more information, check this link: https://ottoproject.io/docs/appfile/customization.html
#customization "%v" {
#  %v_version = "%v"
#}

#The import statement can be used within an Appfile to import fragments of an Appfile from other sources, including the #local filesystem and remote URLs.
#For more information about it, check this link: https://ottoproject.io/docs/appfile/customization.html
#import "github.com/hashicorp/otto-shared/database" {}

`

func (c *InitCommand) Run(args []string) int {
	var language string
	var version string

	fs := c.FlagSet("init", FlagSetNone)
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	fs.StringVar(&language, "lang", "go", "")
	fs.StringVar(&version, "langversion", "1.0", "")
	if err := fs.Parse(args); err != nil {
		return 1
	}

	file := AppFileConfig{language, version}
	generateAppFile(file)
	return 0
}

func (c *InitCommand) Synopsis() string {
	return "Initializes the Appfile"
}

func (c *InitCommand) Help() string {
	helpText := `
Usage: otto init [options]

  Generates the AppFile 

  ==OPTIONS==

  -lang		The applications's language. i.e: Ruby, PHP
  -langversion	The version of the applications's language. i.e: 1.4, 4.0
`
	return strings.TrimSpace(helpText)
}

func generateAppFile(config AppFileConfig) {
	appFileString := fmt.Sprintf(appFile, config.Language, config.Language, config.Language, config.Version)
	ioutil.WriteFile("Appfile", []byte(appFileString), 0644)
}
