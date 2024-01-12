package gossip

import (
	"context"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// NWTopology  Structure of the network topology.
type NWTopology struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

type Gossip struct {
	ticker *time.Ticker
	done   chan bool
}

type GossipMsg struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

func NewGossip(d time.Duration) *Gossip {
	return &Gossip{
		ticker: time.NewTicker(d),
		done:   make(chan bool),
	}
}

func (gossip *Gossip) Start(node *maelstrom.Node, kv *maelstrom.KV) {

	go func() {
		for {
			select {
			case <-gossip.done:
				gossip.Stop()

			case <-gossip.ticker.C:
				err := doGossip(node, kv)
				if err != nil {
					return
				}
			}
		}

	}()

}

func (gossip *Gossip) Stop() {
	gossip.done <- true
}

func doGossip(node *maelstrom.Node, kv *maelstrom.KV) error {

	peers := node.NodeIDs()
	result, err := kv.ReadInt(context.Background(), "final_result")
	if err != nil {
		return err
	}

	body := &GossipMsg{
		Type:  "gossip",
		Value: result,
	}

	for _, peer := range peers {

		if peer == node.ID() {
			continue
		}

		go func(peer string) {
			for {
				err := node.Send(peer, body)
				//_, err := node.SyncRPC(context.Background(), peer, body)
				if err == nil {
					break
				}
			}

		}(peer)

	}

	return nil
}
