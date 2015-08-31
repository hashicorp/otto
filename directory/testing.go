package directory

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

// TestBackend is a public test helper that verifies a backend
// functions properly.
func TestBackend(t *testing.T, b Backend) {
	var buf bytes.Buffer
	var err error

	//---------------------------------------------------------------
	// Blob
	//---------------------------------------------------------------

	// PutBlob
	err = b.PutBlob("foo", &BlobData{Data: strings.NewReader("bar")})
	if err != nil {
		t.Fatalf("PutBlob error: %s", err)
	}

	// GetBlob (exists)
	data, err := b.GetBlob("foo")
	if err != nil {
		t.Fatalf("GetBlob error: %s", err)
	}
	_, err = io.Copy(&buf, data.Data)
	data.Close()
	if err != nil {
		t.Fatalf("GetBlob error: %s", err)
	}
	if buf.String() != "bar" {
		t.Fatalf("GetBlob bad data: %s", buf.String())
	}

	//---------------------------------------------------------------
	// Infra
	//---------------------------------------------------------------

	// GetInfra (doesn't exist)
	infra, err := b.GetInfra("foo")
	if err != nil {
		t.Fatalf("GetInfra (non-exist) error: %s", err)
	}
	if infra != nil {
		t.Fatal("GetInfra (non-exist): infra should be nil")
	}

	// PutInfra (doesn't exist)
	infra = &Infra{Outputs: map[string]string{"foo": "bar"}}
	if err := b.PutInfra("foo", infra); err != nil {
		t.Fatalf("PutInfra err: %s", err)
	}

	// GetInfra (exists)
	actualInfra, err := b.GetInfra("foo")
	if err != nil {
		t.Fatalf("GetInfra (exist) error: %s", err)
	}
	if !reflect.DeepEqual(actualInfra, infra) {
		t.Fatalf("GetInfra (exist) bad: %#v", actualInfra)
	}

	//---------------------------------------------------------------
	// Deploy
	//---------------------------------------------------------------

	// GetDeploy (doesn't exist)
	deploy := &Deploy{App: "foo", Infra: "bar", InfraFlavor: "baz"}
	deployResult, err := b.GetDeploy(deploy)
	if err != nil {
		t.Fatalf("GetDeploy (non-exist) error: %s", err)
	}
	if deployResult != nil {
		t.Fatal("GetDeploy (non-exist): result should be nil")
	}

	// PutDeploy (doesn't exist)
	if deploy.ID != "" {
		t.Fatalf("PutDeploy: ID should be empty before set")
	}
	if err := b.PutDeploy(deploy); err != nil {
		t.Fatalf("PutDeploy err: %s", err)
	}
	if deploy.ID == "" {
		t.Fatalf("PutDeploy: deploy ID not set")
	}

	// GetDeploy (exists)
	deployResult, err = b.GetDeploy(deploy)
	if err != nil {
		t.Fatalf("GetDeploy (exist) error: %s", err)
	}
	if !reflect.DeepEqual(deployResult, deploy) {
		t.Fatalf("GetDeploy (exist) bad: %#v", deployResult)
	}
}
