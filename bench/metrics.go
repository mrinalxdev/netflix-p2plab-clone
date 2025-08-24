package bench

import (
	"runtime"
	"time"

	"ipld-benchmark/myipld"
)

type PerformanceMetrics struct {
	TotalTime       time.Duration
	NodesPerSecond  float64
	MemoryAlloc     uint64
	MemoryTotal     uint64
	GCPercentage    float64
	SerializedSize  int
}

type DAGMetrics struct {
	PerformanceMetrics
	MaxDepth         int
	AverageDepth     float64
	MaxBreadth       int
	LinkDensity      float64
	Diameter         int
}

func CollectMetrics(runFunc func() error) (*PerformanceMetrics, error) {
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	start := time.Now()
	err := runFunc()
	duration := time.Since(start)

	runtime.ReadMemStats(&m2)

	if err != nil {
		return nil, err
	}

	nodesPerSecond := 0.0
	if duration.Seconds() > 0 {
		// This will be adjusted by the caller based on actual node count
		nodesPerSecond = 1 / duration.Seconds()
	}

	return &PerformanceMetrics{
		TotalTime:      duration,
		NodesPerSecond: nodesPerSecond,
		MemoryAlloc:    m2.Alloc - m1.Alloc,
		MemoryTotal:    m2.TotalAlloc - m1.TotalAlloc,
	}, nil
}

func AnalyzeDAGStructure(root *myipld.MyNode, allNodes []*myipld.MyNode) *DAGMetrics {
	depths := make(map[myipld.MyCID]int)
	maxDepth := 0
	totalDepth := 0

	// BFS to calculate depths
	queue := []*myipld.MyNode{root}
	depths[root.Cid] = 0
	visited := make(map[myipld.MyCID]bool)

	
	visited[root.Cid] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		currentDepth := depths[current.Cid]

		for _, link := range current.Links {
			if !visited[link.Cid] {
				var node *myipld.MyNode
				for _, n := range allNodes {
					if n.Cid == link.Cid {
						node = n
						break
					}
				}
				if node != nil {
					depths[node.Cid] = currentDepth + 1
					if currentDepth+1 > maxDepth {
						maxDepth = currentDepth + 1
					}
					totalDepth += currentDepth + 1

					
					visited[node.Cid] = true
					queue = append(queue, node)
				}
			}
		}
	}

	numNodes := len(visited)
	averageDepth := 0.0
	if numNodes > 0 {
		averageDepth = float64(totalDepth) / float64(numNodes)
	}
	maxBreadth := 0
	for _, node := range allNodes {
		if len(node.Links) > maxBreadth {
			maxBreadth = len(node.Links)
		}
	}

	totalLinks := 0
	for _, node := range allNodes {
		totalLinks += len(node.Links)
	}
	linkDensity := float64(totalLinks) / float64(numNodes*(numNodes-1))
	diameter := maxDepth

	return &DAGMetrics{
		MaxDepth:     maxDepth,
		AverageDepth: averageDepth,
		MaxBreadth:   maxBreadth,
		LinkDensity:  linkDensity,
		Diameter:     diameter,
	}
}