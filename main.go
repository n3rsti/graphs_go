package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	exampleGraphMatrix = [][]int{
		{-1, 8, 8, 6, -1, 6, 4, 9, -5},
		{3, -5, 3, -6, -4, -6, 1, 0, -2},
		{1, 8, -5, 1, -6, -6, 4, 8, -3},
		{9, -4, 12, -4, 5, 12, 5, 7, -2},
		{-2, -3, -5, 10, -5, 10, 0, 12, -1},
		{7, -2, -2, 5, 5, -3, 4, 7, -6},
	}
	exampleAdjecencyMatrix = [][]int{
		{0, 1, 0, 0, 0},
		{0, 0, 1, 1, 0},
		{0, 0, 0, 0, 1},
		{0, 0, 0, 0, 1},
		{0, 0, 0, 0, 0},
	}
)

const CYCLE_MESSAGE = "Graf zawiera cykl. Sortowanie niemożliwe."

func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

type AdjacencyMatrix struct {
	graph [][]int
}

func NewAdjacencyMatrix(vertices int) *AdjacencyMatrix {
	graph := make([][]int, vertices)
	for i := range graph {
		graph[i] = make([]int, vertices)
	}
	return &AdjacencyMatrix{graph: graph}
}

func NewGraphMatrix(vertices int) *GraphMatrix {
	graph := make([][]int, vertices)
	for i := range graph {
		graph[i] = make([]int, vertices+3)
	}
	return &GraphMatrix{graph: graph}
}

func (m *AdjacencyMatrix) addEdge(v1, v2 int) {
	m.graph[v1][v2] = 1
}

func (m *GraphMatrix) addEdge(v1, v2 int) {
	m.graph[v1][v2] = 1
}
func buildAdjacencyMatrixFromFile(filename string) (*AdjacencyMatrix, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return nil, fmt.Errorf("eof")
	}
	firstLine := scanner.Text()
	parts := strings.Fields(firstLine)

	numVertices, err := strconv.Atoi(parts[0])
	graph := NewAdjacencyMatrix(numVertices)

	numEdges, err := strconv.Atoi(parts[1])

	for i := 0; i < numEdges; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("eof")
		}
		line := scanner.Text()
		edgeParts := strings.Fields(line)

		v1, _ := strconv.Atoi(edgeParts[0])
		v2, _ := strconv.Atoi(edgeParts[1])
		graph.addEdge(v1-1, v2-1)
	}

	return graph, nil
}

func (m *GraphMatrix) populateMatrix(successors, predecessors, nonIncident [][]int) {
	numVertices := len(m.graph)

	for i := 0; i < numVertices; i++ {
		if len(successors[i]) > 0 {
			m.graph[i][numVertices] = successors[i][0]
			for j := 0; j < len(successors[i]); j++ {
				next := successors[i][j]
				if j == len(successors[i])-1 {
					m.graph[i][next-1] = next
				} else {
					m.graph[i][next-1] = successors[i][j+1]
				}
			}
		} else {
			m.graph[i][numVertices] = 0
		}
	}

	for i := 0; i < numVertices; i++ {
		if len(predecessors[i]) > 0 {
			m.graph[i][numVertices+1] = predecessors[i][0]
			for j := 0; j < len(predecessors[i]); j++ {
				next := predecessors[i][j]
				if j == len(predecessors[i])-1 {
					m.graph[i][next-1] = next + numVertices
				} else {
					m.graph[i][next-1] = predecessors[i][j+1] + numVertices
				}
			}
		} else {
			m.graph[i][numVertices+1] = 0
		}
	}

	for i := 0; i < numVertices; i++ {
		if len(nonIncident[i]) > 0 {
			m.graph[i][numVertices+2] = nonIncident[i][0]
			for j := 0; j < len(nonIncident[i]); j++ {
				next := nonIncident[i][j]
				if j == len(nonIncident[i])-1 {
					m.graph[i][next-1] = -next
				} else {
					m.graph[i][next-1] = -nonIncident[i][j+1]
				}
			}
		} else {
			m.graph[i][numVertices+2] = 0
		}
	}
}

