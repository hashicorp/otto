package main

import (
	"os"
	"os/signal"

	appGo "github.com/hashicorp/otto/builtin/app/go"
	infraAws "github.com/hashicorp/otto/builtin/infra/aws"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/command"
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
	Ui = &cli.PrefixedUi{
		AskPrefix:    OutputPrefix,
		OutputPrefix: OutputPrefix,
		InfoPrefix:   OutputPrefix,
		ErrorPrefix:  ErrorPrefix,
		Ui:           &cli.BasicUi{Writer: os.Stdout},
	}

	meta := command.Meta{
		CoreConfig: &otto.CoreConfig{
			Apps: appGo.Tuples.Map(app.StructFactory(new(appGo.App))),
			Infrastructures: map[string]infrastructure.Factory{
				"aws": infrastructure.StructFactory(new(infraAws.Infra)),
			},
		},
		Ui: Ui,
	}

	Commands = map[string]cli.CommandFactory{
		"compile": func() (cli.Command, error) {
			return &command.CompileCommand{
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
