package main

import (
	"os"
	"os/signal"

	appCustom "github.com/hashicorp/otto/builtin/app/custom"
	appGo "github.com/hashicorp/otto/builtin/app/go"
	appRuby "github.com/hashicorp/otto/builtin/app/ruby"
	foundationConsul "github.com/hashicorp/otto/builtin/foundation/consul"
	infraAws "github.com/hashicorp/otto/builtin/infra/aws"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile/detect"
	"github.com/hashicorp/otto/command"
	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/otto"
	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available Otto commands.
var Commands map[string]cli.CommandFactory

// Ui is the cli.Ui used for communicating to the outside world.
var Ui cli.Ui

const (
	ErrorPrefix  = "e:"
	OutputPrefix = "o:"
)

func init() {
	Ui = &cli.ColoredUi{
		OutputColor: cli.UiColorNone,
		InfoColor:   cli.UiColorNone,
		ErrorColor:  cli.UiColorRed,
		WarnColor:   cli.UiColorNone,
		Ui: &cli.PrefixedUi{
			AskPrefix:    OutputPrefix,
			OutputPrefix: OutputPrefix,
			InfoPrefix:   OutputPrefix,
			ErrorPrefix:  ErrorPrefix,
			Ui:           &cli.BasicUi{Writer: os.Stdout},
		},
	}

	apps := appGo.Tuples.Map(app.StructFactory(new(appGo.App)))
	apps.Add(appCustom.Tuples.Map(app.StructFactory(new(appCustom.App))))
	apps.Add(appRuby.Tuples.Map(app.StructFactory(new(appRuby.App))))

	foundations := foundationConsul.Tuples.Map(foundation.StructFactory(new(foundationConsul.Foundation)))

	meta := command.Meta{
		CoreConfig: &otto.CoreConfig{
			Apps:        apps,
			Foundations: foundations,
			Infrastructures: map[string]infrastructure.Factory{
				"aws": infraAws.Infra,
			},
		},
		Ui: Ui,
	}

	detectors := []*detect.Detector{
		&detect.Detector{
			Type: "go",
			File: []string{"*.go"},
		},
		&detect.Detector{
			Type: "ruby",
			File: []string{"*.rb", "Gemfile", "config.ru"},
		},
	}

	Commands = map[string]cli.CommandFactory{
		"compile": func() (cli.Command, error) {
			return &command.CompileCommand{
				Meta:      meta,
				Detectors: detectors,
			}, nil
		},

		"build": func() (cli.Command, error) {
			return &command.BuildCommand{
				Meta: meta,
			}, nil
		},

		"deploy": func() (cli.Command, error) {
			return &command.DeployCommand{
				Meta: meta,
			}, nil
		},

		"dev": func() (cli.Command, error) {
			return &command.DevCommand{
				Meta: meta,
			}, nil
		},

		"infra": func() (cli.Command, error) {
			return &command.InfraCommand{
				Meta: meta,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Meta:              meta,
				Revision:          GitCommit,
				Version:           Version,
				VersionPrerelease: VersionPrerelease,
				CheckFunc:         commandVersionCheck,
			}, nil
		},
	}
}

// makeShutdownCh creates an interrupt listener and returns a channel.
// A message will be sent on the channel for every interrupt received.
func makeShutdownCh() <-chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()

	return resultCh
}
