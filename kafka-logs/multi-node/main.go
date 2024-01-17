package main

import (
	"context"
	"encoding/json"
	"log"
	"logs/internal/replicator"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {

	// Each node will get its own local kafka store and linearizable KV
	replicator := replicator.NewReplicator()

	replicator.Node.Handle("send", func(msg maelstrom.Message) error {
		var body map[string]any

		err := json.Unmarshal(msg.Body, &body)
		if err != nil {
			return err
		}

		key := body["key"].(string)
		val := int(body["msg"].(float64))

		// Run the leader election algorithm before sending the requests
		replicator.RunLeaderElection()

		var result int
		// Accept write request only if the node is the leader
		if replicator.Node.ID() == replicator.Leader {
			result = replicator.Store.Set(key, val)
		} else {

			// Delegate all writes to the `Leader`. Just `Send` wont work since writes and sync access to STDOUT.
			// Only `syncRPC` works because we need to wait until our leader responds to our write requests to proceed further.
			// Looking into the implementation of both `send` and `syncRPC` we see the latter is a wrapper around the former
			// with repeated retries until the context is either cancelled or the messgae is
			//err := replicator.Node.Send(replicator.Leader, body)
			_, err := replicator.Node.SyncRPC(context.Background(), replicator.Leader, body)
			if err != nil {
				return err
			}
		}

		return replicator.Node.Reply(msg, map[string]any{
			"type":   "send_ok",
			"offset": result,
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

		// Run the leader election algorithm before sending the requests
		replicator.RunLeaderElection()

		// Accept write request only if the node is the leader
		if replicator.Node.ID() == replicator.Leader {
			replicator.Store.CommitOffsets(offsetMap)
		} else {
			// Delegate all writes to the `Leader`
			_, err := replicator.Node.SyncRPC(context.Background(), replicator.Leader, body)
			if err != nil {
				return err
			}
		}

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
