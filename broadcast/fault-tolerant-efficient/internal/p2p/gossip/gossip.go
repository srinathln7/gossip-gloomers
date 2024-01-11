package gossip

import (
	"broadcast/internal/node"
	"time"
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

type GossipMsgs struct {
	Type string `json:"type"`
	Msgs []int  `json:"msgs"`
}

func NewGossip(d time.Duration) *Gossip {
	return &Gossip{
		ticker: time.NewTicker(d),
		done:   make(chan bool),
	}
}

func (gossip *Gossip) Start(node *node.Node) {

	go func() {
		for {
			select {
			case <-gossip.done:
				gossip.Stop()

			case <-gossip.ticker.C:
				err := doGossip(node)
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

func doGossip(node *node.Node) error {

	peers := node.N.NodeIDs()
	data := node.Store.Get()

	body := &GossipMsgs{
		Type: "gossip",
		Msgs: data,
	}

	for _, peer := range peers {

		if peer == node.N.ID() {
			continue
		}

		go func(peer string) {
			for {
				err := node.N.Send(peer, body)
				if err == nil {
					break
				}
			}

		}(peer)

	}

	return nil
}
