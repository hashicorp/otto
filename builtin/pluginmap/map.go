// pluginmap is a package that contains the mapping of internal plugin
// names to their factories that we can register it in various ways.
package pluginmap

import (
	"github.com/hashicorp/otto/plugin"

	appCustom "github.com/hashicorp/otto/builtin/app/custom"
	appDockerExt "github.com/hashicorp/otto/builtin/app/docker-external"
	appGo "github.com/hashicorp/otto/builtin/app/go"
	appNode "github.com/hashicorp/otto/builtin/app/node"
	appPHP "github.com/hashicorp/otto/builtin/app/php"
	appPython "github.com/hashicorp/otto/builtin/app/python"
	appRuby "github.com/hashicorp/otto/builtin/app/ruby"
	appScriptPack "github.com/hashicorp/otto/builtin/app/scriptpack"
)

var Map = map[string]*plugin.ServeOpts{
	"app-custom":          &plugin.ServeOpts{AppFunc: appCustom.AppFactory},
	"app-docker-external": &plugin.ServeOpts{AppFunc: appDockerExt.AppFactory},
	"app-go":              &plugin.ServeOpts{AppFunc: appGo.AppFactory},
	"app-node":            &plugin.ServeOpts{AppFunc: appNode.AppFactory},
	"app-php":             &plugin.ServeOpts{AppFunc: appPHP.AppFactory},
	"app-python":          &plugin.ServeOpts{AppFunc: appPython.AppFactory},
	"app-ruby":            &plugin.ServeOpts{AppFunc: appRuby.AppFactory},
	"app-scriptpack":      &plugin.ServeOpts{AppFunc: appScriptPack.AppFactory},
}
