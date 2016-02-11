package plan

import (
	"bytes"
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
			"validate-invalid-interp.hcl",
			true,
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

		"test-1": &MockTask{
			ValidateResult: &ExecResult{
				Values: map[string]*TaskResult{
					"Result": nil,
				},
			},
		},

		"test-err": &MockTask{
			ValidateErr: errors.New("error"),
		},
	}

	for _, tc := range cases {
		t.Logf("Testing: %s", tc.Name)

		path := filepath.Join("./test-fixtures", tc.Name)
		plans, err := ParseFile(path)
		if err != nil {
			t.Fatalf("%s, err: %s", tc.Name, err)
		}

		exec := &Executor{TaskMap: testTaskMap}
		for _, p := range plans {
			err := exec.Validate(p)
			if (err != nil) != tc.Err {
				t.Fatalf("%s, err: %s", tc.Name, err)
			}
		}
	}
}

func TestExecutorExecute(t *testing.T) {
	cases := []struct {
		Name   string
		Err    bool
		Result string
	}{
		{
			"execute-basic.hcl",
			false,
			"hello",
		},

		{
			"execute-storage.hcl",
			false,
			"hello",
		},
	}

	task := &testTask{}
	testTaskMap := map[string]TaskExecutor{
		"delete": &DeleteTask{},
		"store":  &StoreTask{},
		"test":   task,
	}

	for _, tc := range cases {
		path := filepath.Join("./test-fixtures", tc.Name)
		plans, err := ParseFile(path)
		if err != nil {
			t.Fatalf("%s, err: %s", tc.Name, err)
		}

		exec := &Executor{TaskMap: testTaskMap}
		for _, p := range plans {
			err := exec.Execute(p)
			if (err != nil) != tc.Err {
				t.Fatalf("%s, err: %s", tc.Name, err)
			}

			if task.Result != tc.Result {
				t.Fatalf("%s, bad: %s", tc.Name, task.Result)
			}
		}
	}
}

func TestExecutorOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	fn := func(args *ExecArgs) (*ExecResult, error) {
		args.Println("HELLO!")
		return nil, nil
	}

	task := &MockTask{ValidateFn: fn, ExecuteFn: fn}
	testTaskMap := map[string]TaskExecutor{
		"test": task,
	}

	path := filepath.Join("./test-fixtures", "execute-output.hcl")
	plans, err := ParseFile(path)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	exec := &Executor{Output: buf, TaskMap: testTaskMap}
	for _, p := range plans {
		err := exec.Validate(p)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		// Should not have output
		if buf.Len() > 0 {
			t.Fatalf("bad: %s", buf.String())
		}
	}

	for _, p := range plans {
		err := exec.Execute(p)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		// Should have output
		if buf.Len() == 0 {
			t.Fatalf("bad: %s", buf.String())
		}
	}
}

type testTask struct {
	Result string
}

func (t *testTask) Validate(args *ExecArgs) (*ExecResult, error) {
	return nil, nil
}

func (t *testTask) Execute(args *ExecArgs) (*ExecResult, error) {
	if arg, ok := args.Args["result"]; ok {
		t.Result = arg.Value.(string)
	}

	return nil, nil
}
