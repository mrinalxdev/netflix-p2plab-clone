package bench

import (
	"fmt"
	"ipld-benchmark/myipld"
	"time"
)

/* {comment}
generateDag creates a simple linear dag using our custom myipld code
it creates numNodes nodes, where each node links to the previous one, forming a chain.
the last node created will be the root of this linear DAG
*/

func generateDag(numNodes int) (*myipld.MyNode, []*myipld.MyNode, error) {
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
		"content" : "custom-leaf-data-0",
		"created" : time.Now().UnixNano(),
	}

	leafNode, err := myipld.NewMyNode(leafData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create initial custom node : %w", err)
	}

	nodes = append(nodes, leafNode)
	prevNode = leafNode // the first node becomes the previous node

	for i := 1; i < numNodes; i ++ {
		nodeData := map[string]interface{}{
			"index" : i,
			"timestamp" : time.Now().UnixNano(),
			"message" : fmt.Sprintf("custom-node-%d-data", i),
		}

		currNode, err = myipld.NewMyNode(nodeData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create custom node %d : %w", i, err)
		}

		/* {comment}
		add a link from the current node to the previous node
		the link name here is illustrative; in real IPLD link names have significance

		shorting the cid for the link name
		*/

		linkName := fmt.Sprintf("custom-link-to-%x", prevNode.Cid.Hash[:8])
		if err := currNode.AddLink(linkName, prevNode.Cid); err != nil {
			return nil, nil, fmt.Errorf("failed to add link to custom node %d")
		}

		nodes = append(nodes, currNode)
		prevNode = currNode
	}

	return currNode, nodes, nil 
}
