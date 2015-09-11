package otto

// ExecuteTask is an enum of available tasks to execute.
type ExecuteTask uint

const (
	ExecuteTaskInvalid ExecuteTask = 0
	ExecuteTaskInfra   ExecuteTask = iota
	ExecuteTaskDev
)

//go:generate stringer -type=ExecuteTask

// ExecuteOpts are the options used for executing generic tasks
// on the Otto environment.
type ExecuteOpts struct {
	// Task is the task to execute
	Task ExecuteTask

	// Action is a sub-action that a task can take. For example:
	// infrastructures accept "destroy", development environments
	// accept "reload", etc.
	Action string

	// Args are additional arguments to the task
	Args []string
}