func readGraphFromFile(filename string) (*GraphMatrix, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	header := strings.Fields(scanner.Text())

	numVertices, err := strconv.Atoi(header[0])
	if err != nil {
		return nil, err
	}

	edges := [][]int{}
	for scanner.Scan() {
		edge := strings.Fields(scanner.Text())

		v1, _ := strconv.Atoi(edge[0])

		v2, _ := strconv.Atoi(edge[1])

		edges = append(edges, []int{v1, v2})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	graph := NewGraphMatrix(numVertices)

	successors := make([][]int, numVertices)
	predecessors := make([][]int, numVertices)
	nonIncident := make([][]int, numVertices)

	for _, edge := range edges {
		v1, v2 := edge[0]-1, edge[1]-1
		successors[v1] = append(successors[v1], v2+1)
		predecessors[v2] = append(predecessors[v2], v1+1)
	}

	graph.populateMatrix(successors, predecessors, nonIncident)
	return graph, nil
}

type GraphMatrix struct {
	graph [][]int
}

func (m GraphMatrix) get(i, j int) int {
	return m.graph[i-1][j-1]
}

func (m AdjacencyMatrix) topologicalSort() ([]int, bool) {
	vertices := len(m.graph)
	visited := make([]bool, vertices)
	progressStack := make([]bool, vertices)
	stack := make([]int, 0, vertices)
	cycle := false

	var dfs func(v int) bool
	dfs = func(v int) bool {
		if progressStack[v] {
			cycle = true
			return true
		}
		if visited[v] {
			return false
		}

		visited[v] = true
		progressStack[v] = true

		for i := 0; i < vertices; i++ {
			if m.graph[v][i] == 1 {
				if dfs(i) {
					return true
				}
			}
		}

		progressStack[v] = false
		stack = append(stack, v)
		return false
	}

	for i := 0; i < vertices; i++ {
		if !visited[i] {
			if dfs(i) {
				return nil, true
			}
		}
	}

	reverse(stack)

	return stack, cycle
}

func (m AdjacencyMatrix) topologicalSortKahn() ([]int, bool) {
	vertices := len(m.graph)
	inDegree := make([]int, vertices)
	queue := make([]int, 0, vertices)
	result := make([]int, 0, vertices)

	for i := 0; i < vertices; i++ {
		for j := 0; j < vertices; j++ {
			if m.graph[j][i] == 1 {
				inDegree[i]++
			}
		}
	}

	for i := 0; i < vertices; i++ {
		if inDegree[i] == 0 {
			queue = append(queue, i)
		}
	}

	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		result = append(result, vertex)

		for i := 0; i < vertices; i++ {
			if m.graph[vertex][i] == 1 {
				inDegree[i]--
				if inDegree[i] == 0 {
					queue = append(queue, i)
				}
			}
		}
	}

	if len(result) != vertices {
		return nil, true
	}

	return result, false
}

func (m GraphMatrix) topologicalSort() ([]int, bool) {
	vertices := len(m.graph[0]) - 3
	visited := make([]bool, vertices)
	progressStack := make([]bool, vertices)
	stack := make([]int, 0, vertices)
	cycle := false

	var dfs func(v int) bool
	dfs = func(v int) bool {
		if progressStack[v] {
			cycle = true
			return true
		}
		if visited[v] {
			return false
		}

		visited[v] = true
		progressStack[v] = true

		predecessor := m.graph[v][vertices]
		for predecessor != 0 {
			if !visited[predecessor-1] {
				if dfs(predecessor - 1) {
					return true
				}
			}
			nextPredecessor := m.graph[v][predecessor-1]
			if predecessor == nextPredecessor {
				break
			}
			predecessor = nextPredecessor
		}

		progressStack[v] = false
		stack = append(stack, v+1)
		return false
	}

	for i := 0; i < vertices; i++ {
		if !visited[i] {
			if dfs(i) {
				return nil, true
			}
		}
	}

	reverse(stack)

	return stack, cycle
}

func (m GraphMatrix) topologicalSortKahn() ([]int, bool) {
	vertices := len(m.graph[0]) - 3
	inDegree := make([]int, vertices)
	queue := make([]int, 0, vertices)
	result := make([]int, 0, vertices)

	for i := 1; i <= vertices; i++ {
		predecessor := m.get(i, vertices+2) - vertices
		for predecessor != 0 && predecessor != -vertices {
			inDegree[i-1]++
			nextPredecessor := m.get(i, predecessor) - vertices
			if predecessor == nextPredecessor {
				break
			}
			predecessor = nextPredecessor
		}
	}

	for i := 1; i <= vertices; i++ {
		if inDegree[i-1] == 0 {
			queue = append(queue, i)
		}
	}

	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		result = append(result, vertex)

		successor := m.get(vertex, vertices+1)
		for successor != 0 {
			inDegree[successor-1]--
			if inDegree[successor-1] == 0 {
				queue = append(queue, successor)
			}

			nextSuccessor := m.get(vertex, successor)
			if successor == nextSuccessor {
				break
			}
			successor = nextSuccessor
		}
	}

	if len(result) != vertices {
		return nil, true
	}

	return result, false
}

func main() {
	// graph := AdjacencyMatrix{
	// 	graph: exampleAdjecencyMatrix,
	// }
	//
	// if sorted, cycle := graph.topologicalSortKahn(); !cycle {
	// 	fmt.Println(sorted)
	// } else {
	// 	fmt.Println(CYCLE_MESSAGE)
	// }

	m, _ := buildAdjacencyMatrixFromFile("graph.txt")
	fmt.Println(m.graph)

}
