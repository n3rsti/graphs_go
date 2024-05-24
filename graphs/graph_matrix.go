package graphs

import (
	"bufio"
	"fmt"
	"main/utils"
	"os"
	"strconv"
	"strings"
)

type GraphMatrix struct {
	Graph [][]int16
}

func (m GraphMatrix) Print() {
	for _, row := range m.Graph {
		for _, cell := range row {
			fmt.Print(cell, " ")
		}
		fmt.Println()
	}
}

func (m GraphMatrix) get(i, j int16) int16 {
	return m.Graph[i-1][j-1]
}

func NewGraphMatrix(vertices int) *GraphMatrix {
	graph := make([][]int16, vertices)
	for i := range graph {
		graph[i] = make([]int16, vertices+3)
	}
	return &GraphMatrix{Graph: graph}
}

func (m *GraphMatrix) addEdge(v1, v2 int) {
	m.Graph[v1][v2] = 1
}

func (m *GraphMatrix) BuildFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	header := strings.Fields(scanner.Text())

	numVertices, err := strconv.Atoi(header[0])
	if err != nil {
		return err
	}

	edges := [][]int16{}
	for scanner.Scan() {
		edge := strings.Fields(scanner.Text())

		v1, _ := strconv.ParseInt(edge[0], 10, 16)
		v2, _ := strconv.ParseInt(edge[1], 10, 16)

		edges = append(edges, []int16{int16(v1), int16(v2)})
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	graph := NewGraphMatrix(numVertices)

	successors := make([][]int16, numVertices)
	predecessors := make([][]int16, numVertices)
	nonIncident := make([][]int16, numVertices)

	for _, edge := range edges {
		v1, v2 := edge[0]-1, edge[1]-1
		successors[v1] = append(successors[v1], v2+1)
		predecessors[v2] = append(predecessors[v2], v1+1)
	}

	for i := 0; i < numVertices; i++ {
		incidentVertices := make(map[int16]bool)

		for _, edge := range edges {
			if edge[0]-1 == int16(i) {
				incidentVertices[edge[1]] = true
			}
			if edge[1]-1 == int16(i) {
				incidentVertices[edge[0]] = true
			}
		}
		for j := int16(1); j <= int16(numVertices); j++ {
			if !incidentVertices[int16(j)] {
				nonIncident[i] = append(nonIncident[i], j)
			}
		}
	}

	graph.populate(successors, predecessors, nonIncident)

	m.Graph = graph.Graph
	return nil
}

func (m *GraphMatrix) BuildFromInput(input [][2]int16, vertices int) {
	graph := NewGraphMatrix(vertices)

	successors := make([][]int16, vertices)
	predecessors := make([][]int16, vertices)
	nonIncident := make([][]int16, vertices)

	for _, edge := range input {
		v1, v2 := edge[0]-1, edge[1]-1
		successors[v1] = append(successors[v1], v2+1)
		predecessors[v2] = append(predecessors[v2], v1+1)
	}

	for i := 0; i < vertices; i++ {
		incidentVertices := make(map[int16]bool)

		for _, edge := range input {
			if edge[0]-1 == int16(i) {
				incidentVertices[edge[1]] = true
			}
			if edge[1]-1 == int16(i) {
				incidentVertices[edge[0]] = true
			}
		}
		for j := 1; j <= vertices; j++ {
			if !incidentVertices[int16(j)] {
				nonIncident[i] = append(nonIncident[i], int16(j))
			}
		}
	}

	graph.populate(successors, predecessors, nonIncident)

	m.Graph = graph.Graph
}

func (m *GraphMatrix) populate(successors, predecessors, nonIncident [][]int16) {
	numVertices := int16(len(m.Graph))

	for i := int16(0); i < numVertices; i++ {
		if len(successors[i]) > 0 {
			m.Graph[i][numVertices] = successors[i][0]
			for j := 0; j < len(successors[i]); j++ {
				next := successors[i][j]
				if j == len(successors[i])-1 {
					m.Graph[i][next-1] = next
				} else {
					m.Graph[i][next-1] = successors[i][j+1]
				}
			}
		} else {
			m.Graph[i][numVertices] = 0
		}
	}

	for i := int16(0); i < numVertices; i++ {
		if len(predecessors[i]) > 0 {
			m.Graph[i][numVertices+1] = predecessors[i][0] + numVertices
			for j := 0; j < len(predecessors[i]); j++ {
				next := predecessors[i][j]
				if j == len(predecessors[i])-1 {
					m.Graph[i][next-1] = next + numVertices
				} else {
					m.Graph[i][next-1] = predecessors[i][j+1] + numVertices
				}
			}
		} else {
			m.Graph[i][numVertices+1] = 0
		}
	}

	for i := int16(0); i < numVertices; i++ {
		if len(nonIncident[i]) > 0 {
			m.Graph[i][numVertices+2] = -nonIncident[i][0]
			for j := 0; j < len(nonIncident[i]); j++ {
				next := nonIncident[i][j]
				if j == len(nonIncident[i])-1 {
					m.Graph[i][next-1] = -next
				} else {
					m.Graph[i][next-1] = -nonIncident[i][j+1]
				}
			}
		} else {
			m.Graph[i][numVertices+2] = 0
		}
	}
}

func (m GraphMatrix) TopologicalSort() ([]int16, bool) {
	vertices := len(m.Graph[0]) - 3
	visited := make([]bool, vertices)
	progressStack := make([]bool, vertices)
	stack := make([]int16, 0, vertices)
	cycle := false

	var dfs func(v int16) bool
	dfs = func(v int16) bool {
		if progressStack[v] {
			cycle = true
			return true
		}
		if visited[v] {
			return false
		}

		visited[v] = true
		progressStack[v] = true

		predecessor := m.Graph[v][vertices]
		for predecessor != 0 {
			if !visited[predecessor-1] {
				if dfs(predecessor - 1) {
					return true
				}
			}
			nextPredecessor := m.Graph[v][predecessor-1]
			if predecessor == nextPredecessor {
				break
			}
			predecessor = nextPredecessor
		}

		progressStack[v] = false
		stack = append(stack, v+1)
		return false
	}

	for i := int16(0); i < int16(vertices); i++ {
		if !visited[i] {
			if dfs(i) {
				return nil, true
			}
		}
	}

	utils.Reverse(stack)

	return stack, cycle
}

func (m GraphMatrix) TopologicalSortKahn() ([]int16, bool) {
	vertices := len(m.Graph[0]) - 3
	inDegree := make([]int, vertices)
	queue := make([]int16, 0, vertices)
	result := make([]int16, 0, vertices)

	for i := int16(1); i <= int16(vertices); i++ {
		predecessor := m.get(i, int16(vertices+2)) - int16(vertices)
		for predecessor != 0 && predecessor != int16(-vertices) {
			inDegree[i-1]++
			nextPredecessor := m.get(i, predecessor) - int16(vertices)
			if predecessor == nextPredecessor {
				break
			}
			predecessor = nextPredecessor
		}
	}

	for i := int16(1); i <= int16(vertices); i++ {
		if inDegree[i-1] == 0 {
			queue = append(queue, i)
		}
	}

	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		result = append(result, vertex)

		successor := m.get(vertex, int16(vertices+1))
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

	if int16(len(result)) != int16(vertices) {
		return nil, true
	}

	return result, false
}
