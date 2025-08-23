package bench

import (
	"fmt"
	"ipld-benchmark/myipld"
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
}
