package grapher

import (
	"container/heap"
)

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

func Dijkstra(graph Graph, start, end string) ([]string, int) {
	dist := make(map[string]int)
	prev := make(map[string]string)
	for _, node := range graph.Nodes {
		if node == start {
			dist[node] = 0
		} else {
			dist[node] = 1<<31 - 1 // infinity
		}
	}

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{
		Node:     start,
		Priority: dist[start],
	})

	for pq.Len() > 0 {
		u := heap.Pop(&pq).(*Item)
		if u.Node == end {
			break
		}

		for _, edge := range graph.Edges {
			var neighbor string
			if edge.From == u.Node {
				neighbor = edge.To
			} else if edge.To == u.Node {
				neighbor = edge.From
			} else {
				continue
			}

			alt := dist[u.Node] + edge.Weight
			if alt < dist[neighbor] {
				dist[neighbor] = alt
				prev[neighbor] = u.Node

				heap.Push(&pq, &Item{
					Node:     neighbor,
					Priority: alt,
				})
			}
		}
	}

	var path []string
	if prev[end] != "" || start == end {
		for u := end; u != ""; u = prev[u] {
			path = append([]string{u}, path...)
		}
	}

	return path, dist[end]
}