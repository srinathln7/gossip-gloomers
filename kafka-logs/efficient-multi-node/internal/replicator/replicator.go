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
	lv     *maelstrom.KV
	leader string

	Mu    sync.RWMutex
	Node  *maelstrom.Node
	Store *store.KafkaStore
}

func NewReplicator() *Replicator {

	node := maelstrom.NewNode()
	return &Replicator{
		Node:  node,
		Store: store.NewKafkaStore(),
		lv:    maelstrom.NewLinKV(node),
		Mu:    sync.RWMutex{},
	}
}

func (r *Replicator) RunLeaderElection() {

	// r.mu.Lock()
	// defer r.mu.Unlock()

	// If the node finds in its local state that the leader has not been assigned yet,
	if r.leader == "" {

		// Attempt to acquire the leadership if no leader has been assigned yet
		err := r.lv.CompareAndSwap(context.Background(), CLUSTER_LEADER, r.Node.ID(), r.Node.ID(), true)
		if err != nil {

			// If CompareAndSwap fails, it means someone else may have acquired leadership
			// Retrieve the current leader from the key-value store
			log.Println("error: cas failure and now retrieving current leader from the cluster")
			leader, err := r.lv.Read(context.Background(), CLUSTER_LEADER)
			if err != nil {
				// If reading the leader from the key-value store also fails, handle the error
				log.Panic("error during leader election algorithm")
			} else {
				// If successful, update the local state with the retrieved leader information
				r.leader = leader.(string)
			}
		}
	} else {
		r.leader = r.Node.ID()
	}
}

func (r *Replicator) GetLeader() string {
	// r.mu.RLock()
	// defer r.mu.RUnlock()
	return r.leader
}
