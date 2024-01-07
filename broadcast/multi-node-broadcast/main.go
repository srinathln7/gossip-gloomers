package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// Node represents a node in the distributed system.
type Node struct {
	n  *maelstrom.Node // Reference to the Maelstrom node
	mu sync.RWMutex    // Mutex for concurrent access to shared resources
	wg sync.WaitGroup  // WaitGroup to wait for goroutines to finish

	id      string           // Unique identifier for the node
	recvMsg map[int]struct{} // Map to track received broadcast messages
}

// TopologyMsg represents the JSON message structure for network topology.
type TopologyMsg struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

// NwTopology stores the network topology shared across all nodes.
var NwTopology map[string][]string = make(map[string][]string)

func main() {
	// Create a Maelstrom node
	n := maelstrom.NewNode()
	node := Node{
		n:       n,
		id:      n.ID(),
		recvMsg: make(map[int]struct{}),
	}

	// Define handlers for different message types
	n.Handle("broadcast", node.hBroadcast)
	n.Handle("read", node.hRead)
	n.Handle("topology", node.hTopology)

	// Run the Maelstrom node
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

// hBroadcast handles the "broadcast" message type.
func (node *Node) hBroadcast(msg maelstrom.Message) error {
	var body map[string]interface{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	// Extract the message value
	val := int(body["message"].(float64))

	// Check if the node has already received or broadcasted the message, discard if yes
	node.mu.RLock()
	if _, exists := node.recvMsg[val]; exists {
		node.mu.RUnlock()
		return nil
	}
	node.mu.RUnlock()

	// Support concurrent writes on the map
	// If the message has not been broadcasted before, record it in the node's memory
	node.mu.Lock()
	node.recvMsg[val] = struct{}{}
	node.mu.Unlock()

	// Retrieve the neighbors list of a particular node
	neighbours := NwTopology[node.n.ID()]

	// Broadcast the values to the neighbors
	node.wg.Add(len(neighbours))
	for _, neighbour := range neighbours {
		// Spin up a separate goroutine to broadcast each message to the neighbor
		go func(neighbour string) {
			node.n.Send(neighbour, body)
			node.wg.Done()
		}(neighbour)
	}

	// Wait for all broadcasts to finish
	node.wg.Wait()

	// Return the standard broadcast response
	return node.n.Reply(msg, map[string]interface{}{
		"type": "broadcast_ok",
	})
}

// hRead handles the "read" message type.
func (node *Node) hRead(msg maelstrom.Message) error {
	var res map[string]interface{}

	if err := json.Unmarshal(msg.Body, &res); err != nil {
		return err
	}

	// Prepare the response with the list of received messages
	res["type"] = "read_ok"
	var msgs []int

	node.mu.RLock()
	for val := range node.recvMsg {
		msgs = append(msgs, val)
	}
	node.mu.RUnlock()

	res["messages"] = msgs
	return node.n.Reply(msg, res)
}

// hTopology handles the "topology" message type.
func (node *Node) hTopology(msg maelstrom.Message) error {
	var body TopologyMsg

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	// Store the network topology in-memory
	NwTopology = body.Topology

	return node.n.Reply(msg, map[string]interface{}{
		"type": "topology_ok",
	})
}
