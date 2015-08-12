package ui

import (
	"github.com/mitchellh/colorstring"
)

var colorstringDisable = colorstring.Colorize{
	Colors:  colorstring.DefaultColors,
	Disable: true,
	Reset:   false,
}

// Colorize is a helper to colorize the string according to the colorstring
// library defaults.
func Colorize(t string) string {
	return colorstring.Color(t)
}

// StripColors is a helper to strip all the color tags from the text.
func StripColors(t string) string {
	return colorstringDisable.Color(t)
}
