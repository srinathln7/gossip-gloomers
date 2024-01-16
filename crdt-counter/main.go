package main

import (
	// "counter/internal/node"
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// Counter: Represents the Global Counter
type Counter struct {
	value int
	mu    sync.RWMutex
}

const GLOBAL_RESULT_KEY = "final_result"

func main() {

	// Create a new maelstrom node and a seq.consistent key-value store
	node := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(node)

	// Initialize the local counter
	counter := &Counter{value: 0}

	node.Handle("add", func(msg maelstrom.Message) error {

		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		data := body["delta"]
		delta := int(data.(float64))

		// Add the value

		counter.mu.Lock()
		// Retreive the current counter value
		val, err := kv.ReadInt(context.Background(), GLOBAL_RESULT_KEY)
		if err != nil {
			val = 0
		}
		counter.mu.Unlock()

		result := val + delta

		// Run a Compare and swap until it is successful
		// CompareAndSwap is a mechanism to perform updates in a concurrent environment while ensuring atomicity and avoiding race conditions.
		// If the current value matches the expected "from" value, it indicates that the key hasn't been modified by other operations, and the update can proceed.
		// If the current value doesn't match, it means that another operation has modified the key, and the update is skipped.

		for {

			// CompareAndSwap operation is atomic, meaning that the check for the current value and the update are treated as a single, indivisible operation.
			// Here, the operation goes through only the key "result" value still reads as the one captured by the `val` variable. If not, some other process
			// has already updated the key and keep re-trying
			err = kv.CompareAndSwap(context.Background(), GLOBAL_RESULT_KEY, val, result, true)
			if err == nil {
				break
			}
		}

		return node.Reply(msg, map[string]interface{}{
			"type": "add_ok",
		})

	})

	node.Handle("read", func(msg maelstrom.Message) error {

		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		counter.mu.Lock()

		existingValue, err := kv.ReadInt(context.Background(), GLOBAL_RESULT_KEY)

		if err != nil {
			body["value"] = 0
		}

		counter.mu.Unlock()

		body["value"] = existingValue
		body["type"] = "read_ok"
		return node.Reply(msg, body)

	})

	if err := node.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}

}
