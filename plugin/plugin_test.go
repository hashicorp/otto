package plugin

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/hashicorp/otto/app"
	pluginrpc "github.com/hashicorp/otto/rpc"
)

func helperProcess(s ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, s...)
	env := []string{
		"GO_WANT_HELPER_PROCESS=1",
		"OTTO_PLUGIN_MIN_PORT=10000",
		"OTTO_PLUGIN_MAX_PORT=25000",
	}

	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(env, os.Environ()...)
	return cmd
}

// This is not a real test. This is just a helper process kicked off by
// tests.
func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}

		args = args[1:]
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "bad-version":
		fmt.Printf("%s1|tcp|:1234\n", APIVersion)
		<-make(chan int)
	case "app":
		Serve(&ServeOpts{
			AppFunc: testAppFixed(new(app.Mock)),
		})
	case "invalid-rpc-address":
		fmt.Println("lolinvalid")
	case "mock":
		fmt.Printf("%s|tcp|:1234\n", APIVersion)
		<-make(chan int)
	case "start-timeout":
		time.Sleep(1 * time.Minute)
		os.Exit(1)
	case "stderr":
		fmt.Printf("%s|tcp|:1234\n", APIVersion)
		log.Println("HELLO")
		log.Println("WORLD")
	case "stdin":
		fmt.Printf("%s|tcp|:1234\n", APIVersion)
		data := make([]byte, 5)
		if _, err := os.Stdin.Read(data); err != nil {
			log.Printf("stdin read error: %s", err)
			os.Exit(100)
		}

		if string(data) == "hello" {
			os.Exit(0)
		}

		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %q\n", cmd)
		os.Exit(2)
	}
}

func testAppFixed(p app.App) pluginrpc.AppFunc {
	return func() app.App {
		return p
	}
}
