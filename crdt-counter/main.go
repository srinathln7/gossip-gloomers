package main

import (
	// "counter/internal/node"
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	p2p "counter/internal/p2p/gossip"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Counter struct {
	value int
	mu    sync.Mutex
}

func main() {

	// Initialize the global counter and a global key
	counter := &Counter{value: 0}
	const key = "final_result"

	// Create a new maelstrom node and a seq.consistent key-value store
	node := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(node)

	gossip := p2p.NewGossip(100 * time.Millisecond)

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
		val, err := kv.ReadInt(context.Background(), key)
		if err != nil {
			val = 0
		}
		result := val + delta

		// Run a Compare and swap until it is successful
		// CompareAndSwap is a mechanism to perform updates in a concurrent environment while ensuring atomicity and avoiding race conditions.
		// If the current value matches the expected "from" value, it indicates that the key hasn't been modified by other operations, and the update can proceed.
		// If the current value doesn't match, it means that another operation has modified the key, and the update is skipped.

		for {

			// CompareAndSwap operation is atomic, meaning that the check for the current value and the update are treated as a single, indivisible operation.
			// Here, the operation goes through only the key "result" value still reads as the one captured by the `val` variable. If not, some other process
			// has already updated the key and keep re-trying
			err = kv.CompareAndSwap(context.Background(), key, val, result, true)
			if err == nil {
				break
			}
		}
		counter.mu.Unlock()

		gossip.Start(node, kv)

		return node.Reply(msg, map[string]interface{}{
			"type": "add_ok",
		})

	})

	node.Handle("gossip", func(msg maelstrom.Message) error {

		var body p2p.GossipMsg
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Value received from the peer nodes through the gossip protocol
		recvVal := body.Value

		counter.mu.Lock()
		counter.value = max(counter.value, recvVal)
		counter.mu.Unlock()

		return node.Reply(msg, body)

	})

	node.Handle("read", func(msg maelstrom.Message) error {

		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		counter.mu.Lock()

		existingValue, err := kv.ReadInt(context.Background(), key)

		if err != nil {
			body["value"] = 0
		}

		counter.value = max(existingValue, counter.value)
		body["value"] = counter.value

		counter.mu.Unlock()

		body["type"] = "read_ok"
		return node.Reply(msg, body)

	})

	if err := node.Run(); err != nil {
		gossip.Stop()
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}

}
