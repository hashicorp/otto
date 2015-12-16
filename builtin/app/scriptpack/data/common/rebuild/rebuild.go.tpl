package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"os"
	"path/filepath"

	sp "{{ working_gopath }}"
)

func main() {
	// Logs don't matter for this since it should always work. If it
	// doesn't, it will show up in the error output.
	log.SetOutput(ioutil.Discard)

	path := filepath.Join("{{ path.working }}", "_scriptpack_staging")
	if err := os.RemoveAll(path); err !=nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}

	if err := sp.ScriptPack.Write(path); err !=nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}
