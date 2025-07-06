// bench/dag_generator.go
package bench

import (
	"fmt"
	"time"

	"ipld-benchmark/myipld"
)


func GenerateCustomDAG(numNodes int) (*myipld.MyNode, []*myipld.MyNode, error) {
	if numNodes <= 0 {
		return nil, nil, fmt.Errorf("numNodes must be positive")
	}

	var (
		nodes    []*myipld.MyNode
		prevNode *myipld.MyNode   
		currNode *myipld.MyNode   
		err      error
	)


	leafData := map[string]interface{}{
		"content": "custom-leaf-data-0",
		"created": time.Now().UnixNano(),
	}
	leafNode, err := myipld.NewMyNode(leafData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create initial custom node: %w", err)
	}
	nodes = append(nodes, leafNode)
	prevNode = leafNode

	// creating subsequent nodes, each linking to the immediately previous node
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

		linkName := fmt.Sprintf("custom-link-to-%x", prevNode.Cid.Hash[:8])
		if err := currNode.AddLink(linkName, prevNode.Cid); err != nil {
			return nil, nil, fmt.Errorf("failed to add link to custom node %d: %w", i, err)
		}

		nodes = append(nodes, currNode)
		prevNode = currNode           
	}

	return currNode, nodes, nil
}