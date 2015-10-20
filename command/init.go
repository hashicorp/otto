package command

import (
	"fmt"
	"io/ioutil"
)

type InitCommand struct {
	Meta
}

type FileConfigurations struct {
	Language string
}

var appFile = `application {
	name = "otto"
	type = "%v"
}\n`

func (c *InitCommand) Run(args []string) int {
	file := FileConfigurations{"go"}
	generateAppFile(file)
	return 0
}

func (c *InitCommand) Synopsis() string {
	return "Initializes the Appfile"
}

func (c *InitCommand) Help() string {
	return "I need somebody..."
}

func generateAppFile(config FileConfigurations) {
	appFileString := fmt.Sprintf(appFile, "jey")
	ioutil.WriteFile("Appfile", []byte(appFileString), 0644)
}
