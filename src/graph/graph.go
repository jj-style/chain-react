package graph

import "github.com/samber/lo"

// Graph represents a set of vertices connected by edges.
type Graph[T, E any] struct {
	Vertices map[int]*Vertex[T, E]
}

// Vertex is a node in the graph that stores the int value at that node
// along with a map to the vertices it is connected to via edges.
type Vertex[T, E any] struct {
	Val   T
	Edges map[int]*Edge[T, E]
}

// Edge represents an edge in the graph and the destination vertex.
type Edge[T, E any] struct {
	Weight E
	Vertex *Vertex[T, E]
}

func (g *Graph[T, E]) AddVertex(key int, val T) {
	g.Vertices[key] = &Vertex[T, E]{Val: val, Edges: map[int]*Edge[T, E]{}}
}

func (g *Graph[T, E]) AddEdge(srcKey, destKey int, weight E) {
	// check if src & dest exist
	if _, ok := g.Vertices[srcKey]; !ok {
		return
	}
	if _, ok := g.Vertices[destKey]; !ok {
		return
	}

	// add edge src --> dest
	g.Vertices[srcKey].Edges[destKey] = &Edge[T, E]{Weight: weight, Vertex: g.Vertices[destKey]}
}

func (g *Graph[T, E]) Neighbors(srcKey int) []T {
	return lo.MapToSlice(g.Vertices[srcKey].Edges, func(key int, value *Edge[T, E]) T {
		return value.Vertex.Val
	})
}

func (g *Graph[T, E]) NeighboursMap(srcKey int) map[int]*Edge[T, E] {
	return g.Vertices[srcKey].Edges
}

func (g *Graph[T, E]) HasVertex(key int) bool {
	_, exists := g.Vertices[key]
	return exists
}

func NewGraph[T, E any](opts ...GraphOption[T, E]) *Graph[T, E] {
	g := &Graph[T, E]{Vertices: map[int]*Vertex[T, E]{}}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

type GraphOption[T, E any] func(g *Graph[T, E])
