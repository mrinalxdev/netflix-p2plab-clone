// main.go
package main

import (
	"fmt"
	"log"
	"time"

	"ipld-benchmark/bench"
)

func main() {
	fmt.Println("Starting IPLD DAG Benchmarks (Custom IPLD)...")
	fmt.Println("\n========================================================")
	fmt.Println("--- Benchmarking Custom IPLD Implementation ---")
	fmt.Println("========================================================")
	fmt.Println("\n--- Benchmarking DAG Generation (Custom IPLD, 1000 nodes) ---")
	start := time.Now()
	customRootNode, customNodes, err := bench.GenerateCustomDAG(1000)
	if err != nil {
		log.Fatalf("Error generating DAG (Custom IPLD): %v", err)
	}
	duration := time.Since(start)
	fmt.Printf("DAG Generation (Custom IPLD, 1000 nodes): %s\n", duration)
	fmt.Printf("Root Custom CID: %x\n", customRootNode.Cid.Hash[:8]) 
	fmt.Printf("Total custom nodes generated: %d\n", len(customNodes))
	fmt.Println("\n--- Benchmarking Individual Node Creation (Custom IPLD, 1000 nodes) ---")
	bench.BenchmarkCustomNodeCreation(1000)
	fmt.Println("\n--- Benchmarking DAG Traversal (Custom IPLD, 1000 nodes) ---")
	if customRootNode == nil {
		log.Fatal("Root custom node is nil, cannot perform traversal benchmark (Custom IPLD).")
	}
	bench.BenchmarkCustomTraversal(customRootNode, customNodes) 
	fmt.Println("\n--- Benchmarking DAG Serialization (Custom IPLD, 1000 nodes) ---")
	bench.BenchmarkCustomSerialization(customRootNode)
	fmt.Println("\n--- Benchmarking DAG Deserialization (Custom IPLD, 1000 nodes) ---")
	bench.BenchmarkCustomDeserialization(customRootNode)

	fmt.Println("\nIPLD DAG Benchmarks Completed.")
}