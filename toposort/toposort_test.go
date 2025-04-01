package toposort

import (
	"github.com/stretchr/testify/assert"
	"gitlab.lrz.de/hm/goal-core/hmgraph"
	"testing"
)

func TestSmallStandardCase(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(6)
	vs[3].CreateArc(vs[1])
	vs[1].CreateArc(vs[4])
	vs[0].CreateArc(vs[5])
	vs[2].CreateArcs([]*hmgraph.Vertex{vs[4], vs[1]})
	sorting, err := TopologicalOrder(g)
	assert.NoError(t, err)
	assertValidOrdering(g, sorting, t)
	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed.")
}

func TestErrorOnEdge(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(6)
	vs[3].CreateArc(vs[1])
	vs[1].CreateArc(vs[4])
	vs[0].CreateArc(vs[5])
	vs[2].CreateArcs([]*hmgraph.Vertex{vs[4], vs[1]})
	vs[5].CreateEdge(vs[4])
	_, err := TopologicalOrder(g)
	assert.Error(t, err)
	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed on error.")
}

func TestErrorOnCycle(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(6)
	vs[3].CreateArc(vs[1])
	vs[1].CreateArc(vs[4])
	vs[0].CreateArc(vs[5])
	vs[2].CreateArcs([]*hmgraph.Vertex{vs[4], vs[1]})
	vs[4].CreateArc(vs[3])
	_, err := TopologicalOrder(g)
	assert.Error(t, err)
	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed on cycle.")
}

func assertValidOrdering(g *hmgraph.Graph, sorting []*hmgraph.Vertex, t *testing.T) {
	numbering := hmgraph.CreateVertexMap(g, "testNumber", -1)
	for i, vertex := range sorting {
		numbering.Set(vertex, i)
		vertex.ForOutArcs(func(arc *hmgraph.Arc) {
			assert.True(t, numbering.Get(arc.Target()) == -1)
		})
	}
	numbering.Dispose()
}
