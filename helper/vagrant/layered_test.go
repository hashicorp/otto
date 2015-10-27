package vagrant

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hashicorp/otto/context"
	"github.com/hashicorp/otto/helper/exec"
	"github.com/hashicorp/otto/ui"
	"github.com/hashicorp/terraform/dag"
)

func TestLayerVertex_impl(t *testing.T) {
	var _ dag.Hashable = new(layerVertex)
}

func TestLayeredBuild(t *testing.T) {
	dir := tempDir(t)
	defer os.RemoveAll(dir)

	runner := new(exec.MockRunner)
	defer exec.TestChrunner(runner.Run)()

	// Build an environment using foo and bar
	layer := &Layered{
		DataDir: dir,
		Layers: []*Layer{
			testLayer(t, "foo", dir),
			testLayer(t, "bar", dir),
		},
	}

	ctx := testContextShared(t)
	if err := layer.Build(ctx); err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(runner.Commands) != 6 {
		t.Fatalf("bad: %#v", runner.Commands)
	}

	// Repeat the test since this should be a no-op
	if err := layer.Build(ctx); err != nil {
		t.Fatalf("err: %s", err)
	}
	if len(runner.Commands) != 6 {
		t.Fatalf("bad: %#v", runner.Commands)
	}
}

func TestLayeredPending_new(t *testing.T) {
	dir := tempDir(t)
	defer os.RemoveAll(dir)

	runner := new(exec.MockRunner)
	defer exec.TestChrunner(runner.Run)()

	layer := &Layered{
		DataDir: dir,
		Layers: []*Layer{
			testLayer(t, "foo", dir),
			testLayer(t, "bar", dir),
		},
	}

	pending, err := layer.Pending()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []string{"foo", "bar"}
	if !reflect.DeepEqual(pending, expected) {
		t.Fatalf("bad: %#v", pending)
	}
}

func TestLayeredPending_partial(t *testing.T) {
	dir := tempDir(t)
	defer os.RemoveAll(dir)

	runner := new(exec.MockRunner)
	defer exec.TestChrunner(runner.Run)()

	// Build the foo layer
	layer := &Layered{
		DataDir: dir,
		Layers: []*Layer{
			testLayer(t, "foo", dir),
		},
	}

	ctx := testContextShared(t)
	if err := layer.Build(ctx); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Add the bar layer
	layer.Layers = append(layer.Layers, testLayer(t, "bar", dir))

	// Grab the pending
	pending, err := layer.Pending()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []string{"bar"}
	if !reflect.DeepEqual(pending, expected) {
		t.Fatalf("bad: %#v", pending)
	}
}

func TestLayeredPrune(t *testing.T) {
	dir := tempDir(t)
	defer os.RemoveAll(dir)

	runner := new(exec.MockRunner)
	defer exec.TestChrunner(runner.Run)()

	// Build an environment using foo and bar
	layer := &Layered{
		DataDir: dir,
		Layers: []*Layer{
			testLayer(t, "foo", dir),
			testLayer(t, "bar", dir),
		},
	}

	ctx := testContextShared(t)
	if err := layer.Build(ctx); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Add an environment
	env := &Vagrant{DataDir: filepath.Join(dir, "v1")}
	if err := layer.AddEnv(env); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Clear the commands because we only want to count those
	runner.Commands = nil

	// Prune should not do anything
	if _, err := layer.Prune(ctx); err != nil {
		t.Fatalf("err: %s", err)
	}
	if len(runner.Commands) > 0 {
		t.Fatalf("bad: %#v", runner.Commands)
	}
}

func TestLayeredPrune_empty(t *testing.T) {
	dir := tempDir(t)
	defer os.RemoveAll(dir)

	runner := new(exec.MockRunner)
	defer exec.TestChrunner(runner.Run)()

	layer := &Layered{
		DataDir: dir,
		Layers: []*Layer{
			testLayer(t, "foo", dir),
			testLayer(t, "bar", dir),
		},
	}

	if _, err := layer.Prune(testContextShared(t)); err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(runner.Commands) > 0 {
		t.Fatalf("bad: %#v", runner.Commands)
	}
}

func TestLayeredPrune_all(t *testing.T) {
	dir := tempDir(t)
	defer os.RemoveAll(dir)

	runner := new(exec.MockRunner)
	defer exec.TestChrunner(runner.Run)()

	layer := &Layered{
		DataDir: dir,
		Layers: []*Layer{
			testLayer(t, "foo", dir),
			testLayer(t, "bar", dir),
		},
	}

	ctx := testContextShared(t)
	if err := layer.Build(ctx); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Clear the commands because we only want to count those
	runner.Commands = nil
	if _, err := layer.Prune(ctx); err != nil {
		t.Fatalf("err: %s", err)
	}

	// 2 vagrant destroy -f
	if len(runner.Commands) != 2 {
		t.Fatalf("bad: %#v", runner.Commands)
	}
}

func testContextShared(t *testing.T) *context.Shared {
	return &context.Shared{
		Ui: &ui.Logged{Ui: new(ui.Mock)},
	}
}

func testLayer(t *testing.T, id string, dir string) *Layer {
	dir = filepath.Join(dir, id)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("err: %s", err)
	}

	vagrantfile := filepath.Join(dir, "Vagrantfile")
	if err := ioutil.WriteFile(vagrantfile, []byte("hello"), 0644); err != nil {
		t.Fatalf("err: %s", err)
	}

	return &Layer{
		ID:          id,
		Vagrantfile: vagrantfile,
	}
}

func tempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return dir
}
