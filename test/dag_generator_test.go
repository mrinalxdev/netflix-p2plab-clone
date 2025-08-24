package test

import (
	"ipld-benchmark/bench"
	"testing"
)

func TestDAGSStructure(t *testing.T) {
	testCases := []struct {
		name      string
		structure bench.DAGStructure
		nodeCount int
	}{
		{"LinearDAG-100", bench.LinearDAG, 100},
		{"BinaryTreeDAG-100", bench.BinaryTreeDAG, 100},
		{"StarDAG-100", bench.StarDAG, 100},
		{"RandomDAG-100", bench.RandomDAG, 100},
		{"LinearDAG-1000", bench.LinearDAG, 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			root, nodes, err := bench.GenerateDAG(tc.structure, tc.nodeCount)

			if err != nil {
				t.Fatalf("Failed to generate %s: %v", tc.name, err)
			}

			if root == nil {
				t.Fatalf("Root node is nil for %s", tc.name)
			}

			if len(nodes) != tc.nodeCount {
				t.Fatalf("Expected %d nodes, got %d for %s", tc.nodeCount, len(nodes), tc.name)
			}

			// verifying if all nodes have valid CIDs
			// for i, node := range nodes {
			// 	if node.Cid.Hash == [32]byte{}{
			// 		t.Errorf("Node %d has empty CID in %s", i, tc.name)
			// 	}
			// }

			for i, node := range nodes {
				var emptyHash [32]byte

				if node.Cid.Hash == emptyHash {
					t.Errorf("Node %d has empty CID in %s", i, tc.name)
				}
			}

			metrics := bench.AnalyzeDAGStructure(root, nodes)

			if metrics.MaxDepth <= 0 {
				t.Errorf("Invalid max depth %d for %s", metrics.MaxDepth, tc.name)
			}
		})
	}
}

func TestInvalidDAGParameter(t *testing.T) {
	_, _, err := bench.GenerateDAG(bench.LinearDAG, 0)
	if err == nil {
		t.Error("Expected error for zero nodes")
	}

	_, _, err = bench.GenerateDAG(bench.RandomDAG, 10)
	if err == nil {
		t.Error("Expected error for zero maxLimits")
	}
}

func TestDAGStructureString(t *testing.T) {
	testCases := []struct {
		structure bench.DAGStructure
		expected  string
	}{
		{bench.LinearDAG, "LinearDAG"},
		{bench.BinaryTreeDAG, "BinaryTreeDAG"},
		{bench.StarDAG, "StarDAG"},
		{bench.RandomDAG, "RandomDAG"},
		{bench.DAGStructure(99), "UnknownDAG"},
	}

	for _, tc := range testCases {
		result := tc.structure.String()
		if result != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, result)
		}
	}
}
