package graph

import "errors"

type Chain []int

var (
	ErrChainLength       = errors.New("chain must have at least two vertices")
	ErrVertexNotFound    = errors.New("vertex not found in graph")
	ErrNeighbourNotFound = errors.New("neighbour not found in chain")
)

func (g *Graph[T, E]) Verify(c Chain) error {
	if len(c) < 2 {
		return ErrChainLength
	}

	v, exists := g.Vertices[c[0]]
	if !exists {
		return ErrVertexNotFound
	}

	for _, id := range c[1:] {
		e, exists := v.Edges[id]
		if !exists {
			return ErrNeighbourNotFound
		}
		v = e.Vertex
	}

	return nil
}

func (g *Graph[T, E]) VerifyWithEdges(c Chain) ([]*Edge[T, E], error) {
	edges := make([]*Edge[T, E], 0, len(c)-1)
	if len(c) < 2 {
		return edges, ErrChainLength
	}

	v, exists := g.Vertices[c[0]]
	if !exists {
		return edges, ErrVertexNotFound
	}

	for _, id := range c[1:] {
		e, exists := v.Edges[id]
		if !exists {
			return edges, ErrNeighbourNotFound
		}
		edges = append(edges, e)
		v = e.Vertex
	}

	return edges, nil
}
