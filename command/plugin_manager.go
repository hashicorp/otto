package command

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/builtin/pluginmap"
	"github.com/hashicorp/otto/otto"
	"github.com/hashicorp/otto/plugin"
	"github.com/kardianos/osext"
)

// PluginGlob is the glob pattern used to find plugins.
const PluginGlob = "otto-*"

// PluginManager is responsible for discovering and starting plugins.
//
// Plugin cleanup is done out in the main package: we just defer
// plugin.CleanupClients in main itself.
type PluginManager struct {
	// PluginDirs are the directories where plugins can be found.
	// Any plugins with the same types found later (higher index) will
	// override earlier (lower index) directories.
	PluginDirs []string

	plugins []*Plugin
}

// Plugin is a single plugin that has been loaded.
type Plugin struct {
	// Path and Args are the method used to invocate this plugin.
	// These are the only two values that need to be set manually. Once
	// these are set, call Load to load the plugin.
	Path string
	Args []string

	// The fields below are loaded as part of the Load() call and should
	// not be set manually, but can be accessed after Load.
	App     app.Factory `json:"-"`
	AppMeta *app.Meta   `json:"-"`

	used bool
}

// Load loads the plugin specified by the Path and instantiates the
// other fields on this structure.
func (p *Plugin) Load() error {
	// Create the plugin client to communicate with the process
	pluginClient := plugin.NewClient(&plugin.ClientConfig{
		Cmd:     exec.Command(p.Path, p.Args...),
		Managed: true,
	})

	// Request the client
	client, err := pluginClient.Client()
	if err != nil {
		return err
	}

	// Get the app implementation
	appImpl, err := client.App()
	if err != nil {
		return err
	}

	p.AppMeta, err = appImpl.Meta()
	if err != nil {
		return err
	}

	// Create a custom factory that when called marks the plugin as used
	p.used = false
	p.App = func() (app.App, error) {
		p.used = true
		return client.App()
	}

	return nil
}

// Used tracks whether or not this plugin was used or not. You can call
// this after compilation on each plugin to determine what plugin
// was used.
func (p *Plugin) Used() bool {
	return p.used
}

func (p *Plugin) String() string {
	return fmt.Sprintf("%s %v", p.Path, p.Args)
}

// ConfigureCore configures the Otto core configuration with the loaded
// plugin data.
func (m *PluginManager) ConfigureCore(core *otto.CoreConfig) error {
	if core.Apps == nil {
		core.Apps = make(map[app.Tuple]app.Factory)
	}

	for _, p := range m.Plugins() {
		for _, tuple := range p.AppMeta.Tuples {
			core.Apps[tuple] = p.App
		}
	}

	return nil
}

// Plugins returns the loaded plugins.
func (m *PluginManager) Plugins() []*Plugin {
	return m.plugins
}

// Discover will find all the available plugin binaries. Each time this
// is called it will override any previously discovered plugins.
func (m *PluginManager) Discover() error {
	result := make([]*Plugin, 0, 20)

	if !testingMode {
		// Get our own path
		exePath, err := osext.Executable()
		if err != nil {
			return err
		}

		// First we add all the builtin plugins which we get by executing ourself
		for k, _ := range pluginmap.Apps {
			result = append(result, &Plugin{
				Path: exePath,
				Args: []string{"plugin-builtin", "app", k},
			})
		}
	}

	// Log it
	for _, r := range result {
		log.Printf("[DEBUG] Detected plugin: %s", r)
	}

	// Save our result
	m.plugins = result

	return nil
}

// StoreUsed will persist the used plugins into a file. LoadUsed can
// then be called to load the plugins that were used only, making plugin
// loading much more efficient.
func (m *PluginManager) StoreUsed(path string) error {
	// Get the used plugins
	plugins := make([]*Plugin, 0, 2)
	for _, p := range m.Plugins() {
		if p.Used() {
			plugins = append(plugins, p)
		}
	}

	// Write the used plugins to the given path as JSON
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(&usedPluginWrapper{
		Version: usedPluginVersion,
		Plugins: plugins,
	})
}

// LoadUsed will load the plugins in the given used file that was saved
// with StoreUsed.
func (m *PluginManager) LoadUsed(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	var wrapper usedPluginWrapper
	dec := json.NewDecoder(f)
	err = dec.Decode(&wrapper)
	f.Close()
	if err != nil {
		return err
	}

	if wrapper.Version > usedPluginVersion {
		// TODO: error
	}

	m.plugins = wrapper.Plugins
	return m.LoadAll()
}

// LoadAll will launch every plugin and add it to the CoreConfig given.
func (m *PluginManager) LoadAll() error {
	// If we've never loaded plugin paths, then let's discover those first
	if m.Plugins() == nil {
		if err := m.Discover(); err != nil {
			return err
		}
	}

	// Go through each plugin path and load single
	// TODO: parallelize
	for _, plugin := range m.Plugins() {
		if err := plugin.Load(); err != nil {
			return err
		}
	}

	return nil
}

// usedPluginVersion is the current version of the used plugin format
// that we understand. We can increment and handle older versions as we go.
const usedPluginVersion int = 1

type usedPluginWrapper struct {
	Version int       `json:"version"`
	Plugins []*Plugin `json:"plugins"`
}
