package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	commons "github.com/shaohong/aoc2025/commons"
)

type Node struct {
	x  int
	y  int
	z  int
	id int
}

type NodeDistance struct {
	distance int
	nodePair [2]int // save two node ids, first is smaller id
}

type Circuit struct {
	nodes map[int]bool // set of node ids
}

type Circuits []Circuit

func Distance(a, b Node) int {
	dx := a.x - b.x
	dy := a.y - b.y
	dz := a.z - b.z
	return dx*dx + dy*dy + dz*dz
}

func ParseInput() []Node {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(input), "\n")
	nodes := make([]Node, 0)
	for n, line := range lines {
		if line == "" {
			continue
		}
		newNode := Node{id: n}
		fmt.Sscanf(line, "%d,%d,%d", &newNode.x, &newNode.y, &newNode.z)
		nodes = append(nodes, newNode)
	}
	return nodes
}

func Part1(connectionsToCheck int, topNCircuit int) {
	nodes := ParseInput()

	nodeDistances := make([]NodeDistance, 0)

	for i, nodeA := range nodes {
		for j := i + 1; j < len(nodes); j++ {
			nodeB := nodes[j]
			dist := Distance(nodeA, nodeB)
			nodeDistances = append(nodeDistances, NodeDistance{
				distance: dist,
				nodePair: [2]int{nodeA.id, nodeB.id},
			})
		}
	}

	// sort the nodeDistances by distance
	sort.Slice(nodeDistances, func(i, j int) bool {
		return nodeDistances[i].distance < nodeDistances[j].distance
	})

	// put these node distances into a queue
	nodeDistanceQueue := commons.Queue[NodeDistance]{}
	for _, nd := range nodeDistances {
		nodeDistanceQueue.Enqueue(nd)
	}

	circuits := makingConnections(nodeDistanceQueue, connectionsToCheck)

	fmt.Println("Total circuits formed:", len(circuits))

	totalProducts := 1
	for i := 0; i < min(topNCircuit, len(circuits)); i++ {
		fmt.Printf("Circuit %d size: %d, content: %+v\n", i, len(circuits[i].nodes), circuits[i].nodes)
		totalProducts *= len(circuits[i].nodes)
	}
	fmt.Printf("Total product of largest %d circuits: %d\n", topNCircuit, totalProducts)
}

func makingConnections(nodeDistanceQueue commons.Queue[NodeDistance], connectionsToCheck int) Circuits {
	circuits := Circuits{}
	totalConnectionsMade := 0

	for i := 0; i < connectionsToCheck; i++ {
		smalledDistanceNodePair, ok := nodeDistanceQueue.Dequeue()
		if !ok {
			break
		}
		nodeAID := smalledDistanceNodePair.nodePair[0]
		nodeBID := smalledDistanceNodePair.nodePair[1]

		var circuitAIndex, circuitBIndex = -1, -1
		for ci, circuit := range circuits {
			if circuit.nodes[nodeAID] {
				circuitAIndex = ci
			}
			if circuit.nodes[nodeBID] {
				circuitBIndex = ci
			}
		}

		if circuitAIndex == -1 && circuitBIndex == -1 {
			// create a new circuit
			newCircuit := Circuit{nodes: make(map[int]bool)}
			newCircuit.nodes[nodeAID] = true
			newCircuit.nodes[nodeBID] = true
			circuits = append(circuits, newCircuit)
			totalConnectionsMade++
			fmt.Printf("added connection between %d and node %d\n", nodeAID, nodeBID)
			continue
		} else if circuitAIndex != -1 && circuitBIndex != -1 {
			if circuitAIndex != circuitBIndex { // merge two circuits
				circuitA := circuits[circuitAIndex]
				circuitB := circuits[circuitBIndex]
				for nodeID := range circuitB.nodes {
					circuitA.nodes[nodeID] = true
				}

				circuits = append(circuits[:circuitBIndex], circuits[circuitBIndex+1:]...)
				totalConnectionsMade++
				fmt.Printf("added connection between %d and node %d\n", nodeAID, nodeBID)
				continue
			}
		} else { // add node to existing circuit
			circuitIdx := max(circuitAIndex, circuitBIndex)
			circuit := circuits[circuitIdx]
			circuit.nodes[nodeAID] = true
			circuit.nodes[nodeBID] = true
			totalConnectionsMade++
			fmt.Printf("added connection between %d and node %d\n", nodeAID, nodeBID)
			continue
		}

	}

	// sort the circuits by size, in descending order
	sort.Slice(circuits, func(i, j int) bool {
		return len(circuits[i].nodes) > len(circuits[j].nodes)
	})

	return circuits
}

