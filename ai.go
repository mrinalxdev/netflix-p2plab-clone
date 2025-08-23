package bench

import (
	"fmt"
	"math/rand"
	"time"

	"ipld-benchmark/myipld"
)

// DAGStructure represents the type of DAG to generate.
type DAGStructure int

const (
	LinearDAG DAGStructure = iota
	BinaryTreeDAG
	StarDAG
	RandomDAG
)

// String returns a human-readable name for the DAG structure.
func (d DAGStructure) String() string {
	switch d {
	case LinearDAG:
		return "LinearDAG"
	case BinaryTreeDAG:
		return "BinaryTreeDAG"
	case StarDAG:
		return "StarDAG"
	case RandomDAG:
		return "RandomDAG"
	default:
		return "UnknownDAG"
	}
}

// GenerateDAG generates a DAG of the specified structure and number of nodes.
func GenerateDAG(structure DAGStructure, numNodes int) (*myipld.MyNode, []*myipld.MyNode, error) {
	switch structure {
	case LinearDAG:
		return GenerateLinearDAG(numNodes)
	case BinaryTreeDAG:
		return GenerateBinaryTreeDAG(numNodes)
	case StarDAG:
		return GenerateStarDAG(numNodes)
	case RandomDAG:
		return GenerateRandomDAG(numNodes, 3) // default maxLinks=3
	default:
		return nil, nil, fmt.Errorf("unknown DAG structure")
	}
}

// GenerateLinearDAG generates a linear DAG where each node points to the previous one.
func GenerateLinearDAG(numNodes int) (*myipld.MyNode, []*myipld.MyNode, error) {
	if numNodes <= 0 {
		return nil, nil, fmt.Errorf("numNodes must be positive")
	}

	nodes := make([]*myipld.MyNode, 0, numNodes)
	var prevNode *myipld.MyNode

	for i := 0; i < numNodes; i++ {
		nodeData := map[string]interface{}{
			"index":     i,
			"timestamp": time.Now().UnixNano(),
			"message":   fmt.Sprintf("node-%d-data", i),
		}
		currNode, err := myipld.NewMyNode(nodeData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create node %d: %w", i, err)
		}

		if prevNode != nil {
			linkName := fmt.Sprintf("link-to-%x", prevNode.Cid.Hash[:8])
			if err := currNode.AddLink(linkName, prevNode.Cid); err != nil {
				return nil, nil, fmt.Errorf("failed to add link to node %d: %w", i, err)
			}
		}

		nodes = append(nodes, currNode)
		prevNode = currNode
	}

	return prevNode, nodes, nil
}

// GenerateBinaryTreeDAG generates a binary tree DAG.
func GenerateBinaryTreeDAG(numNodes int) (*myipld.MyNode, []*myipld.MyNode, error) {
	if numNodes <= 0 {
		return nil, nil, fmt.Errorf("numNodes must be positive")
	}

	nodes := make([]*myipld.MyNode, 0, numNodes)
	queue := make([]*myipld.MyNode, 0)

	// Create root node
	rootData := map[string]interface{}{
		"index":     0,
		"timestamp": time.Now().UnixNano(),
		"message":   "root-node",
	}
	rootNode, err := myipld.NewMyNode(rootData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create root node: %w", err)
	}
	nodes = append(nodes, rootNode)
	queue = append(queue, rootNode)

	index := 1
	for len(queue) > 0 && index < numNodes {
		current := queue[0]
		queue = queue[1:]

		// Add left child
		if index < numNodes {
			leftData := map[string]interface{}{
				"index":     index,
				"timestamp": time.Now().UnixNano(),
				"message":   fmt.Sprintf("node-%d-data", index),
			}
			leftNode, err := myipld.NewMyNode(leftData)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to create left child %d: %w", index, err)
			}
			linkName := fmt.Sprintf("left-child-%x", leftNode.Cid.Hash[:8])
			if err := current.AddLink(linkName, leftNode.Cid); err != nil {
				return nil, nil, fmt.Errorf("failed to add left link: %w", err)
			}
			nodes = append(nodes, leftNode)
			queue = append(queue, leftNode)
			index++
		}

		// Add right child
		if index < numNodes {
			rightData := map[string]interface{}{
				"index":     index,
				"timestamp": time.Now().UnixNano(),
				"message":   fmt.Sprintf("node-%d-data", index),
			}
			rightNode, err := myipld.NewMyNode(rightData)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to create right child %d: %w", index, err)
			}
			linkName := fmt.Sprintf("right-child-%x", rightNode.Cid.Hash[:8])
			if err := current.AddLink(linkName, rightNode.Cid); err != nil {
				return nil, nil, fmt.Errorf("failed to add right link: %w", err)
			}
			nodes = append(nodes, rightNode)
			queue = append(queue, rightNode)
			index++
		}
	}

	return rootNode, nodes, nil
}

