package bench

import (
	"fmt"
	"log"
	"time"

	"ipld-benchmark/myipld"
)


func BenchmarkDAGOperations(structure DAGStructure, numNodes int) (*PerformanceMetrics, *DAGMetrics, error) {
	generateMetrics, err := CollectMetrics(func() error {
		_, _, err := GenerateDAG(structure, numNodes)
		return err
	})
	if err != nil {
		return nil, nil, fmt.Errorf("DAG generation failed: %w", err)
	}
	generateMetrics.NodesPerSecond = float64(numNodes) / generateMetrics.TotalTime.Seconds()

	root, nodes, err := GenerateDAG(structure, numNodes)
	if err != nil {
		return nil, nil, err
	}

	traversalMetrics, err := CollectMetrics(func() error {
		BenchmarkTraversal(root, nodes)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	serializationMetrics, err := CollectMetrics(func() error {
		BenchmarkSerialization(root)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	deserializationMetrics, err := CollectMetrics(func() error {
		BenchmarkDeserialization(root)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	combinedMetrics := &PerformanceMetrics{
		TotalTime:      generateMetrics.TotalTime + traversalMetrics.TotalTime + serializationMetrics.TotalTime + deserializationMetrics.TotalTime,
		NodesPerSecond: generateMetrics.NodesPerSecond,
		MemoryAlloc:    generateMetrics.MemoryAlloc + traversalMetrics.MemoryAlloc + serializationMetrics.MemoryAlloc + deserializationMetrics.MemoryAlloc,
		MemoryTotal:    generateMetrics.MemoryTotal + traversalMetrics.MemoryTotal + serializationMetrics.MemoryTotal + deserializationMetrics.MemoryTotal,
	}
	dagMetrics := AnalyzeDAGStructure(root, nodes)

	return combinedMetrics, dagMetrics, nil
}

func BenchmarkTraversal(root *myipld.MyNode, allNodes []*myipld.MyNode) {
	if root == nil || len(allNodes) == 0 {
		return
	}

	visitedNodes := 0
	nodeMap := make(map[myipld.MyCID]*myipld.MyNode)
	for _, node := range allNodes {
		nodeMap[node.Cid] = node
	}

	stack := []*myipld.MyNode{root}
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
			}
		}
	}
}

func BenchmarkSerialization(node *myipld.MyNode) {
	if node == nil {
		return
	}
	node.ToBytes()
}

func BenchmarkDeserialization(node *myipld.MyNode) {
	if node == nil {
		return
	}
	data, err := node.ToBytes()
	if err != nil {
		return
	}
	myipld.FromBytes(data)
}

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