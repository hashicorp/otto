package pongo2_ext

import (
	"io/ioutil"

	"github.com/flosch/pongo2"
)

func init() {
	pongo2.RegisterFilter("read", filterRead)
}

func filterRead(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	data, err := ioutil.ReadFile(in.String())
	if err != nil {
		return nil, &pongo2.Error{
			Sender:   "filter:read",
			ErrorMsg: err.Error(),
		}
	}

	return pongo2.AsSafeValue(string(data)), nil
}
