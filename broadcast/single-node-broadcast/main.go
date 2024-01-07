package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var cache map[string][]int = make(map[string][]int)

func main() {

	// Create a maelstrom node to begin with
	n := maelstrom.NewNode()

	n.Handle("broadcast", func(msg maelstrom.Message) error {

		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Store the broadcasted values locally in each node against its unique identifier in-memory so that it can be read later
		cache[n.ID()] = append(cache[n.ID()], int(body["message"].(float64)))

		return n.Reply(msg, map[string]any{
			"type": "broadcast_ok",
		})
	})

	n.Handle("read", func(msg maelstrom.Message) error {

		var res map[string]any

		if err := json.Unmarshal(msg.Body, &res); err != nil {
			return err
		}

		res["type"] = "read_ok"
		res["messages"] = cache[n.ID()]

		return n.Reply(msg, res)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]any{
			"type": "topology_ok",
		})
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}

}
