package graph

import (
	hp "container/heap"

	"main/model"
)

func GetPath(g []model.Point, c []model.Point, obstacles []model.Polygon) ([]model.Point, float64) {
	m := len(g)

	graph := newGraph()

	for i := 0; i < len(g); i++ {
		for j := 0; j < len(c); j++ {
			segment := model.Segment{
				A: g[i],
				B: c[j],
			}

			flag := true
			for k := 0; k < len(obstacles); k++ {
				if obstacles[k].IntersectsWithSegment(segment) {
					flag = false
					break
				}
			}
			if flag {
				graph.addEdge(i, m+j, model.Dist(g[i], c[j]))
			}
		}
	}

	l, path := graph.getPath(0, 1)

	ret := make([]model.Point, 0, len(path))
	for i := 0; i < len(path); i++ {
		if path[i] < m {
			ret = append(ret, g[path[i]])
		} else {
			ret = append(ret, c[path[i]-m])
		}
	}

	return ret, l
}

type path struct {
	value float64
	nodes []int
}

type minPath []path

func (h minPath) Len() int {
	return len(h)
}

func (h minPath) Less(i, j int) bool {
	return h[i].value < h[j].value
}

func (h minPath) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *minPath) Push(x interface{}) {
	*h = append(*h, x.(path))
}

func (h *minPath) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type heap struct {
	values *minPath
}

func newHeap() *heap {
	return &heap{values: &minPath{}}
}

func (h *heap) push(p path) {
	hp.Push(h.values, p)
}

func (h *heap) pop() path {
	i := hp.Pop(h.values)
	return i.(path)
}

type edge struct {
	node   int
	weight float64
}

type graph struct {
	nodes map[int][]edge
}

func newGraph() *graph {
	return &graph{nodes: make(map[int][]edge)}
}

func (g *graph) addEdge(origin, destiny int, weight float64) {
	g.nodes[origin] = append(g.nodes[origin], edge{node: destiny, weight: weight})
	g.nodes[destiny] = append(g.nodes[destiny], edge{node: origin, weight: weight})
}

func (g *graph) getEdges(node int) []edge {
	return g.nodes[node]
}

func (g *graph) getPath(origin, destiny int) (float64, []int) {
	h := newHeap()
	h.push(path{value: 0, nodes: []int{origin}})
	visited := make(map[int]bool)

	for len(*h.values) > 0 {
		p := h.pop()
		node := p.nodes[len(p.nodes)-1]

		if visited[node] {
			continue
		}

		if node == destiny {
			return p.value, p.nodes
		}

		for _, e := range g.getEdges(node) {
			if !visited[e.node] {
				h.push(
					path{
						value: p.value + e.weight,
						nodes: append([]int{}, append(p.nodes, e.node)...),
					},
				)
			}
		}

		visited[node] = true
	}

	return 0, nil
}
