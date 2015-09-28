package appfile

import (
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform/dag"
)

// This file contains the MarshalJSON and UnmarshalJSON functions so that
// the Compiled struct can be safely encoded/decoded on disk. We put this
// in a separate file since it requires a bunch of auxilliary structs that
// we didn't want to confuse compile.go with.

func (c *Compiled) MarshalJSON() ([]byte, error) {
	raw := &compiledJSON{
		File:  c.File,
		Edges: make([]map[string]string, 0, len(c.Graph.Edges())),
	}

	// Compile the list of vertices, keeping track of their position
	set := make(map[dag.Vertex]string)
	for i, rawV := range c.Graph.Vertices() {
		v := rawV.(*CompiledGraphVertex)
		raw.Vertices = append(raw.Vertices, v)
		set[v] = strconv.FormatInt(int64(i), 10)
	}

	// Map the edges by position
	for _, e := range c.Graph.Edges() {
		raw.Edges = append(raw.Edges,
			map[string]string{
				set[e.Source()]: set[e.Target()],
			})
	}

	return json.Marshal(raw)
}

func (c *Compiled) UnmarshalJSON(data []byte) error {
	var raw compiledJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	c.File = raw.File
	c.Graph = new(dag.AcyclicGraph)
	for _, v := range raw.Vertices {
		c.Graph.Add(v)
	}
	for _, e := range raw.Edges {
		for a, b := range e {
			ai, err := strconv.ParseInt(a, 0, 0)
			if err != nil {
				return err
			}

			bi, err := strconv.ParseInt(b, 0, 0)
			if err != nil {
				return err
			}

			c.Graph.Connect(dag.BasicEdge(raw.Vertices[ai], raw.Vertices[bi]))
		}
	}

	return nil
}

type compiledJSON struct {
	File     *File
	Vertices []*CompiledGraphVertex
	Edges    []map[string]string
}
