package terraform

import (
	"fmt"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/directory"
)

// DeployArtifactExtractor is the function type that is used to
// extract artifacts from a Build for deploys.
type DeployArtifactExtractor func(
	*app.Context, *directory.Build, *directory.Infra) (map[string]string, error)

var deployArtifactExtractors = map[string]DeployArtifactExtractor{
	"aws": deployArtifactExtractAWS,
	"digitalocean": deployArtifactExtractDo,
}

func deployArtifactExtractAWS(
	ctx *app.Context,
	build *directory.Build,
	infra *directory.Infra) (map[string]string, error) {
	ami, ok := build.Artifact[infra.Outputs["region"]]
	if !ok {
		return nil, fmt.Errorf(
			"An artifact for the region '%s' could not be found. Please run\n"+
				"`otto build` and try again.",
			infra.Outputs["region"])
	}

	return map[string]string{"ami": ami}, nil
}

func deployArtifactExtractDo(
	ctx *app.Context,
	build *directory.Build,
	infra *directory.Infra) (map[string]string, error) {
	droplet, ok := build.Artifact[infra.Outputs["region"]]
	if !ok {
		return nil, fmt.Errorf(
			"An artifact for the region '%s' could not be found. Please run\n"+
				"`otto build` and try again.",
			infra.Outputs["region"])
	}

	return map[string]string{"image": droplet}, nil
}
