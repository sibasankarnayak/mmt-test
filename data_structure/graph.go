package data_structure

import (
	"fmt"
	"math"
)

type ID interface{}

type Vertex interface {
	ID() ID
	Edges() []Edge
}

type Edge interface {
	Get() (from ID, to ID, weight float64, code ID)
}

type Graph struct {
	vertices map[ID]*vertex
	egress   map[ID]map[ID]*edge
	ingress  map[ID]map[ID]*edge
}

type vertex struct {
	self   interface{}
	enable bool
}

type edge struct {
	self    interface{}
	weight  float64
	enable  bool
	changed bool
	code    interface{}
}

func (edge *edge) getWeight() float64 {
	return edge.weight
}

func (edge *edge) getCode() interface{} {
	return edge.code
}

func NewGraph() *Graph {
	graph := new(Graph)
	graph.vertices = make(map[ID]*vertex)
	graph.egress = make(map[ID]map[ID]*edge)
	graph.ingress = make(map[ID]map[ID]*edge)

	return graph
}

func (graph *Graph) GetVertex(id ID) (vertex interface{}, err error) {
	if v, exists := graph.vertices[id]; exists {
		vertex = v.self
		return
	}

	err = fmt.Errorf("Vertex %v is not found", id)
	return
}

func (graph *Graph) GetEdge(from ID, to ID) (interface{}, error) {
	if _, exists := graph.vertices[from]; !exists {
		return nil, fmt.Errorf("Vertex(from) %v is not found", from)
	}

	if _, exists := graph.vertices[to]; !exists {
		return nil, fmt.Errorf("Vertex(to) %v is not found", to)
	}

	if edge, exists := graph.egress[from][to]; exists {
		return edge.self, nil
	}

	return nil, fmt.Errorf("Edge from %v to %v is not found", from, to)
}

func (graph *Graph) GetEdgeWeight(from ID, to ID) (float64, error) {
	if _, exists := graph.vertices[from]; !exists {
		return math.Inf(1), fmt.Errorf("Vertex(from) %v is not found", from)
	}

	if _, exists := graph.vertices[to]; !exists {
		return math.Inf(1), fmt.Errorf("Vertex(to) %v is not found", to)
	}

	if edge, exists := graph.egress[from][to]; exists {
		return edge.weight, nil
	}

	return math.Inf(1), nil
}

func (graph *Graph) AddVertex(id ID, v interface{}) error {
	if _, exists := graph.vertices[id]; exists {
		return fmt.Errorf("Vertex %v is duplicate", id)
	}

	graph.vertices[id] = &vertex{v, true}
	graph.egress[id] = make(map[ID]*edge)
	graph.ingress[id] = make(map[ID]*edge)

	return nil
}

func (graph *Graph) AddEdge(from ID, to ID, code ID, weight float64, e interface{}) error {
	if weight == math.Inf(-1) {
		return fmt.Errorf("-inf weight is reserved for internal usage")
	}

	if _, exists := graph.vertices[from]; !exists {
		return fmt.Errorf("Vertex(from) %v is not found", from)
	}

	if _, exists := graph.vertices[to]; !exists {
		return fmt.Errorf("Vertex(to) %v is not found", to)
	}

	if _, exists := graph.egress[from][to]; exists {
		return fmt.Errorf("Edge from %v to %v is duplicate", from, to)
	}

	graph.egress[from][to] = &edge{e, weight, true, false, code}
	graph.ingress[to][from] = graph.egress[from][to]

	return nil
}

func (graph *Graph) UpdateEdgeWeight(from ID, to ID, weight float64) error {
	if weight == math.Inf(-1) {
		return fmt.Errorf("-inf weight is reserved for internal usage")
	}

	if _, exists := graph.vertices[from]; !exists {
		return fmt.Errorf("Vertex(from) %v is not found", from)
	}

	if _, exists := graph.vertices[to]; !exists {
		return fmt.Errorf("Vertex(to) %v is not found", to)
	}

	if edge, exists := graph.egress[from][to]; exists {
		edge.weight = weight
		return nil
	}

	return fmt.Errorf("Edge from %v to %v is not found", from, to)
}

