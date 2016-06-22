package pongo2

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// A template set allows you to create your own group of templates with their own global context (which is shared
// among all members of the set), their own configuration (like a specific base directory) and their own sandbox.
// It's useful for a separation of different kind of templates (e. g. web templates vs. mail templates).
type TemplateSet struct {
	name string

	// Globals will be provided to all templates created within this template set
	Globals Context

	// If debug is true (default false), ExecutionContext.Logf() will work and output to STDOUT. Furthermore,
	// FromCache() won't cache the templates. Make sure to synchronize the access to it in case you're changing this
	// variable during program execution (and template compilation/execution).
	Debug bool

	// Base directory: If you set the base directory (string is non-empty), all filename lookups in tags/filters are
	// relative to this directory. If it's empty, all lookups are relative to the current filename which is importing.
	baseDirectory string

	// Sandbox features
	// - Limit access to directories (using SandboxDirectories)
	// - Disallow access to specific tags and/or filters (using BanTag() and BanFilter())
	//
	// You can limit file accesses (for all tags/filters which are using pongo2's file resolver technique)
	// to these sandbox directories. All default pongo2 filters/tags are respecting these restrictions.
	// For example, if you only have your base directory in the list, a {% ssi "/etc/passwd" %} will not work.
	// No items in SandboxDirectories means no restrictions at all.
	//
	// For efficiency reasons you can ban tags/filters only *before* you have added your first
	// template to the set (restrictions are statically checked). After you added one, it's not possible anymore
	// (for your personal security).
	//
	// SandboxDirectories can be changed at runtime. Please synchronize the access to it if you need to change it
	// after you've added your first template to the set. You *must* use this match pattern for your directories:
	// http://golang.org/pkg/path/filepath/#Match
	SandboxDirectories   []string
	firstTemplateCreated bool
	bannedTags           map[string]bool
	bannedFilters        map[string]bool

	// Template cache (for FromCache())
	templateCache      map[string]*Template
	templateCacheMutex sync.Mutex
}

// Create your own template sets to separate different kind of templates (e. g. web from mail templates) with
// different globals or other configurations (like base directories).
func NewSet(name string) *TemplateSet {
	return &TemplateSet{
		name:          name,
		Globals:       make(Context),
		bannedTags:    make(map[string]bool),
		bannedFilters: make(map[string]bool),
		templateCache: make(map[string]*Template),
	}
}

// Use this function to set your template set's base directory. This directory will be used for any relative
// path in filters, tags and From*-functions to determine your template.
func (set *TemplateSet) SetBaseDirectory(name string) error {
	// Make the path absolute
	if !filepath.IsAbs(name) {
		abs, err := filepath.Abs(name)
		if err != nil {
			return err
		}
		name = abs
	}

	// Check for existence
	fi, err := os.Stat(name)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("The given path '%s' is not a directory.")
	}

	set.baseDirectory = name
	return nil
}

func (set *TemplateSet) BaseDirectory() string {
	return set.baseDirectory
}

// Ban a specific tag for this template set. See more in the documentation for TemplateSet.
func (set *TemplateSet) BanTag(name string) {
	_, has := tags[name]
	if !has {
		panic(fmt.Sprintf("Tag '%s' not found.", name))
	}
	if set.firstTemplateCreated {
		panic("You cannot ban any tags after you've added your first template to your template set.")
	}
	_, has = set.bannedTags[name]
	if has {
		panic(fmt.Sprintf("Tag '%s' is already banned.", name))
	}
	set.bannedTags[name] = true
}

// Ban a specific filter for this template set. See more in the documentation for TemplateSet.
func (set *TemplateSet) BanFilter(name string) {
	_, has := filters[name]
	if !has {
		panic(fmt.Sprintf("Filter '%s' not found.", name))
	}
	if set.firstTemplateCreated {
		panic("You cannot ban any filters after you've added your first template to your template set.")
	}
	_, has = set.bannedFilters[name]
	if has {
		panic(fmt.Sprintf("Filter '%s' is already banned.", name))
	}
	set.bannedFilters[name] = true
}

