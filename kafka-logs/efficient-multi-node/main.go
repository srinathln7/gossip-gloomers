package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"

	"logs/internal/replicator"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {

	replicator := replicator.NewReplicator()

	replicator.Node.Handle("send", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		replicator.Mu.Lock()
		defer replicator.Mu.Unlock()

		//internal.CheckLeader(state, n, kv)

		replicator.RunLeaderElection()
		leader := replicator.GetLeader()

		if leader == replicator.Node.ID() {
			key := body["key"].(string)
			value := int(body["msg"].(float64))
			offset := replicator.Store.Set(key, value)

			body["type"] = "send_update"
			// update followers
			for _, f := range replicator.Node.NodeIDs() {
				if f == replicator.Node.ID() {
					continue
				}

				respL, err := replicator.Node.SyncRPC(context.Background(), f, body)
				if err != nil {
					return err
				}

				var bd map[string]any
				if err := json.Unmarshal(respL.Body, &bd); err != nil {
					return err
				}
				fo := bd["offset"]
				fv := int(fo.(float64))
				if fv != offset {
					return errors.New("consistency problem : send. leader:" + strconv.Itoa(offset) + " , follower " + strconv.Itoa(fv))
				}
			}

			resp := make(map[string]any)
			resp["type"] = "send_ok"
			resp["offset"] = offset
			return replicator.Node.Reply(msg, resp)
		} else {
			respL, err := replicator.Node.SyncRPC(context.Background(), leader, body)
			if err != nil {
				return err
			}
			return replicator.Node.Reply(msg, respL.Body)
		}
	})

	replicator.Node.Handle("poll", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		replicator.Mu.Lock()
		defer replicator.Mu.Unlock()

		data := body["offsets"].(map[string]any)
		res := replicator.Store.Get(data)

		resp := make(map[string]any)
		resp["type"] = "poll_ok"
		resp["msgs"] = res
		return replicator.Node.Reply(msg, resp)
	})

	replicator.Node.Handle("list_committed_offsets", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		replicator.Mu.Lock()
		defer replicator.Mu.Unlock()

		data := body["keys"].([]interface{})
		res := replicator.Store.GetCommitedOffsets(data)

		resp := make(map[string]any)
		resp["type"] = "list_committed_offsets_ok"
		resp["offsets"] = res
		return replicator.Node.Reply(msg, resp)

	})

	replicator.Node.Handle("commit_offsets", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		//internal.CheckLeader(state, n, kv)

		replicator.Mu.Lock()
		defer replicator.Mu.Unlock()

		replicator.RunLeaderElection()
		leader := replicator.GetLeader()

		if leader == replicator.Node.ID() {
			data := body["offsets"].(map[string]any)
			replicator.Store.CommitOffsets(data)

			body["type"] = "commit_offsets_update"
			// update followers
			for _, f := range replicator.Node.NodeIDs() {
				if f == replicator.Node.ID() {
					continue
				}

				_, err := replicator.Node.SyncRPC(context.Background(), f, body)
				if err != nil {
					return err
				}
			}
		} else {
			_, err := replicator.Node.SyncRPC(context.Background(), leader, body)
			if err != nil {
				return err
			}
		}

		return replicator.Node.Reply(msg, map[string]any{
			"type": "commit_offsets_ok",
		})
	})

	replicator.Node.Handle("send_update", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		key := body["key"].(string)
		value := int(body["msg"].(float64))
		offset := replicator.Store.Set(key, value)

		resp := make(map[string]any)
		resp["type"] = "send_update_ok"
		resp["offset"] = offset
		return replicator.Node.Reply(msg, resp)
	})

	replicator.Node.Handle("commit_offsets_update", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		data := body["offsets"].(map[string]any)
		replicator.Store.CommitOffsets(data)

		return replicator.Node.Reply(msg, map[string]any{
			"type": "commit_offsets_update_ok",
		})
	})

	if err := replicator.Node.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
