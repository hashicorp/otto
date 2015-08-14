package main

import (
	"log"
	"os"
	"os/signal"
)

func init() {
	// Setup a listener for interrupts (SIGINT) and just black hole
	// the signals. We do this because subcommands should setup additional
	// listeners if they care. And if we don't do this then the subprocesses
	// we execute such as Terraform and Vagrant never receive the ctrl-C
	// and can't gracefully clean up.
	signalCh := make(chan os.Signal, 2)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		for {
			<-signalCh
			log.Printf("[DEBUG] main: interrupt received. ignoring since command should also listen")
		}
	}()
}
