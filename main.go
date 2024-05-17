package main

import (
	"fmt"
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
	graph := AdjacencyMatrix{
		graph: exampleAdjecencyMatrix,
	}

	if sorted, cycle := graph.topologicalSortKahn(); !cycle {
		fmt.Println(sorted)
	} else {
		fmt.Println(CYCLE_MESSAGE)
	}
}
