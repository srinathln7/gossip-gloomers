package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// Node represents a node in the distributed system.
type Node struct {
	n  *maelstrom.Node // Reference to the Maelstrom node
	wg sync.WaitGroup  // WaitGroup to wait for goroutines to finish

	br *BroadcastSvc // Broadcast service

	id      string           // Unique identifier for the node
	recvMsg map[int]struct{} // Map to track received broadcast messages
}

type Broadcast struct {
	dest string
	val  map[string]any // val to be broadcasted
}

type BroadcastSvc struct {
	ch chan Broadcast
}

// TopologyMsg represents the JSON message structure for network topology.
type TopologyMsg struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

// NwTopology stores the network topology shared across all nodes.
var NwTopology map[string][]string = make(map[string][]string)

// initBroadcastSvc : Initializes a broadcasting service for a given node with a specified number of worker goroutines
// Each worker goroutine waits for messages on the channel and broadcasts them to neighbors
func initBroadcastSvc(ctx context.Context, n *maelstrom.Node, numOfWorkers int) *BroadcastSvc {

	// init a channel to send and receive the msg
	ch := make(chan Broadcast)

	for i := 0; i < numOfWorkers; i++ {

		go func() {
			for {
				select {
				case msg := <-ch:

					// Run an infinite for loop until the msg is broadcasted
					for {
						if err := n.Send(msg.dest, msg.val); err != nil {
							continue
						}

						break
					}

				// `Done()` returns a closed channel when work done on behalf of this context should be cancelled.
				// Done may return `nil` if this context can never be cancelled. Successive calls to done return the same value.
				// Rem: `read` from a closed channel will return `default` value while `read` from a `nil` channel will BLOCK.
				// This cond'n gets UNBLOCKED because the parent call `genGreeting` cancels the context after a deadline of 1s.
				case <-ctx.Done():
					return
				}
			}
		}()

	}

	return &BroadcastSvc{
		ch: ch,
	}

}

func main() {

	// Create a Maelstrom node
	n := maelstrom.NewNode()

	//Create a context with cancel function. `WithCancel` fns return a new `context` that closes its `done` channel when the returned `cancel` function is called
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a node service
	node := Node{
		n:  n,
		id: n.ID(),

		br:      initBroadcastSvc(ctx, n, 10),
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
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	// Extract the message value
	val := int(body["message"].(float64))

	// Check if the node has already received or broadcasted the message, discard if yes
	if _, exists := node.recvMsg[val]; exists {
		return nil
	}

	// Otherwise, record it in the node's memory
	node.recvMsg[val] = struct{}{}

	// Retrieve the neighbors list of a particular node
	neighbours := NwTopology[node.n.ID()]

	// Broadcast the values to the neighbors
	node.wg.Add(len(neighbours))
	for _, neighbour := range neighbours {
		// Spin up a separate goroutine to broadcast each message to the neighbor
		go func(neighbour string) {
			node.br.ch <- Broadcast{dest: neighbour, val: body}
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

	for val := range node.recvMsg {
		msgs = append(msgs, val)
	}

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
