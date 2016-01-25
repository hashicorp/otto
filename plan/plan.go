// The plan package contains structures and helpers for Otto "plans,"
// the structure representing a goal and a set of tasks to achieve that
// goal.
package plan

// Plan is an executable object that represents a goal and the
// steps to take (tasks) to achieve that goal.
type Plan struct {
	Description string
	Tasks       []*Task
}

// Task is a single executable unit for a Plan. Tasks are meant to remain
// small in scope so that they can be composed and reasoned about within
// a plan.
type Task struct {
	Name string // Name of the task
	Type string // Type of the task

	Description         string // Short description of what this task will do
	DetailedDescription string // Long details about what this task will do (optional)
}
