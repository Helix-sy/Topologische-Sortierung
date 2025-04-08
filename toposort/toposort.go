package toposort

import (
	"fmt"
	"gitlab.lrz.de/hm/goal-core/hmgraph"
)

func TopologicalOrder(g *hmgraph.Graph) ([]*hmgraph.Vertex, error) {
	vertices := g.GetVertices()

	//Check for undirected edges
	hasEdge := false
	for _, v := range vertices {
		v.ForEdges(func(e *hmgraph.Edge) {
			hasEdge = true
		})
	}
	if hasEdge {
		return nil, fmt.Errorf("graph contains undirected edges")
	}

	inDegree := hmgraph.CreateVertexMap(g, "inDegree", 0)
	defer inDegree.Dispose()

	for _, v := range vertices {
		v.ForOutArcs(func(arc *hmgraph.Arc) {
			target := arc.Target()
			inDegree.Set(target, inDegree.Get(target)+1)
		})
	}

	var queue []*hmgraph.Vertex
	for _, v := range vertices {
		if inDegree.Get(v) == 0 {
			queue = append(queue, v)
		}
	}

	var sorting []*hmgraph.Vertex
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		sorting = append(sorting, v)

		v.ForOutArcs(func(arc *hmgraph.Arc) {
			target := arc.Target()
			newIn := inDegree.Get(target) - 1
			inDegree.Set(target, newIn)
			if newIn == 0 {
				queue = append(queue, target)
			}
		})
	}

	if len(sorting) != len(vertices) {
		return nil, fmt.Errorf("graph contains a cycle")
	}

	return sorting, nil
}
