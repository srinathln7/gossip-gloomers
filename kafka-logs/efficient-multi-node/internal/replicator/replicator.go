package replicator

import (
	"context"
	"log"
	store "logs/internal/kafka-store"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const CLUSTER_LEADER = "cluster_leader"

type Replicator struct {
	mu sync.Mutex
	lv *maelstrom.KV

	Node   *maelstrom.Node
	Store  *store.KafkaStore
	Leader string
}

func NewReplicator() *Replicator {

	node := maelstrom.NewNode()
	return &Replicator{
		Node:  node,
		Store: store.NewKafkaStore(),
		lv:    maelstrom.NewLinKV(node),
	}
}

func (r *Replicator) RunLeaderElection() {

	r.mu.Lock()
	defer r.mu.Unlock()

	// If the node finds in its local state that the leader has not been assigned yet,
	if r.Leader == "" {

		// Attempt to acquire the leadership if no leader has been assigned yet
		err := r.lv.CompareAndSwap(context.Background(), CLUSTER_LEADER, r.Node.ID(), r.Node.ID(), true)

		if err != nil {

			log.Println("error: cas failure and now retrieving current leader from the cluster")
			// If CompareAndSwap fails, it means someone else may have acquired leadership
			// Retrieve the current leader from the key-value store
			leader, err := r.lv.Read(context.Background(), CLUSTER_LEADER)

			if err != nil {
				// If reading the leader from the key-value store also fails, handle the error
				log.Panic("error during leader election algorithm")
			} else {
				// If successful, update the local state with the retrieved leader information
				r.Leader = leader.(string)
			}
		}
	} else {
		r.Leader = r.Node.ID()
	}
}
