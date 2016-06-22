package pongo2

import (
	"bytes"
	"fmt"
	"io"
)

type Template struct {
	set *TemplateSet

	// Input
	is_tpl_string bool
	name          string
	tpl           string
	size          int

	// Calculation
	tokens []*Token
	parser *Parser

	// first come, first serve (it's important to not override existing entries in here)
	level           int
	parent          *Template
	child           *Template
	blocks          map[string]*NodeWrapper
	exported_macros map[string]*tagMacroNode

	// Output
	root *nodeDocument
}

func newTemplateString(set *TemplateSet, tpl string) (*Template, error) {
	return newTemplate(set, "<string>", true, tpl)
}

func newTemplate(set *TemplateSet, name string, is_tpl_string bool, tpl string) (*Template, error) {
	// Create the template
	t := &Template{
		set:             set,
		is_tpl_string:   is_tpl_string,
		name:            name,
		tpl:             tpl,
		size:            len(tpl),
		blocks:          make(map[string]*NodeWrapper),
		exported_macros: make(map[string]*tagMacroNode),
	}

	// Tokenize it
	tokens, err := lex(name, tpl)
	if err != nil {
		return nil, err
	}
	t.tokens = tokens

	// For debugging purposes, show all tokens:
	/*for i, t := range tokens {
		fmt.Printf("%3d. %s\n", i, t)
	}*/

	// Parse it
	err = t.parse()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (tpl *Template) execute(context Context) (*bytes.Buffer, error) {
	// Create output buffer
	// We assume that the rendered template will be 30% larger
	buffer := bytes.NewBuffer(make([]byte, 0, int(float64(tpl.size)*1.3)))

	// Determine the parent to be executed (for template inheritance)
	parent := tpl
	for parent.parent != nil {
		parent = parent.parent
	}

	// Create context if none is given
	newContext := make(Context)
	newContext.Update(tpl.set.Globals)

	if context != nil {
		newContext.Update(context)

		if len(newContext) > 0 {
			// Check for context name syntax
			err := newContext.checkForValidIdentifiers()
			if err != nil {
				return nil, err
			}

			// Check for clashes with macro names
			for k, _ := range newContext {
				_, has := tpl.exported_macros[k]
				if has {
					return nil, &Error{
						Filename: tpl.name,
						Sender:   "execution",
						ErrorMsg: fmt.Sprintf("Context key name '%s' clashes with macro '%s'.", k, k),
					}
				}
			}
		}
	}

	// Create operational context
	ctx := newExecutionContext(parent, newContext)

	// Run the selected document
	err := parent.root.Execute(ctx, buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

// Executes the template with the given context and writes to writer (io.Writer)
// on success. Context can be nil. Nothing is written on error; instead the error
// is being returned.
func (tpl *Template) ExecuteWriter(context Context, writer io.Writer) error {
	buffer, err := tpl.execute(context)
	if err != nil {
		return err
	}

	l := buffer.Len()
	n, werr := buffer.WriteTo(writer)
	if int(n) != l {
		panic(fmt.Sprintf("error on writing template: n(%d) != buffer.Len(%d)", n, l))
	}
	if werr != nil {
		return &Error{
			Filename: tpl.name,
			Sender:   "execution",
			ErrorMsg: werr.Error(),
		}
	}
	return nil
}

// Executes the template and returns the rendered template as a []byte
func (tpl *Template) ExecuteBytes(context Context) ([]byte, error) {
	// Execute template
	buffer, err := tpl.execute(context)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Executes the template and returns the rendered template as a string
func (tpl *Template) Execute(context Context) (string, error) {
	// Execute template
	buffer, err := tpl.execute(context)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil

}
