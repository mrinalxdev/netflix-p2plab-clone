package bench

import (
	"fmt"
	"ipld-benchmark/myipld"
	"math/rand"
	"time"
)

type DAGStructure int

const (
	LinearDAG DAGStructure = iota

	BinaryTreeDAG
	StarDAG
	RandomDAG
)

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


func GenerateLinearDAG(numNodes int) (*myipld.MyNode, []*myipld.MyNode, error){
	if numNodes <= 0 {
		return nil, nil, fmt.Errorf("numNodes must be positive")
	}

	nodes := make([]*myipld.MyNode, 0, numNodes)
	var prevNode *myipld.MyNode


	for i := 0; i < numNodes; i++ {
		nodeData := map[string]interface{}{
			"index" : i,
			"timestamp": time.Now().UnixNano(),
			"message" : fmt.Sprintf("node-%d-data", i),
		}

		currNode, err := myipld.NewMyNode(nodeData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create node %d : %w", i , err)
		}


		if prevNode != nil {
			linkName := fmt.Sprintf("link-to-%w ", prevNode.Cid.Hash[:8])

			if err := currNode.AddLink(linkName , prevNode.Cid); err != nil {
				return nil, nil, fmt.Errorf("Failed to add link to node %d: %w", i, err)
			}
		}

		nodes = append(nodes, currNode)
		prevNode = currNode
	}

	return prevNode, nodes, nil
}


func GenerateBinaryTreeDAG(numNodes int) (*myipld.MyNode, []*myipld.MyNode, error){
	if numNodes <= 0 {
		return nil, nil, fmt.Errorf("numNodes must be positive")
	}

	nodes := make([]*myipld.MyNode, 0, numNodes)
	queue := make([]*myipld.MyNode, 0)


	rootData := map[string]interface {}{
		"index" : 0,
		"timestamp" : time.Now().UnixNano(),
		"message" : "root-node",
	}

	rootNode, err := myipld.NewMyNode(rootData)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create root node : %w", err)
	}


	nodes = append(nodes, rootNode)
	queue = append(queue, rootNode)


	index := 1

	for len(queue) > 0 && index < numNodes {
		current := queue[0]
		queue = queue[1:]


		// adding the left child to the most left node
		// why am I adding children to nodes ??
		// ....
		if index < numNodes {
			leftData := map[string]interface{}{
				"index" : index,
				"timestamp" : time.Now().UnixNano(),
				"message" : fmt.Sprintf("node-%d-data", index),
			}


			leftNode, err := myipld.NewMyNode(leftData)


			if err != nil {
				return nil, nil, fmt.Errorf("failed to create left child %d : %w", index, err)
			}


			linkName := fmt.Sprintf("left-child-%x", leftNode.Cid.Hash[:8])

			if err := current.AddLink(linkName, leftNode.Cid); err != nil {
				return nil, nil, fmt.Errorf("failed to add left link : %w", err)
			}

			nodes = append(nodes, leftNode)
			queue = append(queue, leftNode)

			index++

		}

		// adding the right child to the most right node


		if index < numNodes {
			rightData := map[string]interface{}{
				"index" : index,
				"timestamp": time.Now().UnixNano(),
				"message": fmt.Sprintf("node-%d-data", index),
			}
			rightNode, err := myipld.NewMyNode(rightData)

			if err != nil {
				return nil, nil, fmt.Errorf("failed to create right child %d : %w", index , err)
			}
			linkName := fmt.Sprintf("right-child-%x", rightNode.Cid.Hash[:8])


			if err := current.AddLink(linkName, rightNode.Cid); err != nil {
				return nil, nil, fmt.Errorf("failed to add right link : %w", err)
			}

			nodes = append(nodes, rightNode)
			queue = append(queue, rightNode)

			index++
		}
	}

	return rootNode, nodes, nil
}

// Okkay let me educate you on the StarDag
// these are just DAG's shaped in a star what did you think of 
// some rocket science ??

func GenerateStarDAG(numNodes int) (*myipld.MyNode, []*myipld.MyNode, error){
	if numNodes <= 0 {
		return nil, nil, fmt.Errorf("numNodes must be positive")
	}


	nodes := make([]*myipld.MyNode, 0, numNodes)


	centerData := map[string]interface{}{
		"index" : 0,
		"timestamp" : time.Now().UnixNano(),
		"message" : "center-node",
	}


	centerNode, err := myipld.NewMyNode(centerData)


	if err != nil {
		return nil, nil, fmt.Errorf("failed to create center node : %w", err)
	}


	nodes = append(nodes, centerNode)


	// creating the leaf nodes and linking them to center
	for i := 1; i < numNodes; i++ {
		leafData := map[string]interface{}{
			"index" : i,
			"timestamp" : time.Now().UnixNano(),
			"message" : fmt.Sprintf("leaf-node-%d", i),
		}


		leafNode, err := myipld.NewMyNode(leafData)


		if err != nil {
			return nil, nil, fmt.Errorf("failed to create leaf node %d : %w", i, err)
		}

		linkName := fmt.Sprintf("leaf-link-%x", leafNode.Cid.Hash[:8])
		if err := centerNode.AddLink(linkName, leafNode.Cid); err != nil {
			return nil, nil, fmt.Errorf("failed to add leaf link: %w", err)
		}
		nodes = append(nodes, leafNode)
	}

	return centerNode, nodes, nil
}

func GenerateRandomDAG(numNodes int, maxLinks int)(*myipld.MyNode, []*myipld.MyNode, error){
	if numNodes <= 0 {
		return nil, nil, fmt.Errorf("numNodes must be positive")
	}

	if maxLinks <= 0 {
		return nil, nil, fmt.Errorf("maxLinks must be positive")
	}


	// why seed ain't working normally like brooo com on T_T
	// kam karjaa bhai T_T

	rand.Seed(time.Now().UnixNano())

	nodes := make([]*myipld.MyNode, 0, numNodes)


	for i := 0; i < numNodes; i++ {
		nodeData := map[string]interface{}{
			"index" : i,
			"timestamp" : time.Now().UnixNano(),
			"message" : fmt.Sprintf("node-%d-data", i),
		}


		node, err := myipld.NewMyNode(nodeData)

		if err != nil {
			return nil, nil, fmt.Errorf("failed to create node %d : %w", i, err)
		}

		nodes = append(nodes, node)
	}


	for i := 1; i < numNodes; i++ {
		current := nodes[i]

		numLinks := rand.Intn(maxLinks) + 1

		for j := 0; j < numLinks; j++ {
			targetIndex := rand.Intn(i)
			targetNode := nodes[targetIndex]

			linkName := fmt.Sprintf("random-link-to-%x", targetNode.Cid.Hash[:8])

			if err := current.AddLink(linkName, targetNode.Cid); err != nil {
				return nil, nil, fmt.Errorf("failed to add random link : %w", err)
			}
		}
	}


	return nodes[0], nodes, nil
}












