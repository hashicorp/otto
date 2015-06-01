package appfile

// File is the structure of a single Appfile.
type File struct {
	Application    *Application
	Project        *Project
	Infrastructure []*Infrastructure
}

// Application is the structure of an application definition.
type Application struct {
	Name string
	Type string
}

// Project is the structure of a project that many applications
// can belong to.
type Project struct {
	Name  string
	Stack *Stack
}

// Infrastructure is the structure of defining the infrastructure
// that an application must run on.
type Infrastructure struct {
	Name   string
	Type   string
	Flavor string
}

// Stack is the structure that defines the stack that a project is
// built on.
type Stack struct {
	Name string
}
