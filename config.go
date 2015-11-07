package main

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

// Config is the structure of the configuration for the Otto CLI.
//
// This is not the configuration for Otto itself. That is in the
// "config" package.
type Config struct {
	DisableCheckpoint          bool `hcl:"disable_checkpoint"`
	DisableCheckpointSignature bool `hcl:"disable_checkpoint_signature"`
}

// BuiltinConfig is the built-in defaults for the configuration. These
// can be overridden by user configurations.
var BuiltinConfig Config

// ConfigDir returns the configuration directory for Otto.
func ConfigDir() (string, error) {
	return configDir()
}
