package flag

import (
	"flag"
)

// FilterArgs filters the args slice to only include the the flags
// in the given flagset and returns a new arg slice that has the
// included args as well as a slice that has only the excluded args.
//
// Any positional arguments are added to BOTH slices.
func FilterArgs(fs *flag.FlagSet, args []string) ([]string, []string) {
	// Optimistically make bothy the length of the arguments. There
	// should never be so many arguments that this is too ineffecient.
	inc := make([]string, 0, len(args))
	exc := make([]string, 0, len(args))
	pos := make([]string, 0, len(args))

	// Make a map of the valid flags
	flags := make(map[string]struct{})
	fs.VisitAll(func(f *flag.Flag) {
		flags[f.Name] = struct{}{}
	})

	// Go through each, parse out a single argument, and determine where
	// it should go plus how many of the slots.
	i := 0
	for i < len(args) {
		n, loc := filterOne(flags, args, i)

		// Determine what slice to add the values to
		var result *[]string
		switch loc {
		case filterLocBoth:
			result = &pos
		case filterLocInc:
			result = &inc
		case filterLocExc:
			result = &exc
		}

		// Copy the values
		*result = append(*result, args[i:i+n]...)

		// Increment i so we continue moving through the arguments
		i += n
	}

	// Copy the positional elements onto both
	inc = append(inc, pos...)
	exc = append(exc, pos...)

	return inc, exc
}

type filterLoc byte

const (
	filterLocBoth filterLoc = iota
	filterLocInc
	filterLocExc
)

// filterOne is based very heavily on the official flag package
// "parseOne" function. We do this on purpose so that we parse things
// as similarly as possible in order to split the args.
func filterOne(flags map[string]struct{}, args []string, i int) (int, filterLoc) {
	// Get the arg
	s := args[i]

	// If the arg is empty, not a flag, or just a "-" then we have to
	// add it to BOTH lists.
	if len(s) == 0 || s[0] != '-' || len(s) == 1 {
		return 1, filterLocBoth
	}

	// If we hit double minuses, then we return the rest of the args to
	// BOTH lists.
	num_minuses := 1
	if s[1] == '-' {
		num_minuses++
		if len(s) == 2 { // "--" terminates the flags
			return len(args) - i, filterLocBoth
		}
	}

	// Otherwise, get the name. If the syntax is invalid, let's just add it
	// to both.
	name := s[num_minuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return 1, filterLocBoth
	}

	// Check for an argument to the flag
	has_value := false
	for i := 1; i < len(name); i++ { // equals cannot be first
		if name[i] == '=' {
			has_value = true
			name = name[0:i]
			break
		}
	}

	// Determine where this will go from here on out
	pos := filterLocInc
	if _, valid := flags[name]; !valid {
		pos = filterLocExc
	}

	// It must have a value, which might be the next argument.
	if !has_value && len(args) > i {
		return 2, pos
	}

	return 1, pos
}
