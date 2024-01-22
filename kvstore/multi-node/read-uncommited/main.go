package main

import (
	"encoding/json"
	"log"
	"os"

	kvstore "kvstore/internal/kvs"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {

	node := maelstrom.NewNode()
	kvs := kvstore.NewKVStore()

	node.Handle("txn", func(msg maelstrom.Message) error {

		var body map[string]any
		err := json.Unmarshal(msg.Body, &body)
		if err != nil {
			return err
		}

		var resp [][3]any

		reqs := body["txn"].([]any)
		for _, req := range reqs {
			reqBody := req.([]any)
			reqType := reqBody[0].(string)
			reqKey := int(reqBody[1].(float64))

			switch reqType {
			case "w":
				reqValue := int(reqBody[2].(float64))
				kvs.SetAndBroadcast(node, reqKey, reqValue)
				res := [3]any{reqType, reqKey, reqValue}
				resp = append(resp, res)

			case "r":
				val := kvs.Get(reqKey)
				if val != nil {
					res := [3]any{reqType, reqKey, val.Val}
					resp = append(resp, res)
				} else {
					res := [3]any{reqType, reqKey, nil}
					resp = append(resp, res)
				}
			}
		}

		return node.Reply(msg, map[string]any{
			"type": "txn_ok",
			"txn":  resp,
		})
	})

	node.Handle("gossip", func(msg maelstrom.Message) error {

		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		key := int(body["key"].(float64))
		val := int(body["value"].(float64))
		version := int(body["version"].(float64))

		kvs.SyncLocal(key, val, version)
		return node.Reply(msg, map[string]any{
			"type": "gossip_ok",
		})
	})

	if err := node.Run(); err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}
}
