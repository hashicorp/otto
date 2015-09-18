package otto

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/otto/directory"
)

// statusInfo holds the complete status information for the Core.Status
// function.
type statusInfo struct {
	Err    error
	Dev    *directory.Dev
	Build  *directory.Build
	Deploy *directory.Deploy
	Infra  *directory.Infra
}

// statusInfo gets the information for the Status call.
//
// This is meant to be called in a goroutine.
func (c *Core) statusInfo(resultCh chan<- *statusInfo) {
	infra := c.appfile.ActiveInfrastructure()
	if infra == nil {
		panic("infra not found")
	}

	var err error
	var result statusInfo

	// Dev
	result.Dev, err = c.dir.GetDev(&directory.Dev{Lookup: directory.Lookup{
		AppID: c.appfile.ID}})
	if err != nil {
		result.Err = multierror.Append(result.Err, fmt.Errorf(
			"Error loading development status: %s", err))
	}

	// Build
	result.Build, err = c.dir.GetBuild(&directory.Build{Lookup: directory.Lookup{
		AppID: c.appfile.ID, Infra: infra.Name, InfraFlavor: infra.Flavor}})
	if err != nil {
		result.Err = multierror.Append(result.Err, fmt.Errorf(
			"Error loading build status: %s", err))
	}

	// Deploy
	result.Deploy, err = c.dir.GetDeploy(&directory.Deploy{Lookup: directory.Lookup{
		AppID: c.appfile.ID, Infra: infra.Name, InfraFlavor: infra.Flavor}})
	if err != nil {
		result.Err = multierror.Append(result.Err, fmt.Errorf(
			"Error loading deploy status: %s", err))
	}

	// Infra
	result.Infra, err = c.dir.GetInfra(&directory.Infra{Lookup: directory.Lookup{
		Infra: infra.Name}})
	if err != nil {
		result.Err = multierror.Append(result.Err, fmt.Errorf(
			"Error loading infra status: %s", err))
	}

	resultCh <- &result
}
