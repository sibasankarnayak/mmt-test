package data_structure

import (
	"fmt"
	"math"
	"strings"
)

func (graph *Graph) Dijkstra(source ID) (dist map[ID]float64, prev map[ID]ID, err error) {
	if _, exists := graph.vertices[source]; !exists {
		return nil, nil, fmt.Errorf("Vertex %v is not existed", source)
	}

	dist = make(map[ID]float64)
	prev = make(map[ID]ID)
	heap := NewFibHeap()

	for id := range graph.vertices {
		prev[id] = nil
		if id != source {
			dist[id] = math.Inf(1)
			heap.Insert(id, math.Inf(1))
		} else {
			dist[id] = 0
			heap.Insert(id, 0)
		}
	}

	for heap.Num() != 0 {
		min, _ := heap.ExtractMin()
		for to, edge := range graph.egress[min] {
			if edge.getWeight() < 0 {
				return nil, nil, fmt.Errorf("Negative weight form vertex %v to vertex %v is not allowed", min, to)
			}
			if !edge.enable {
				continue
			}

			if dist[min]+edge.getWeight() < dist[to] {
				heap.DecreaseKey(to, dist[min]+edge.getWeight())
				prev[to] = min.(string) + "_" + edge.getCode().(string)
				// prev[to] = min
				dist[to] = dist[min] + edge.getWeight()
			}
		}
	}

	return
}

func getPath(prev map[ID]ID, lastNode ID) (path []ID) {
	prevNode := prev[lastNode]
	if prevNode == nil {
		return nil
	}
	reversePath := []ID{lastNode}
	for ; prevNode != nil; prevNode = prev[strings.Split(prevNode.(string), "_")[0]] {
		reversePath = append(reversePath, prevNode)

	}

	path = make([]ID, len(reversePath))
	for index, node := range reversePath {
		path[len(reversePath)-index-1] = node
	}

	return
}
