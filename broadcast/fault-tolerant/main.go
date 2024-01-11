package main

import (
	"broadcast/internal/node"
	p2p "broadcast/internal/p2p/gossip"
	"encoding/json"
	"log"
	"os"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {

	node := node.NewNode()
	gossip := p2p.NewGossip(100 * time.Millisecond)

	// Define handlers for different message types

	node.N.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		m := body["message"]
		v := int(m.(float64))
		node.Store.Set(v)

		gossip.Start(node)

		return node.N.Reply(msg, map[string]interface{}{
			"type": "broadcast_ok",
		})
	})

	node.N.Handle("gossip", func(msg maelstrom.Message) error {
		var body p2p.GossipMsgs
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		node.Store.SetAll(body.Msgs)

		return node.N.Reply(msg, body)
	})

	node.N.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["messages"] = node.Store.Get()
		body["type"] = "read_ok"

		return node.N.Reply(msg, body)
	})

	node.N.Handle("topology", func(msg maelstrom.Message) error {

		var body p2p.NWTopology

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		node.SetTopology(body.Topology[node.N.ID()])

		return node.N.Reply(msg, map[string]interface{}{
			"type": "topology_ok",
		})
	})

	if err := node.N.Run(); err != nil {
		gossip.Stop()
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}

}
