package toposort

import (
	"errors"
	"gitlab.lrz.de/hm/goal-core/hmgraph"
)

// TopologicalOrder computes any topological vertex order.
// It will return an error if either the graph contains any
// edge or if the arcs form a cycle.
func TopologicalOrder(g *hmgraph.Graph) (sorting []*hmgraph.Vertex, err error) {
	// Check if the graph contains any undirected edges
	hasEdge := false
	g.ForVertices(func(v *hmgraph.Vertex) {
		v.ForEdges(func(e *hmgraph.Edge) {
			hasEdge = true
		})
	})
	if hasEdge {
		return nil, errors.New("graph contains undirected edges")
	}

	// Count vertices for initialization
	vertexCount := 0
	g.ForVertices(func(v *hmgraph.Vertex) {
		vertexCount++
	})

	// Initialize data structures for DFS
	visited := hmgraph.CreateVertexMap(g, "visited", false)
	defer visited.Dispose()

	temp := hmgraph.CreateVertexMap(g, "temp", false)
	defer temp.Dispose()

	order := make([]*hmgraph.Vertex, 0, vertexCount)

	// DFS function to detect cycles and build topological order
	var dfs func(*hmgraph.Vertex) error
	dfs = func(v *hmgraph.Vertex) error {
		// If this vertex is already in the temporary mark, we have a cycle
		if temp.Get(v) == true {
			return errors.New("graph contains a cycle")
		}

		// If we haven't visited this vertex yet
		if visited.Get(v) == false {
			// Mark temporarily for cycle detection
			temp.Set(v, true)

			// Visit all successors (outgoing arcs)
			var cycleError error
			v.ForOutArcs(func(arc *hmgraph.Arc) {
				if cycleError == nil {
					cycleError = dfs(arc.Target())
				}
			})

			// If a cycle was detected, propagate the error
			if cycleError != nil {
				return cycleError
			}

			// Mark as permanently visited
			visited.Set(v, true)
			// Remove temporary mark
			temp.Set(v, false)
			// Add to ordering (prepend)
			order = append([]*hmgraph.Vertex{v}, order...)
		}

		return nil
	}

	// Run DFS from each unvisited vertex
	var finalError error
	g.ForVertices(func(v *hmgraph.Vertex) {
		if finalError == nil && visited.Get(v) == false {
			if err := dfs(v); err != nil {
				finalError = err
			}
		}
	})

	if finalError != nil {
		return nil, finalError
	}

	return order, nil
}
