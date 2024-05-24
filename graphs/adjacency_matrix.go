package graphs

import (
	"bufio"
	"fmt"
	"main/utils"
	"os"
	"strconv"
	"strings"
)

type AdjacencyMatrix struct {
	Graph [][]int16
}

func (m AdjacencyMatrix) Print() {
	for _, row := range m.Graph {
		for _, cell := range row {
			fmt.Print(cell, " ")
		}
		fmt.Println()
	}
}

func NewAdjacencyMatrix(vertices int) *AdjacencyMatrix {
	graph := make([][]int16, vertices)
	for i := range graph {
		graph[i] = make([]int16, vertices)
	}
	return &AdjacencyMatrix{Graph: graph}
}

func (m *AdjacencyMatrix) addEdge(v1, v2 int16) {
	m.Graph[v1][v2] = 1
}

func (m *AdjacencyMatrix) BuildFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	firstLine := scanner.Text()
	parts := strings.Fields(firstLine)

	numVertices, _ := strconv.Atoi(parts[0])
	graph := NewAdjacencyMatrix(numVertices)

	for scanner.Scan() {
		line := scanner.Text()
		edgeParts := strings.Fields(line)

		v1, _ := strconv.ParseInt(edgeParts[0], 10, 16)
		v2, _ := strconv.ParseInt(edgeParts[1], 10, 16)
		graph.addEdge(int16(v1-1), int16(v2-1))
	}

	m.Graph = graph.Graph

	return nil
}

func (m *AdjacencyMatrix) BuildFromInput(input [][2]int16, vertices int) {
	graph := NewAdjacencyMatrix(vertices)

	for _, edge := range input {
		graph.addEdge(edge[0]-1, edge[1]-1)
	}

	m.Graph = graph.Graph
}

func (m AdjacencyMatrix) TopologicalSort() ([]int16, bool) {
	vertices := int16(len(m.Graph))
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

		for i := int16(0); i < vertices; i++ {
			if m.Graph[v][i] == 1 {
				if dfs(i) {
					return true
				}
			}
		}

		progressStack[v] = false
		stack = append(stack, v+1)
		return false
	}

	for i := int16(0); i < vertices; i++ {
		if !visited[i] {
			if dfs(i) {
				return nil, true
			}
		}
	}

	utils.Reverse(stack)

	return stack, cycle
}

func (m AdjacencyMatrix) TopologicalSortKahn() ([]int16, bool) {
	vertices := int16(len(m.Graph))
	inDegree := make([]int, vertices)
	queue := make([]int16, 0, vertices)
	result := make([]int16, 0, vertices)

	for i := int16(0); i < vertices; i++ {
		for j := int16(0); j < vertices; j++ {
			if m.Graph[j][i] == 1 {
				inDegree[i]++
			}
		}
	}

	for i := int16(0); i < vertices; i++ {
		if inDegree[i] == 0 {
			queue = append(queue, i)
		}
	}

	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		result = append(result, vertex+1)

		for i := int16(0); i < vertices; i++ {
			if m.Graph[vertex][i] == 1 {
				inDegree[i]--
				if inDegree[i] == 0 {
					queue = append(queue, i)
				}
			}
		}
	}

	if int16(len(result)) != vertices {
		return nil, true
	}

	return result, false
}
