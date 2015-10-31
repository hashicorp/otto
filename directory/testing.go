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

	// Through this method we use "Errorf" instead of "Fatalf". This is
	// important and deliberate: for our RPC tests, Fatalf causes the
	// test to hang in a failure due to our RPC model. Errorf causes it
	// to end properly.

	//---------------------------------------------------------------
	// Blob
	//---------------------------------------------------------------

	// GetBlob (doesn't exist)
	data, err := b.GetBlob("foo")
	if err != nil {
		t.Errorf("GetBlob error: %s", err)
		return
	}
	if data != nil {
		data.Close()
		t.Errorf("GetBlob should be nil data")
		return
	}

	// PutBlob
	err = b.PutBlob("foo", &BlobData{Data: strings.NewReader("bar")})
	if err != nil {
		t.Errorf("PutBlob error: %s", err)
		return
	}

	// GetBlob (exists)
	data, err = b.GetBlob("foo")
	if err != nil {
		t.Errorf("GetBlob error: %s", err)
		return
	}
	_, err = io.Copy(&buf, data.Data)
	data.Close()
	if err != nil {
		t.Errorf("GetBlob error: %s", err)
		return
	}
	if buf.String() != "bar" {
		t.Errorf("GetBlob bad data: %s", buf.String())
		return
	}

	//---------------------------------------------------------------
	// Infra
	//---------------------------------------------------------------

	// GetInfra (doesn't exist)
	infra := &Infra{Lookup: Lookup{Infra: "foo"}}
	actualInfra, err := b.GetInfra(infra)
	if err != nil {
		t.Errorf("GetInfra (non-exist) error: %s", err)
		return
	}
	if actualInfra != nil {
		t.Error("GetInfra (non-exist): infra should be nil")
		return
	}

	// PutInfra (doesn't exist)
	infra.Outputs = map[string]string{"foo": "bar"}
	if infra.ID != "" {
		t.Errorf("PutInfra: ID should be empty before set")
		return
	}
	if err := b.PutInfra(infra); err != nil {
		t.Errorf("PutInfra err: %s", err)
		return
	}
	if infra.ID == "" {
		t.Errorf("PutInfra: infra ID not set")
		return
	}

	// GetInfra (exists)
	actualInfra, err = b.GetInfra(infra)
	if err != nil {
		t.Errorf("GetInfra (exist) error: %s", err)
		return
	}
	if !reflect.DeepEqual(actualInfra, infra) {
		t.Errorf("GetInfra (exist) bad: %#v", actualInfra)
		return
	}

	// GetInfra with foundation (doesn't exist)
	infra = &Infra{Lookup: Lookup{Infra: "foo", Foundation: "bar"}}
	actualInfra, err = b.GetInfra(infra)
	if err != nil {
		t.Errorf("GetInfra (non-exist) error: %s", err)
		return
	}
	if actualInfra != nil {
		t.Error("GetInfra (non-exist): infra should be nil")
		return
	}

	// PutInfra with foundation (doesn't exist)
	infra.Outputs = map[string]string{"foo": "bar"}
	if infra.ID != "" {
		t.Errorf("PutInfra: ID should be empty before set")
		return
	}
	if err := b.PutInfra(infra); err != nil {
		t.Errorf("PutInfra err: %s", err)
		return
	}
	if infra.ID == "" {
		t.Errorf("PutInfra: infra ID not set")
		return
	}

	// GetInfra with foundation (exists)
	actualInfra, err = b.GetInfra(infra)
	if err != nil {
		t.Errorf("GetInfra (exist) error: %s", err)
		return
	}
	if !reflect.DeepEqual(actualInfra, infra) {
		t.Errorf("GetInfra (exist) bad: %#v", actualInfra)
		return
	}

	//---------------------------------------------------------------
	// Deploy
	//---------------------------------------------------------------

	// GetDeploy (doesn't exist)
	deploy := &Deploy{Lookup: Lookup{
		AppID: "foo", Infra: "bar", InfraFlavor: "baz"}}
	deployResult, err := b.GetDeploy(deploy)
	if err != nil {
		t.Errorf("GetDeploy (non-exist) error: %s", err)
		return
	}
	if deployResult != nil {
		t.Error("GetDeploy (non-exist): result should be nil")
		return
	}

	// PutDeploy (doesn't exist)
	if deploy.ID != "" {
		t.Errorf("PutDeploy: ID should be empty before set")
		return
	}
	if err := b.PutDeploy(deploy); err != nil {
		t.Errorf("PutDeploy err: %s", err)
		return
	}
	if deploy.ID == "" {
		t.Errorf("PutDeploy: deploy ID not set")
		return
	}

	// GetDeploy (exists)
	deployResult, err = b.GetDeploy(deploy)
	if err != nil {
		t.Errorf("GetDeploy (exist) error: %s", err)
		return
	}
	if !reflect.DeepEqual(deployResult, deploy) {
		t.Errorf("GetDeploy (exist) bad: %#v", deployResult)
		return
	}

	//---------------------------------------------------------------
	// Dev
	//---------------------------------------------------------------

	// GetDev (doesn't exist)
	dev := &Dev{Lookup: Lookup{AppID: "foo"}}
	devResult, err := b.GetDev(dev)
	if err != nil {
		t.Errorf("GetDev (non-exist) error: %s", err)
		return
	}
	if devResult != nil {
		t.Error("GetDev (non-exist): result should be nil")
		return
	}

	// PutDev (doesn't exist)
	if dev.ID != "" {
		t.Errorf("PutDev: ID should be empty before set")
		return
	}
	if err := b.PutDev(dev); err != nil {
		t.Errorf("PutDev err: %s", err)
		return
	}
	if dev.ID == "" {
		t.Errorf("PutDev: dev ID not set")
		return
	}

	// GetDev (exists)
	devResult, err = b.GetDev(dev)
	if err != nil {
		t.Errorf("GetDev (exist) error: %s", err)
		return
	}
	if !reflect.DeepEqual(devResult, dev) {
		t.Errorf("GetDev (exist) bad: %#v", devResult)
		return
	}

	// DeleteDev (exists)
	err = b.DeleteDev(dev)
	if err != nil {
		t.Errorf("DeleteDev error: %s", err)
		return
	}
	devResult, err = b.GetDev(dev)
	if err != nil {
		t.Errorf("GetDev (non-exist) error: %s", err)
		return
	}
	if devResult != nil {
		t.Error("GetDev (non-exist): result should be nil")
		return
	}

	// DeleteDev (doesn't exist)
	err = b.DeleteDev(dev)
	if err != nil {
		t.Errorf("DeleteDev error: %s", err)
		return
	}
}
