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
	// App
	//---------------------------------------------------------------

	{
		// GetApp (doesn't exist)
		lookup := AppLookup{AppID: "42", Version: "1.2.3", ConfigHash: 42}
		app, err := b.GetApp(&lookup)
		if err != nil {
			t.Errorf("GetApp error: %s", err)
			return
		}
		if app != nil {
			t.Errorf("App shouldn't be found: %#v", app)
			return
		}

		// PutApp
		expected := &App{Name: "foo", Type: "bar"}
		err = b.PutApp(&lookup, expected)
		if err != nil {
			t.Errorf("PutApp error: %s", err)
			return
		}

		// GetApp (exists)
		app, err = b.GetApp(&lookup)
		if err != nil {
			t.Errorf("GetApp error: %s", err)
			return
		}
		if !reflect.DeepEqual(expected, app) {
			t.Errorf(
				"GetApp doesn't match. Expected, then actual:\n\n%#v\n\n%#v",
				expected, app)
			return
		}

		// ListApps
		expected.AppLookup = lookup
		apps, err := b.ListApps()
		if err != nil {
			t.Errorf("ListApps error: %s", err)
			return
		}
		if len(apps) != 1 {
			t.Errorf("ListApps length should be 1: %d", len(apps))
			return
		}
		if !reflect.DeepEqual(expected, apps[0]) {
			t.Errorf(
				"ListApps doesn't match. Expected, then actual:\n\n%#v\n\n%#v",
				expected, apps[0])
			return
		}
	}

	//---------------------------------------------------------------
	// Infra
	//---------------------------------------------------------------

	{
		// GetInfra (doesn't exist)
		lookup := &InfraLookup{Name: "foo"}
		infra, err := b.GetInfra(lookup)
		if err != nil {
			t.Errorf("GetInfra (non-exist) error: %s", err)
			return
		}
		if infra != nil {
			t.Error("GetInfra (non-exist): infra should be nil")
			return
		}

		// PutInfra (doesn't exist)
		expected := &Infra{Name: "foo", Type: "bar"}
		err = b.PutInfra(lookup, expected)
		if err != nil {
			t.Errorf("PutInfra error: %s", err)
			return
		}

		// GetInfra (exists)
		infra, err = b.GetInfra(lookup)
		if err != nil {
			t.Errorf("GetInfra error: %s", err)
			return
		}
		if !reflect.DeepEqual(expected, infra) {
			t.Errorf(
				"GetInfra doesn't match. Expected, then actual:\n\n%#v\n\n%#v",
				expected, infra)
			return
		}
	}

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
