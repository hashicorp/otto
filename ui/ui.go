package ui

import "os"

// Ui is the component of Otto responsible for reading/writing to the
// console.
//
// All Ui implementations MUST expect colorstring[1] style inputs. If
// the output interface doesn't support colors, these must be stripped.
// The StripColors helper in this package can be used to do this. For
// terminals, the Colorize helper in this package can be used.
//
// [1]: github.com/mitchellh/colorstring
type Ui interface {
	// Header, Message, and Raw are all methods for outputting messages
	// to the Ui, all with different styles. Header and Message should
	// be used liberally, and Raw should be scarcely used if possible.
	//
	// Header outputs the message with a style denoting it is a sectional
	// message. An example: "==> TEXT" might be header text.
	//
	// Message outputs the message with a style, but one that is less
	// important looking than Header. Example: "    TEXT" might be
	// the message text, prefixed with spaces so that it lines up with
	// header items.
	//
	// Raw outputs a message with no styling.
	Header(string)
	Message(string)
	Raw(string)

	// Input is used to ask the user for input. This should be done
	// sparingly or very early on since Otto is meant to be an automated
	// tool.
	Input(*InputOpts) (string, error)
}

// InputOpts are options for asking for input.
type InputOpts struct {
	// Id is a unique ID for the question being asked that might be
	// used for logging or to look up a prior answered question.
	Id string

	// Query is a human-friendly question for inputting this value.
	Query string

	// Description is a description about what this option is. Be wary
	// that this will probably be in a terminal so split lines as you see
	// necessary.
	Description string

	// Default will be the value returned if no data is entered.
	Default string

	// Hide will hide the text while it is being typed.
	Hide bool

	// EnvVars is a list of environment variables where the value can be looked
	// up, in priority order. If any of these environment Variables are
	// non-empty, they will be returned as the value for this input and the user
	// will not be prompted.
	EnvVars []string
}

// EnvVarValue reads the configured list of EnvVars, returns the first
// non-empty value it finds, otherwise returns an empty string.
func (o *InputOpts) EnvVarValue() string {
	for _, envVar := range o.EnvVars {
		if val := os.Getenv(envVar); val != "" {
			return val
		}
	}
	return ""
}
