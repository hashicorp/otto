package main

import (
	"fmt"
	"encoding/json"
	"log"
	"io/ioutil"
	"os"
	"path/filepath"

	sp "{{ working_gopath }}"
)

const(
	scriptpackDir = "_scriptpack_staging"
)

func main() {
	// Logs don't matter for this since it should always work. If it
	// doesn't, it will show up in the error output.
	log.SetOutput(ioutil.Discard)

	path := filepath.Join("{{ path.working }}", scriptpackDir)
	if err := os.RemoveAll(path); err !=nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}

	// Write the actual ScriptPack data
	if err := sp.ScriptPack.Write(path); err !=nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}

	// Write the ScriptPack env vars
	envPath := filepath.Join(path, "env")
	output, err := json.MarshalIndent(
		sp.ScriptPack.Env("/devroot/"+scriptpackDir), "", "\t")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
	if err := ioutil.WriteFile(envPath, output, 0644); err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}
