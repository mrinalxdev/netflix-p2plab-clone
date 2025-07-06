package bench

import (
	// "bytes"
	"fmt"
	"log"
	"time"

	"ipld-benchmark/myipld"
)

// --- Benchmarks for Custom IPLD Implementation ---

/* {comment}
BenchmarkCustomNodeCreation measures the time to create `numNodes` individual custom IPLD nodes.
This benchmark focuses on the overhead of creating new `myipld.MyNode` instances,
including data marshalling and CID computation.
*/

func BenchmarkCustomNodeCreation(numNodes int) {
	start := time.Now()
	for i := 0; i < numNodes; i++ {
		data := map[string]interface{}{
			"id":        i,
			"timestamp": time.Now().UnixNano(),
			"random":    "some-random-string-to-vary-data-size-and-hash",
		}
		_, err := myipld.NewMyNode(data)
		if err != nil {
			log.Printf("Error creating custom node %d: %v", i, err)
			return
		}
	}
	duration := time.Since(start)
	fmt.Printf("  Node Creation (%d Custom IPLD nodes): %s\n", numNodes, duration)
}

/* {comment}

BenchmarkCustomTraversal measures the time to traverse a custom IPLD DAG.

	Since we don't have a persistent block store, this simulates traversal
	by looking up linked nodes in an in-memory map of all generated nodes.
	In a real IPLD system, this would involve fetching blocks from a store.

	Parameters:
	  - rootNode: The starting node for traversal.
	  - allNodes: A slice of all nodes in the DAG, used to simulate a node store.
*/


func BenchmarkCustomTraversal(rootNode *myipld.MyNode, allNodes []*myipld.MyNode) {
	if rootNode == nil || len(allNodes) == 0 {
		fmt.Println("  Cannot benchmark custom traversal: root node or allNodes is nil/empty.")
		return
	}

	visitedNodes := 0
	nodeMap := make(map[myipld.MyCID]*myipld.MyNode)
	for _, node := range allNodes {
		nodeMap[node.Cid] = node
	}

	start := time.Now()
	stack := []*myipld.MyNode{rootNode}       
	visitedCIDs := make(map[myipld.MyCID]bool)

	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if visitedCIDs[curr.Cid] {
			continue
		}
		visitedCIDs[curr.Cid] = true
		visitedNodes++              
		for _, link := range curr.Links {
			if !visitedCIDs[link.Cid] {
				if linkedNode, ok := nodeMap[link.Cid]; ok {
					stack = append(stack, linkedNode)
				}
				// next todo is to try to fetch it from a persistent store
			}
		}
	}

	duration := time.Since(start)
	fmt.Printf("  DAG Traversal (Custom IPLD, visited %d nodes): %s\n", visitedNodes, duration)
}


func BenchmarkCustomSerialization(node *myipld.MyNode) {
	if node == nil {
		fmt.Println("  Cannot benchmark custom serialization: node is nil.")
		return
	}

	start := time.Now()
	_, err := node.ToBytes()
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("  Error during custom node serialization: %v\n", err)
	} else {
		fmt.Printf("  Node Serialization (Custom IPLD): %s\n", duration)
	}
}

/* {comment}
BenchmarkCustomDeserialization measures the time to deserialize a custom IPLD node.
This involves converting a byte slice (JSON) back into a `myipld.MyNode` structure. 
*/

func BenchmarkCustomDeserialization(node *myipld.MyNode) {
	if node == nil {
		fmt.Println("  Cannot benchmark custom deserialization: node is nil.")
		return
	}

	data, err := node.ToBytes()
	if err != nil {
		fmt.Printf("  Error preparing data for custom deserialization: %v\n", err)
		return
	}

	start := time.Now()
	_, err = myipld.FromBytes(data) // Deserialize the bytes back into a node
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("  Error during custom node deserialization: %v\n", err)
	} else {
		fmt.Printf("  Node Deserialization (Custom IPLD, from %d bytes): %s\n", len(data), duration)
	}
}
