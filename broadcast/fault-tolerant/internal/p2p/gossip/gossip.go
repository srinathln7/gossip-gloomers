package gossip

import (
	"broadcast/internal/cache"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Gossip struct {
	ticker *time.Ticker
	done   chan bool
}

type GossipMsg struct {
	Type string `json:"type"`
	Msg  int    `json:"msg"`
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

func (gossip *Gossip) Start(kv *cache.KVStore, node *maelstrom.Node, maxNodes int) {

	go func() {
		for {
			select {
			case <-gossip.done:
				gossip.Stop()

			case <-gossip.ticker.C:
				err := doGossip(kv, node, maxNodes)
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

func doGossip(kv *cache.KVStore, node *maelstrom.Node, maxNodes int) error {

	peers := node.NodeIDs()
	data := kv.Get()

	body := &GossipMsgs{
		Type: "gossip",
		Msgs: data,
	}

	for _, peer := range peers {

		if peer == node.ID() {
			continue
		}

		go func(peer string) {
			for {
				err := node.Send(peer, body)
				if err == nil {
					break
				}

				node.Run()
			}

		}(peer)

	}

	return nil
}
