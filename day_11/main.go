//https://adventofcode.com/2025/day/11

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Node struct {
	id         string
	neighbours []*Node
}

func NewNode(id string) *Node {
	return &Node{
		id:         id,
		neighbours: []*Node{},
	}
}

type DAG struct {
	Nodes map[string]*Node
}

func (dag *DAG) GetNode(id string) *Node {
	return dag.Nodes[id]
}

func (dag *DAG) String() string {
	var sb strings.Builder
	for id, node := range dag.Nodes {
		sb.WriteString(fmt.Sprintf("Node %s: [", id))
		for _, neighbour := range node.neighbours {
			sb.WriteString(fmt.Sprintf("%s ", neighbour.id))
		}
		sb.WriteString("]\n")
	}
	return sb.String()
}

func (dag *DAG) AddNode(id string) {
	if dag.Nodes == nil {
		dag.Nodes = make(map[string]*Node)
	}
	_, exists := dag.Nodes[id]
	// only add if not exists
	if exists {
		return
	}
	dag.Nodes[id] = NewNode(id)
}

func (dag *DAG) AddEdge(fromID, toID string) {
	fromNode, fromExists := dag.Nodes[fromID]
	if !fromExists {
		dag.AddNode(fromID)
		fromNode = dag.Nodes[fromID]
	}
	toNode, toExists := dag.Nodes[toID]
	if !toExists {
		dag.AddNode(toID)
		toNode = dag.Nodes[toID]
	}

	fromNode.neighbours = append(fromNode.neighbours, toNode)
	dag.Nodes[fromID] = fromNode
}

func (dag *DAG) TopologicalSort() []*Node {
	inDegree := make(map[string]int)
	for id := range dag.Nodes {
		inDegree[id] = 0
	}
	for _, node := range dag.Nodes {
		for _, neighbour := range node.neighbours {
			inDegree[neighbour.id]++
		}
	}

	// var queue []*Node
	queue := make([]*Node, 0, len(dag.Nodes))
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, dag.Nodes[id])
		}
	}
	var sorted []*Node
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		sorted = append(sorted, current)

		for _, neighbour := range current.neighbours {
			inDegree[neighbour.id]--
			if inDegree[neighbour.id] == 0 {
				queue = append(queue, neighbour)
			}
		}
	}

	if len(sorted) != len(dag.Nodes) {
		panic("graph has at least one cycle, topological sort not possible")
	}

	return sorted
}

func (dag *DAG) FindAllPaths(startID, endID string) [][]string {

	// Implement the logic to find all paths from startID to endID using topoSorted
	// 1) Get Topological Sort Order
	topoSorted := dag.TopologicalSort()

	// 2) Reachability DP to dst in reverse topo order
	canReach := make(map[string]bool)
	canReach[endID] = true
	for i := len(topoSorted) - 1; i >= 0; i-- {
		node := topoSorted[i]
		if node.id == endID {
			continue
		}
		for _, neighbour := range node.neighbours {
			if canReach[neighbour.id] {
				canReach[node.id] = true
				break
			}
		}
	}

	// Quick exit: src can't reach dst at all.
	if !canReach[startID] {
		return nil
	}

	// 3) DFS enumerate with pruning
	var results [][]string
	var dfs func(currentID string, path []string)
	dfs = func(currentID string, path []string) {
		if currentID == endID {
			// found a path
			pathCopy := make([]string, len(path))
			copy(pathCopy, path)
			results = append(results, pathCopy)
			return
		}
		currentNode := dag.GetNode(currentID)
		for _, neighbour := range currentNode.neighbours {
			// Prune:  don't explore neighbours that can't reach dst
			if canReach[neighbour.id] {
				path = append(path, neighbour.id)
				dfs(neighbour.id, path)
				path = path[:len(path)-1] // backtrack
			}
		}
	}

	dfs(startID, []string{startID})
	return results
}

func ParseInput() DAG {
	dag := DAG{}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(input), "\n")

	// reach line is like: ```aaa: you hhh```
	// so we split by ": " to get the node_id from part[0]
	// then trim part[1] and then split by " " to get the neighbour ids
	for _, line := range lines {
		parts := strings.Split(line, ":")
		nodeID := strings.TrimSpace(parts[0])
		dag.AddNode(nodeID)

		neighbours := strings.Fields(strings.TrimSpace(parts[1]))
		for _, neighbourID := range neighbours {
			dag.AddEdge(nodeID, neighbourID)
		}
	}
	return dag
}

func Part1() {

	dag := ParseInput()
	fmt.Printf("DAG:\n%s\n", dag.String())

	allPaths := dag.FindAllPaths("you", "out")
	for _, path := range allPaths {
		fmt.Println(strings.Join(path, " -> "))
	}

	fmt.Printf("Total Paths from 'you' to 'out': %d\n", len(allPaths))
}

func Part2() {
	dag := ParseInput()
	// fmt.Printf("DAG:\n%s\n", dag.String())

	// topological sort to see the relationship between 'fft' and 'dac'
	sortedNodes := dag.TopologicalSort()
	precedingNode := ""
	for _, node := range sortedNodes {
		if node.id == "fft" || node.id == "dac" {
			precedingNode = node.id
			break
		}
	}
	fmt.Printf("In Topological Order, '%s' comes first.\n", precedingNode)

	// In Topological Order, 'fft' comes first.

	// so we find all the paths from 'svr' to 'fft'
	pathsToFFT := dag.FindAllPaths("svr", "fft")
	fmt.Println("pathsToFFT: ", len(pathsToFFT))

	// find all the paths from 'dac' to 'out'
	pathsDacOut := dag.FindAllPaths("dac", "out")
	fmt.Println("pathsDacOut: ", len(pathsDacOut))

	// find the paths from 'fft' to 'dac'
	pathsFFTDac := dag.FindAllPaths("fft", "dac")
	fmt.Println("pathsFFTDac: ", len(pathsFFTDac))

	// result shall be the combination of these paths.
	fmt.Println("combination of paths:", len(pathsToFFT)*len(pathsFFTDac)*len(pathsDacOut))

}

func main() {
	fmt.Println("--- Day 11: Reactor ---")
	// Part1()
	Part2()
}
