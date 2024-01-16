package replicator

import (
	store "logs/internal/kafka-store"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Replicator struct {
	Node  *maelstrom.Node
	Store *store.KafkaStore
}

func NewReplicator() *Replicator {
	return &Replicator{
		Node:  maelstrom.NewNode(),
		Store: store.NewKafkaStore(),
	}
}