func allInOneCircuit(circuits Circuits, totalNodes int) bool {
	if len(circuits) != 1 {
		return false
	}
	for _, circuit := range circuits {
		if len(circuit.nodes) != totalNodes {
			return false
		}
	}
	return true
}

// keep making connections until all nodes are connected in one circuit, return the last distances pair used
func makingConnectionsUntilAllConnected(nodeDistanceQueue commons.Queue[NodeDistance], totalNodes int) NodeDistance {
	circuits := Circuits{}
	totalConnectionsMade := 0

	var lastDistanceNodePair NodeDistance
	var ok bool
	for !allInOneCircuit(circuits, totalNodes) {
		lastDistanceNodePair, ok = nodeDistanceQueue.Dequeue()
		if !ok {
			break
		}
		nodeAID := lastDistanceNodePair.nodePair[0]
		nodeBID := lastDistanceNodePair.nodePair[1]

		var circuitAIndex, circuitBIndex = -1, -1
		for ci, circuit := range circuits {
			if circuit.nodes[nodeAID] {
				circuitAIndex = ci
			}
			if circuit.nodes[nodeBID] {
				circuitBIndex = ci
			}
		}

		if circuitAIndex == -1 && circuitBIndex == -1 {
			// create a new circuit
			newCircuit := Circuit{nodes: make(map[int]bool)}
			newCircuit.nodes[nodeAID] = true
			newCircuit.nodes[nodeBID] = true
			circuits = append(circuits, newCircuit)
			totalConnectionsMade++
			fmt.Printf("added connection between %d and node %d\n", nodeAID, nodeBID)
			continue
		} else if circuitAIndex != -1 && circuitBIndex != -1 {
			if circuitAIndex != circuitBIndex { // merge two circuits
				circuitA := circuits[circuitAIndex]
				circuitB := circuits[circuitBIndex]
				for nodeID := range circuitB.nodes {
					circuitA.nodes[nodeID] = true
				}

				circuits = append(circuits[:circuitBIndex], circuits[circuitBIndex+1:]...)
				totalConnectionsMade++
				fmt.Printf("added connection between %d and node %d\n", nodeAID, nodeBID)
				continue
			}
		} else { // add node to existing circuit
			circuitIdx := max(circuitAIndex, circuitBIndex)
			circuit := circuits[circuitIdx]
			circuit.nodes[nodeAID] = true
			circuit.nodes[nodeBID] = true
			totalConnectionsMade++
			fmt.Printf("added connection between %d and node %d\n", nodeAID, nodeBID)
			continue
		}

	}

	return lastDistanceNodePair
}

func Part2() {
	nodes := ParseInput()

	nodeDistances := make([]NodeDistance, 0)

	for i, nodeA := range nodes {
		for j := i + 1; j < len(nodes); j++ {
			nodeB := nodes[j]
			dist := Distance(nodeA, nodeB)
			nodeDistances = append(nodeDistances, NodeDistance{
				distance: dist,
				nodePair: [2]int{nodeA.id, nodeB.id},
			})
		}
	}

	// sort the nodeDistances by distance
	sort.Slice(nodeDistances, func(i, j int) bool {
		return nodeDistances[i].distance < nodeDistances[j].distance
	})

	// put these node distances into a queue
	nodeDistanceQueue := commons.Queue[NodeDistance]{}
	for _, nd := range nodeDistances {
		nodeDistanceQueue.Enqueue(nd)
	}

	// keep making connections until there is only one circuit, keep track of the last two pairs connected
	lastDistanceNodePair := makingConnectionsUntilAllConnected(nodeDistanceQueue, len(nodes))
	fmt.Println("last nodepair connected:", lastDistanceNodePair)
	// print the product of the x value of the last two nodes connected
	var nodeA, nodeB Node
	for _, node := range nodes {
		if node.id == lastDistanceNodePair.nodePair[0] {
			nodeA = node
		}
		if node.id == lastDistanceNodePair.nodePair[1] {
			nodeB = node
		}
	}
	fmt.Printf("The two nodes connected to make one large circuit is %+v, %+v\n", nodeA, nodeB)

	product := nodeA.x * nodeB.x
	fmt.Printf("Product of x values of last two nodes connected (%d and %d): %d\n", nodeA.id, nodeB.id, product)

}

func main() {
	fmt.Println("--- Day 8: Playground ---")
	//Part1(1000, 3)

	Part2()
}
