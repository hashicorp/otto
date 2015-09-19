package foundation

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/otto/context"
	"github.com/hashicorp/otto/directory"
)

// WriteVars will write the outputs from a foundation and put them into
// a key/value file within the foundation directories that is available
// when uploaded. This allows the build scripts to access the outputs
// at runtime.
//
// Ideally this sort of thing would be possible at compile-time but the
// values of these variables come at runtime so we need to populate this
// at runtime.
//
// By having a var file, it removes the burden of knowing what variables
// to send in to the foundation from the App implementation. They just need
// to call this function, upload the folder, and call main.sh.
func WriteVars(ctx *context.Shared) error {
	infra := ctx.Appfile.ActiveInfrastructure()
	if infra == nil {
		panic("no active infra")
	}

	if len(ctx.FoundationDirs) < len(infra.Foundations) {
		panic("foundationDirs is missing entries")
	}

	// subdirs are the directories to write the vars file to. For now
	// we just write it to one but we put this here because I can see
	// a future where more may need it.
	subdirs := []string{"app-build"}

	// Go through each foundation, grab its infrastructure data, and
	// write it out to the proper path.
	for i, f := range infra.Foundations {
		entry, err := ctx.Directory.GetInfra(&directory.Infra{
			Lookup: directory.Lookup{Infra: infra.Type, Foundation: f.Name}})
		if err != nil {
			return err
		}
		if entry == nil {
			continue
		}

		// Get the var file contents into an in-memory buffer
		var buf bytes.Buffer
		if err := varFile(&buf, entry); err != nil {
			return fmt.Errorf(
				"error generating var file for %s: %s",
				f.Name, err)
		}

		// Get a reader around the buffer so we can seek
		r := bytes.NewReader(buf.Bytes())
		for _, subdir := range subdirs {
			if _, err := r.Seek(0, 0); err != nil {
				return err
			}

			path := filepath.Join(ctx.FoundationDirs[i], subdir)
			if _, err := os.Stat(path); err != nil {
				if os.IsNotExist(err) {
					// Ignore directories that don't exist
					continue
				}

				return err
			}
			path = filepath.Join(path, "vars")

			log.Printf(
				"[DEBUG] writing foundation %s var file to: %s",
				f.Name, path)
			w, err := os.Create(path)
			if err != nil {
				return err
			}
			_, err = io.Copy(w, r)
			w.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func varFile(w io.Writer, entry *directory.Infra) error {
	bufW := bufio.NewWriter(w)
	for k, v := range entry.Outputs {
		_, err := bufW.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		if err != nil {
			return err
		}
	}

	return bufW.Flush()
}
