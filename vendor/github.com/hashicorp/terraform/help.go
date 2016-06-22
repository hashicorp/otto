package main

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/mitchellh/cli"
)

// helpFunc is a cli.HelpFunc that can is used to output the help for Terraform.
func helpFunc(commands map[string]cli.CommandFactory) string {
	// Determine the maximum key length, and classify based on type
	porcelain := make(map[string]cli.CommandFactory)
	plumbing := make(map[string]cli.CommandFactory)
	maxKeyLen := 0
	for key, f := range commands {
		if len(key) > maxKeyLen {
			maxKeyLen = len(key)
		}

		if _, ok := PlumbingCommands[key]; ok {
			plumbing[key] = f
		} else {
			porcelain[key] = f
		}
	}

	var buf bytes.Buffer
	buf.WriteString("usage: terraform [--version] [--help] <command> [args]\n\n")
	buf.WriteString(
		"The available commands for execution are listed below.\n" +
			"The most common, useful commands are shown first, followed by\n" +
			"less common or more advanced commands. If you're just getting\n" +
			"started with Terraform, stick with the common commands. For the\n" +
			"other commands, please read the help and docs before usage.\n\n")
	buf.WriteString("Common commands:\n")
	buf.WriteString(listCommands(porcelain, maxKeyLen))
	buf.WriteString("\nAll other commands:\n")
	buf.WriteString(listCommands(plumbing, maxKeyLen))
	return buf.String()
}

// listCommands just lists the commands in the map with the
// given maximum key length.
func listCommands(commands map[string]cli.CommandFactory, maxKeyLen int) string {
	var buf bytes.Buffer

	// Get the list of keys so we can sort them, and also get the maximum
	// key length so they can be aligned properly.
	keys := make([]string, 0, len(commands))
	for key, _ := range commands {
		// This is an internal command that users should never call directly so
		// we will hide it from the command listing.
		if key == "internal-plugin" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		commandFunc, ok := commands[key]
		if !ok {
			// This should never happen since we JUST built the list of
			// keys.
			panic("command not found: " + key)
		}

		command, err := commandFunc()
		if err != nil {
			log.Printf("[ERR] cli: Command '%s' failed to load: %s",
				key, err)
			continue
		}

		key = fmt.Sprintf("%s%s", key, strings.Repeat(" ", maxKeyLen-len(key)))
		buf.WriteString(fmt.Sprintf("    %s    %s\n", key, command.Synopsis()))
	}

	return buf.String()
}
