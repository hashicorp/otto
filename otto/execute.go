package otto

// ExecuteTask is an enum of available tasks to execute.
type ExecuteTask uint

const (
	ExecuteTaskInvalid ExecuteTask = 0
	ExecuteTaskInfra   ExecuteTask = iota
)

// ExecuteOpts are the options used for executing generic tasks
// on the Otto environment.
type ExecuteOpts struct {
	// Task is the task to execute
	Task ExecuteTask

	// Args are additional arguments to the task
	Args []string
}