// GenerateStarDAG generates a star DAG with a central node connected to all others.
func GenerateStarDAG(numNodes int) (*myipld.MyNode, []*myipld.MyNode, error) {
	if numNodes <= 0 {
		return nil, nil, fmt.Errorf("numNodes must be positive")
	}

	nodes := make([]*myipld.MyNode, 0, numNodes)

	// Create center node
	centerData := map[string]interface{}{
		"index":     0,
		"timestamp": time.Now().UnixNano(),
		"message":   "center-node",
	}
	centerNode, err := myipld.NewMyNode(centerData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create center node: %w", err)
	}
	nodes = append(nodes, centerNode)

	// Create leaf nodes and link to center
	for i := 1; i < numNodes; i++ {
		leafData := map[string]interface{}{
			"index":     i,
			"timestamp": time.Now().UnixNano(),
			"message":   fmt.Sprintf("leaf-node-%d", i),
		}
		leafNode, err := myipld.NewMyNode(leafData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create leaf node %d: %w", i, err)
		}
		linkName := fmt.Sprintf("leaf-link-%x", leafNode.Cid.Hash[:8])
		if err := centerNode.AddLink(linkName, leafNode.Cid); err != nil {
			return nil, nil, fmt.Errorf("failed to add leaf link: %w", err)
		}
		nodes = append(nodes, leafNode)
	}

	return centerNode, nodes, nil
}

// GenerateRandomDAG generates a random DAG with each node having up to maxLinks links.
func GenerateRandomDAG(numNodes int, maxLinks int) (*myipld.MyNode, []*myipld.MyNode, error) {
	if numNodes <= 0 {
		return nil, nil, fmt.Errorf("numNodes must be positive")
	}
	if maxLinks <= 0 {
		return nil, nil, fmt.Errorf("maxLinks must be positive")
	}

	rand.Seed(time.Now().UnixNano())
	nodes := make([]*myipld.MyNode, 0, numNodes)

	// Create all nodes first
	for i := 0; i < numNodes; i++ {
		nodeData := map[string]interface{}{
			"index":     i,
			"timestamp": time.Now().UnixNano(),
			"message":   fmt.Sprintf("node-%d-data", i),
		}
		node, err := myipld.NewMyNode(nodeData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create node %d: %w", i, err)
		}
		nodes = append(nodes, node)
	}

	// For each node, randomly select up to maxLinks previous nodes to link to
	for i := 1; i < numNodes; i++ {
		current := nodes[i]
		// Number of links for this node (at least 1, up to maxLinks)
		numLinks := rand.Intn(maxLinks) + 1
		// Select random nodes from indices [0, i-1]
		for j := 0; j < numLinks; j++ {
			targetIndex := rand.Intn(i)
			targetNode := nodes[targetIndex]
			linkName := fmt.Sprintf("random-link-to-%x", targetNode.Cid.Hash[:8])
			if err := current.AddLink(linkName, targetNode.Cid); err != nil {
				return nil, nil, fmt.Errorf("failed to add random link: %w", err)
			}
		}
	}

	return nodes[0], nodes, nil
}