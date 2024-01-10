package main

import (
	"broadcast/internal/cache"
	p2p "broadcast/internal/p2p/gossip"
	"encoding/json"
	"log"
	"os"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {

	n := maelstrom.NewNode()

	kvstore := cache.NewCache()
	gossip := p2p.NewGossip(100 * time.Millisecond)

	// Define handlers for different message types

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		m := body["message"]
		v := int(m.(float64))
		kvstore.Set(v)

		//rsp := &model.SimpleResp{Type: "broadcast_ok"}
		return n.Reply(msg, map[string]interface{}{
			"type": "broadcast_ok",
		})
	})

	n.Handle("gossip", func(msg maelstrom.Message) error {
		var body p2p.GossipMsgs
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		kvstore.SetAll(body.Msgs)

		return n.Reply(msg, body)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["messages"] = kvstore.Get()
		body["type"] = "read_ok"

		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		gossip.Start(kvstore, n, 5)
		return n.Reply(msg, map[string]interface{}{
			"type": "topology_ok",
		})
	})

	if err := n.Run(); err != nil {
		gossip.Stop()
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}

}
