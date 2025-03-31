package toposort

import (
	"gitlab.lrz.de/hm/goal-core/hmgraph"
)

// TopologicalOrder computes any topological vertex order.
// It will return an error if either the graph contains any
// edge or if the arcs form a cycle.
func TopologicalOrder(g *hmgraph.Graph) (sorting []*hmgraph.Vertex, err error) {
	return nil, nil
}
