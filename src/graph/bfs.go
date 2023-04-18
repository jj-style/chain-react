package graph

import (
	"github.com/jj-style/chain-react/src/utils"
	"github.com/samber/lo"
)

type Element[T, U any] struct {
	Id   int
	Edge *Edge[T, U]
}

func (g *Graph[T, U]) Bfs(src, dest int, found chan<- []Element[T, U]) {
	defer close(found)

	queue := [][]Element[T, U]{}

	path := []Element[T, U]{{Id: src, Edge: nil}}

	queue = append(queue, utils.Clone(path))

	visited := map[int]bool{} // track visited nodes to prevent cycles, if any

	for len(queue) > 0 {
		// dequeue
		path := queue[0]
		queue = queue[1:]

		// if at dest, we have a path, stop here
		last, _ := lo.Last(path)
		if last.Id == dest {
			found <- path
		}

		// for all adjacent vertices to last, if not visited, add to queue
		for adj, e := range g.NeighboursMap(last.Id) {
			if _, seen := visited[adj]; !seen {
				visited[adj] = true
				newpath := utils.Clone(path)
				newpath = append(newpath, Element[T, U]{Id: adj, Edge: e})
				queue = append(queue, newpath)
			}
		}
	}

}