// FromCache() is a convenient method to cache templates. It is thread-safe
// and will only compile the template associated with a filename once.
// If TemplateSet.Debug is true (for example during development phase),
// FromCache() will not cache the template and instead recompile it on any
// call (to make changes to a template live instantaneously).
// Like FromFile(), FromCache() takes a relative path to a set base directory.
// Sandbox restrictions apply (if given).
func (set *TemplateSet) FromCache(filename string) (*Template, error) {
	if set.Debug {
		// Recompile on any request
		return set.FromFile(filename)
	} else {
		// Cache the template
		cleaned_filename := set.resolveFilename(nil, filename)

		set.templateCacheMutex.Lock()
		defer set.templateCacheMutex.Unlock()

		tpl, has := set.templateCache[cleaned_filename]

		// Cache miss
		if !has {
			tpl, err := set.FromFile(cleaned_filename)
			if err != nil {
				return nil, err
			}
			set.templateCache[cleaned_filename] = tpl
			return tpl, nil
		}

		// Cache hit
		return tpl, nil
	}
}

// Loads  a template from string and returns a Template instance.
func (set *TemplateSet) FromString(tpl string) (*Template, error) {
	set.firstTemplateCreated = true

	return newTemplateString(set, tpl)
}

// Loads a template from a filename and returns a Template instance.
// If a base directory is set, the filename must be either relative to it
// or be an absolute path. Sandbox restrictions (SandboxDirectories) apply
// if given.
func (set *TemplateSet) FromFile(filename string) (*Template, error) {
	set.firstTemplateCreated = true

	buf, err := ioutil.ReadFile(set.resolveFilename(nil, filename))
	if err != nil {
		return nil, &Error{
			Filename: filename,
			Sender:   "fromfile",
			ErrorMsg: err.Error(),
		}
	}
	return newTemplate(set, filename, false, string(buf))
}

// Shortcut; renders a template string directly. Panics when providing a
// malformed template or an error occurs during execution.
func (set *TemplateSet) RenderTemplateString(s string, ctx Context) string {
	set.firstTemplateCreated = true

	tpl := Must(set.FromString(s))
	result, err := tpl.Execute(ctx)
	if err != nil {
		panic(err)
	}
	return result
}

// Shortcut; renders a template file directly. Panics when providing a
// malformed template or an error occurs during execution.
func (set *TemplateSet) RenderTemplateFile(fn string, ctx Context) string {
	set.firstTemplateCreated = true

	tpl := Must(set.FromFile(fn))
	result, err := tpl.Execute(ctx)
	if err != nil {
		panic(err)
	}
	return result
}

func (set *TemplateSet) logf(format string, args ...interface{}) {
	if set.Debug {
		logger.Printf(fmt.Sprintf("[template set: %s] %s", set.name, format), args...)
	}
}

// Resolves a filename relative to the base directory. Absolute paths are allowed.
// If sandbox restrictions are given (SandboxDirectories), they will be respected and checked.
// On sandbox restriction violation, resolveFilename() panics.
func (set *TemplateSet) resolveFilename(tpl *Template, filename string) (resolved_path string) {
	if len(set.SandboxDirectories) > 0 {
		defer func() {
			// Remove any ".." or other crap
			resolved_path = filepath.Clean(resolved_path)

			// Make the path absolute
			abs_path, err := filepath.Abs(resolved_path)
			if err != nil {
				panic(err)
			}
			resolved_path = abs_path

			// Check against the sandbox directories (once one pattern matches, we're done and can allow it)
			for _, pattern := range set.SandboxDirectories {
				matched, err := filepath.Match(pattern, resolved_path)
				if err != nil {
					panic("Wrong sandbox directory match pattern (see http://golang.org/pkg/path/filepath/#Match).")
				}
				if matched {
					// OK!
					return
				}
			}

			// No pattern matched, we have to log+deny the request
			set.logf("Access attempt outside of the sandbox directories (blocked): '%s'", resolved_path)
			resolved_path = ""
		}()
	}

	if filepath.IsAbs(filename) {
		return filename
	}

	if set.baseDirectory == "" {
		if tpl != nil {
			if tpl.is_tpl_string {
				return filename
			}
			base := filepath.Dir(tpl.name)
			return filepath.Join(base, filename)
		}
		return filename
	} else {
		return filepath.Join(set.baseDirectory, filename)
	}
}

// Logging function (internally used)
func logf(format string, items ...interface{}) {
	if debug {
		logger.Printf(format, items...)
	}
}

var (
	debug  bool // internal debugging
	logger = log.New(os.Stdout, "[pongo2] ", log.LstdFlags)

	// Creating a default set
	DefaultSet = NewSet("default")

	// Methods on the default set
	FromString           = DefaultSet.FromString
	FromFile             = DefaultSet.FromFile
	FromCache            = DefaultSet.FromCache
	RenderTemplateString = DefaultSet.RenderTemplateString
	RenderTemplateFile   = DefaultSet.RenderTemplateFile

	// Globals for the default set
	Globals = DefaultSet.Globals
)
