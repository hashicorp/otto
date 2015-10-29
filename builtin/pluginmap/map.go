// pluginmap is a package that contains the mapping of internal plugin
// names to their factories that we can register it in various ways.
package pluginmap

import (
	"github.com/hashicorp/otto/rpc"

	appCustom "github.com/hashicorp/otto/builtin/app/custom"
	appDockerExt "github.com/hashicorp/otto/builtin/app/docker-external"
	appGo "github.com/hashicorp/otto/builtin/app/go"
	appNode "github.com/hashicorp/otto/builtin/app/node"
	appPHP "github.com/hashicorp/otto/builtin/app/php"
	appRuby "github.com/hashicorp/otto/builtin/app/ruby"
)

var Apps = map[string]rpc.AppFunc{
	"custom":          appCustom.AppFactory,
	"docker-external": appDockerExt.AppFactory,
	"go":              appGo.AppFactory,
	"node":            appNode.AppFactory,
	"php":             appPHP.AppFactory,
	"ruby":            appRuby.AppFactory,
}
