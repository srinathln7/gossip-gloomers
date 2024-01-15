package main

import (
	"encoding/json"
	"log"
	"logs/internal/replicator"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {

	replicator := replicator.NewReplicator()

	replicator.Node.Handle("send", func(msg maelstrom.Message) error {
		var body map[string]any

		err := json.Unmarshal(msg.Body, &body)
		if err != nil {
			return err
		}

		key := body["key"].(string)
		val := int(body["msg"].(float64))

		return replicator.Node.Reply(msg, map[string]any{
			"type":   "send_ok",
			"offset": replicator.Store.Set(key, val),
		})
	})

	replicator.Node.Handle("poll", func(msg maelstrom.Message) error {

		var body map[string]any

		err := json.Unmarshal(msg.Body, &body)
		if err != nil {
			return err
		}

		offsetMap := body["offsets"].(map[string]any)

		return replicator.Node.Reply(msg, map[string]any{
			"type": "poll_ok",
			"msgs": replicator.Store.Get(offsetMap),
		})
	})

	replicator.Node.Handle("commit_offsets", func(msg maelstrom.Message) error {

		var body map[string]any

		err := json.Unmarshal(msg.Body, &body)
		if err != nil {
			return err
		}

		offsetMap := body["offsets"].(map[string]any)
		replicator.Store.CommitOffsets(offsetMap)

		return replicator.Node.Reply(msg, map[string]any{
			"type": "commit_offsets_ok",
		})
	})

	replicator.Node.Handle("list_committed_offsets", func(msg maelstrom.Message) error {

		var body map[string]any

		err := json.Unmarshal(msg.Body, &body)
		if err != nil {
			return err
		}

		keys := body["keys"].([]any)
		return replicator.Node.Reply(msg, map[string]any{
			"type":    "list_committed_offsets_ok",
			"offsets": replicator.Store.GetCommitedOffsets(keys),
		})
	})

	if err := replicator.Node.Run(); err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}

}
