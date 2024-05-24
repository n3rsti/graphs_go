package generators

import (
	"fmt"
	"math/rand"
)

func GenerateInput(numVertices int) [][2]int16 {
	numEdges := numVertices * (numVertices - 1) / 2
	vertices := rand.Perm(numVertices)
	edges := make(map[string]bool)
	result := make([][2]int16, 0, numEdges)

	edgeCount := 0
	for edgeCount < numEdges {
		for i := 0; i < numVertices; i++ {
			for j := i + 1; j < numVertices; j++ {
				if edgeCount >= numEdges {
					break
				}

				v1 := int16(vertices[i]) + 1
				v2 := int16(vertices[j]) + 1

				edge := fmt.Sprintf("%d %d", v1, v2)
				if !edges[edge] {
					edges[edge] = true
					result = append(result, [2]int16{v1, v2})

					edgeCount++
				}
			}
		}
	}

	return result
}
