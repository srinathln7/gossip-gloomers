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

		// Run the leader election algorithm before sending the requests
		replicator.RunLeaderElection()
		leader := replicator.GetLeader()
		if len(leader) == 0 {
			log.Panic("leader found to be empty")
		}

		log.Printf("[main] node %s log reads leader set to %s", replicator.Node.ID(), leader)

		// Accept write request only if the node is the leader
		if replicator.Node.ID() == leader {
			key := body["key"].(string)
			val := int(body["msg"].(float64))
			offset := replicator.Store.Set(key, val)

			// Propogate the updates to the followers
			body["type"] = "propogate_send_update"

			nodeList := replicator.Node.NodeIDs()
			for _, follower := range nodeList {
				if follower == replicator.Node.ID() {
					continue
				}

				log.Printf("propogating updates to the follower: %s", follower)
				resp, err := replicator.Node.SyncRPC(context.Background(), follower, body)
				if err != nil {
					log.Panicf("error %s while sending update to follower: %s", err.Error(), follower)
					return err
				}

				var respMap map[string]any
				if err := json.Unmarshal(resp.Body, &respMap); err != nil {
					return err
				}

				recvOffSet := int(respMap["offset"].(float64))
				log.Println("checking for db consistency")
				if recvOffSet != offset {
					log.Panicf("consistency error: leader sent %v while follower received %v", offset, recvOffSet)
				}
			}

			return replicator.Node.Reply(msg, map[string]any{
				"type":   "send_ok",
				"offset": offset,
			})
		} else {

			resp, err := replicator.Node.SyncRPC(context.Background(), leader, body)
			if err != nil {
				return err
			}
			return replicator.Node.Reply(msg, resp.Body)
		}

	})

	replicator.Node.Handle("propogate_send_update", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		key := body["key"].(string)
		val := int(body["msg"].(float64))
		offset := replicator.Store.Set(key, val)

		return replicator.Node.Reply(msg, map[string]any{
			"type":   "propogate_send_update_ok",
			"offset": offset,
		})
	})

	replicator.Node.Handle("poll", func(msg maelstrom.Message) error {

		var body map[string]any
		err := json.Unmarshal(msg.Body, &body)
		if err != nil {
			return err
		}

		replicator.Mu.RLock()
		defer replicator.Mu.RUnlock()

		data := body["offsets"].(map[string]any)
		res := replicator.Store.Get(data)

		resp := make(map[string]any)
		resp["type"] = "poll_ok"
		resp["msgs"] = res
		return replicator.Node.Reply(msg, resp)
	})

	replicator.Node.Handle("commit_offsets", func(msg maelstrom.Message) error {

		var body map[string]any
		err := json.Unmarshal(msg.Body, &body)
		if err != nil {
			return err
		}

		// Run the leader election algorithm before sending the requests
		replicator.RunLeaderElection()
		leader := replicator.GetLeader()

		// Accept write request only if the node is the leader
		if replicator.Node.ID() == leader {
			offsetMap := body["offsets"].(map[string]any)
			replicator.Store.CommitOffsets(offsetMap)

			// Propogate the commit update to the followers
			body["type"] = "propogate_commit_update"

			nodeList := replicator.Node.NodeIDs()
			for _, follower := range nodeList {

				if follower == replicator.Node.ID() {
					continue
				}

				_, err := replicator.Node.SyncRPC(context.Background(), follower, body)
				if err != nil {
					return err
				}
			}
		} else {
			// Delegate all writes to the `Leader`
			log.Printf("redirecting all write requests to leader %s", leader)
			_, err := replicator.Node.SyncRPC(context.Background(), leader, body)
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

		replicator.Mu.RLock()
		defer replicator.Mu.RUnlock()

		keys := body["keys"].([]any)
		return replicator.Node.Reply(msg, map[string]any{
			"type":    "list_committed_offsets_ok",
			"offsets": replicator.Store.GetCommitedOffsets(keys),
		})
	})

	replicator.Node.Handle("propogate_commit_update", func(msg maelstrom.Message) error {
		var body map[string]any
		err := json.Unmarshal(msg.Body, &body)
		if err != nil {
			return err
		}

		offsetMap := body["offsets"].(map[string]any)
		replicator.Store.CommitOffsets(offsetMap)

		return replicator.Node.Reply(msg, map[string]any{
			"type": "propogate_commit_update_ok",
		})
	})

	if err := replicator.Node.Run(); err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}

}