func (graph *Graph) DeleteVertex(id ID) interface{} {
	if vertex, exists := graph.vertices[id]; exists {
		for to := range graph.egress[id] {
			delete(graph.ingress[to], id)
		}
		for from := range graph.ingress[id] {
			delete(graph.egress[from], id)
		}
		delete(graph.egress, id)
		delete(graph.ingress, id)
		delete(graph.vertices, id)

		return vertex.self
	}

	return nil
}

func (graph *Graph) DeleteEdge(from ID, to ID) interface{} {
	if _, exists := graph.vertices[from]; !exists {
		return nil
	}

	if _, exists := graph.vertices[to]; !exists {
		return nil
	}

	if edge, exists := graph.egress[from][to]; exists {
		delete(graph.egress[from], to)
		delete(graph.ingress[to], from)
		return edge.self
	}

	return nil
}

func (graph *Graph) AddVertexWithEdges(v Vertex) error {
	if _, exists := graph.vertices[v.ID()]; exists {
		return fmt.Errorf("Vertex %v is duplicate", v.ID())
	}

	graph.vertices[v.ID()] = &vertex{v, true}
	graph.egress[v.ID()] = make(map[ID]*edge)
	graph.ingress[v.ID()] = make(map[ID]*edge)

	for _, eachEdge := range v.Edges() {
		from, to, weight, code := eachEdge.Get()
		if weight == math.Inf(-1) {
			return fmt.Errorf("-inf weight is reserved for internal usage")
		}
		if from != v.ID() && to != v.ID() {
			return fmt.Errorf("Edge from %v to %v is unrelated to the vertex %v", from, to, v.ID())
		}

		if _, exists := graph.egress[to]; !exists {
			graph.egress[to] = make(map[ID]*edge)
		}
		if _, exists := graph.egress[from]; !exists {
			graph.egress[from] = make(map[ID]*edge)
		}
		if _, exists := graph.ingress[from]; !exists {
			graph.ingress[from] = make(map[ID]*edge)
		}
		if _, exists := graph.ingress[to]; !exists {
			graph.ingress[to] = make(map[ID]*edge)
		}

		graph.egress[from][to] = &edge{eachEdge, weight, true, false, code}
		graph.ingress[to][from] = graph.egress[from][to]
	}

	return nil
}

func (graph *Graph) CheckIntegrity() error {
	for from, out := range graph.egress {
		if _, exists := graph.vertices[from]; !exists {
			return fmt.Errorf("Vertex %v is not found", from)
		}
		for to := range out {
			if _, exists := graph.vertices[to]; !exists {
				return fmt.Errorf("Vertex %v is not found", to)
			}
		}
	}

	for to, in := range graph.ingress {
		if _, exists := graph.vertices[to]; !exists {
			return fmt.Errorf("Vertex %v is not found", to)
		}
		for from := range in {
			if _, exists := graph.vertices[from]; !exists {
				return fmt.Errorf("Vertex %v is not found", from)
			}
		}
	}

	return nil
}

func (graph *Graph) GetPathWeight(path []ID) (totalWeight float64) {
	if len(path) == 0 {
		return math.Inf(-1)
	}

	if _, exists := graph.vertices[path[0]]; !exists {
		return math.Inf(-1)
	}

	for i := 0; i < len(path)-1; i++ {
		if _, exists := graph.vertices[path[i+1]]; !exists {
			return math.Inf(-1)
		}
		if edge, exists := graph.egress[path[i]][path[i+1]]; exists {
			totalWeight += edge.getWeight()
		} else {
			return math.Inf(1)
		}
	}

	return totalWeight
}

func (graph *Graph) DisableEdge(from, to ID) {
	graph.egress[from][to].enable = false
}

func (graph *Graph) DisableVertex(vertex ID) {
	for _, edge := range graph.egress[vertex] {
		edge.enable = false
	}
}

func (graph *Graph) DisablePath(path []ID) {
	for _, vertex := range path {
		graph.DisableVertex(vertex)
	}
}

func (graph *Graph) Reset() {
	for _, out := range graph.egress {
		for _, edge := range out {
			edge.enable = true
		}
	}
}
