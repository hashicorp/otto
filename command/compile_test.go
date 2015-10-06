package command

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile/detect"
	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/otto"
	"github.com/mitchellh/cli"
)

func TestCompile(t *testing.T) {
	core := otto.TestCoreConfig(t)
	infra := otto.TestInfra(t, "test", core)
	appImpl := otto.TestApp(t, app.Tuple{"test", "test", "test"}, core)
	ui := new(cli.MockUi)
	c := &CompileCommand{
		Meta: Meta{
			CoreConfig: core,
			Ui:         ui,
		},
	}

	dir := fixtureDir("compile-basic")
	defer os.Remove(filepath.Join(dir, ".ottoid"))
	defer testChdir(t, dir)()

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	if !infra.CompileCalled {
		t.Fatal("Compile should be called")
	}
	if !appImpl.CompileCalled {
		t.Fatal("Compile should be called")
	}
}

func TestCompile_noExplicitType(t *testing.T) {
	core := otto.TestCoreConfig(t)
	infra := otto.TestInfra(t, "aws", core)
	appImpl := otto.TestApp(t, app.Tuple{"test-detection-merge", "aws", "simple"}, core)
	appImpl.CompileFunc = func(ctx *app.Context) (*app.CompileResult, error) {
		if ctx.Application == nil {
			t.Fatal("application unexpectedly nil")
		}
		if ctx.Application.Name != "compile-no-explicit-type" {
			t.Fatalf("expected: compile-no-explicit-type; got: %s", ctx.Application.Name)
		}
		if ctx.Application.Type != "test-detection-merge" {
			t.Fatalf("expected: test-detection-merge; got: %s", ctx.Application.Type)
		}
		return nil, nil
	}
	foundImpl := otto.TestFoundation(
		t, foundation.Tuple{"consul", "aws", "simple"}, core)
	ui := new(cli.MockUi)
	detectors := []*detect.Detector{
		&detect.Detector{
			Type: "test-detection-merge",
			File: []string{"test-file"},
		},
	}
	c := &CompileCommand{
		Meta: Meta{
			CoreConfig: core,
			Ui:         ui,
		},
		Detectors: detectors,
	}

	dir := fixtureDir("compile-detection")
	defer os.Remove(filepath.Join(dir, ".ottoid"))
	defer testChdir(t, dir)()

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	if !infra.CompileCalled {
		t.Fatal("Compile should be called")
	}
	if !appImpl.CompileCalled {
		t.Fatal("Compile should be called")
	}
	if !foundImpl.CompileCalled {
		t.Fatal("Foundation should be called")
	}
}

func TestCompile_noAppFile(t *testing.T) {
	core := otto.TestCoreConfig(t)
	infra := otto.TestInfra(t, "aws", core)
	appImpl := otto.TestApp(t, app.Tuple{"test", "aws", "simple"}, core)
	appImpl.CompileFunc = func(ctx *app.Context) (*app.CompileResult, error) {
		if ctx.Application == nil {
			t.Fatal("application unexpectedly nil")
		}
		if ctx.Application.Name != "compile-no-appfile" {
			t.Fatalf("expected: compile-no-appfile; got: %s", ctx.Application.Name)
		}
		return nil, nil
	}
	foundImpl := otto.TestFoundation(
		t, foundation.Tuple{"consul", "aws", "simple"}, core)
	ui := new(cli.MockUi)
	detectors := []*detect.Detector{
		&detect.Detector{
			Type: "test",
			File: []string{"test-file"},
		},
	}
	c := &CompileCommand{
		Meta: Meta{
			CoreConfig: core,
			Ui:         ui,
		},
		Detectors: detectors,
	}

	dir := fixtureDir("compile-no-appfile")
	defer os.Remove(filepath.Join(dir, ".ottoid"))
	defer testChdir(t, dir)()

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	if !infra.CompileCalled {
		t.Fatal("Compile should be called")
	}
	if !appImpl.CompileCalled {
		t.Fatal("Compile should be called")
	}
	if !foundImpl.CompileCalled {
		t.Fatal("Foundation should be called")
	}
}

func TestCompile_pathFile(t *testing.T) {
	ui := new(cli.MockUi)
	c := &CompileCommand{
		Meta: Meta{
			CoreConfig: otto.TestCoreConfig(t),
			Ui:         ui,
		},
	}

	dir := fixtureDir("compile-file")
	defer os.Remove(filepath.Join(dir, ".ottoid"))
	defer testChdir(t, dir)()

	args := []string{"-appfile", "Appfile.other"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestCompile_pathDir(t *testing.T) {
	ui := new(cli.MockUi)
	c := &CompileCommand{
		Meta: Meta{
			CoreConfig: otto.TestCoreConfig(t),
			Ui:         ui,
		},
	}

	dir := fixtureDir("compile-dir")
	defer os.Remove(filepath.Join(dir, "dir", ".ottoid"))
	defer testChdir(t, dir)()

	args := []string{"-appfile", "dir"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestCompile_altFile(t *testing.T) {
	ui := new(cli.MockUi)
	c := &CompileCommand{
		Meta: Meta{
			CoreConfig: otto.TestCoreConfig(t),
			Ui:         ui,
		},
	}

	dir := fixtureDir("compile-alt")
	defer os.Remove(filepath.Join(dir, ".ottoid"))
	defer testChdir(t, dir)()

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func testChdir(t *testing.T, dir string) func() {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	return func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatal(err)
		}
	}
}
