package plan

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestExecutorValidate(t *testing.T) {
	cases := []struct {
		Name string
		Err  bool
	}{
		{
			"validate-valid.hcl",
			false,
		},

		{
			"validate-invalid-type.hcl",
			true,
		},

		{
			"validate-invalid-arg.hcl",
			true,
		},

		{
			"validate-invalid-result-ref.hcl",
			true,
		},

		{
			"validate-invalid-store-ref.hcl",
			true,
		},
	}

	testTaskMap := map[string]TaskExecutor{
		"delete": &DeleteTask{},
		"store":  &StoreTask{},

		"test-1": &MockTaskExecutor{
			ValidateResult: &ExecResult{
				Values: map[string]*TaskResult{
					"Result": nil,
				},
			},
		},

		"test-err": &MockTaskExecutor{
			ValidateErr: errors.New("error"),
		},
	}

	for _, tc := range cases {
		path := filepath.Join("./test-fixtures", tc.Name)
		plans, err := ParseFile(path)
		if err != nil {
			t.Fatalf("%s, err: %s", tc.Name, err)
		}

		exec := &Executor{TaskMap: testTaskMap}
		for _, p := range plans {
			err := exec.Validate(p, nil)
			if (err != nil) != tc.Err {
				t.Fatalf("%s, err: %s", tc.Name, err)
			}
		}
	}
}
