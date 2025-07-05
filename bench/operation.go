package bench

import (
	"fmt"
	"ipld-benchmark/myipld"
	"log"
	"time"
)

/*

-- Benchmark for custom ipld implementation

*/

func BenchmarkCustomNodeCreation(numNodes int){
	start := time.Now()
	for i := 0; i < numNodes; i++ {
		data := map[string]interface{}{
			"id" : i,
			"timestamp" : time.Now().UnixNano(),
			"random" : "some-random-string-to-vary-data-size-and-hash",
		}

		_, err := myipld.NewMyNode(data)
		if err != nil {
			log.Printf("error creating custom node %d : %v", i, err)
			return
		}
	}

	duration := time.Since(start)
	fmt.Printf("	Node creation (%d custom ipld nodes ): %s\n", numNodes, duration)
}