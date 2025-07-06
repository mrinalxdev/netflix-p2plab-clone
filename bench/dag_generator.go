// bench/dag_generator.go
package bench

import (
	"fmt"
	"time"

	"ipld-benchmark/myipld" // Import our custom IPLD package
)

// GenerateCustomDAG creates a simple linear DAG using our custom myipld implementation.
// It creates `numNodes` nodes, where each node links to the previous one,
// forming a chain. The last node created will be the root of this linear DAG.
//
// Returns:
//   - *myipld.MyNode: The root node of the generated DAG.
//   - []*myipld.MyNode: A slice containing all generated nodes, useful for traversal simulation.
//   - error: An error if DAG generation fails.
func GenerateCustomDAG(numNodes int) (*myipld.MyNode, []*myipld.MyNode, error) {
	if numNodes <= 0 {
		return nil, nil, fmt.Errorf("numNodes must be positive")
	}

	var (
		nodes    []*myipld.MyNode // Slice to hold all generated nodes
		prevNode *myipld.MyNode   // To keep track of the previous node for linking
		currNode *myipld.MyNode   // The current node being created
		err      error
	)

	// Create the first node (often considered a "leaf" in a linear chain)
	leafData := map[string]interface{}{
		"content": "custom-leaf-data-0",
		"created": time.Now().UnixNano(),
	}
	leafNode, err := myipld.NewMyNode(leafData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create initial custom node: %w", err)
	}
	nodes = append(nodes, leafNode)
	prevNode = leafNode // The first node becomes the previous node for the next iteration

	// Create subsequent nodes, each linking to the immediately previous node
	for i := 1; i < numNodes; i++ {
		nodeData := map[string]interface{}{
			"index":     i,
			"timestamp": time.Now().UnixNano(),
			"message":   fmt.Sprintf("custom-node-%d-data", i),
		}
		currNode, err = myipld.NewMyNode(nodeData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create custom node %d: %w", i, err)
		}

		// Add a link from the current node to the previous node
		// The link name here is illustrative; in real IPLD, link names have significance.
		linkName := fmt.Sprintf("custom-link-to-%x", prevNode.Cid.Hash[:8]) // Shortened CID for link name
		if err := currNode.AddLink(linkName, prevNode.Cid); err != nil {
			return nil, nil, fmt.Errorf("failed to add link to custom node %d: %w", i, err)
		}

		nodes = append(nodes, currNode) // Add the newly created node to our list
		prevNode = currNode             // The current node becomes the previous for the next iteration
	}

	// The last node created in the loop is the root of this linear DAG
	return currNode, nodes, nil
}