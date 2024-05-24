package main

import (
	"fmt"
	"main/generators"
	"main/graphs"
	"time"
)

const (
	CYCLE_MESSAGE = "Graf zawiera cykl. Sortowanie niemożliwe."
	INPUT_FILE    = "input"
)

var N_LIST = []int{1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000}

func testFunc(f func() ([]int16, bool)) int64 {
	now := time.Now()
	f()
	return time.Since(now).Nanoseconds()
}

func main() {
	// m := graphs.AdjacencyMatrix{}
	// m.BuildFromFile(INPUT_FILE)
	//
	// if s, cycle := m.TopologicalSort(); !cycle {
	// 	fmt.Println(s)
	// } else {
	// 	fmt.Println(CYCLE_MESSAGE)
	// }

	graphMatrixTime := make(map[int][]int64)
	adjacencyMatrixTime := make(map[int][]int64)

	for _, n := range N_LIST {
		input := generators.GenerateInput(n)

		m := graphs.AdjacencyMatrix{}
		m.BuildFromInput(input, n)

		adjacencyMatrixTime[n] = append(adjacencyMatrixTime[n], testFunc(m.TopologicalSort))
		adjacencyMatrixTime[n] = append(adjacencyMatrixTime[n], testFunc(m.TopologicalSortKahn))

		m2 := graphs.GraphMatrix{}
		m2.BuildFromInput(input, n)

		graphMatrixTime[n] = append(graphMatrixTime[n], testFunc(m2.TopologicalSort))
		graphMatrixTime[n] = append(graphMatrixTime[n], testFunc(m2.TopologicalSortKahn))
	}

	fmt.Println("Graph matrix time:")
	for _, n := range N_LIST {
		fmt.Printf("%d;%d;%d\n", n, graphMatrixTime[n][0], graphMatrixTime[n][1])
	}

	fmt.Println("Adjacency matrix time:")
	for _, n := range N_LIST {
		fmt.Printf("%d;%d;%d\n", n, adjacencyMatrixTime[n][0], adjacencyMatrixTime[n][1])
	}
}
