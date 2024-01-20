package replicator

import (
	"context"
	"fmt"
	store "logs/internal/kafka-store"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const CLUSTER_LEADER = "cluster_leader"

type Replicator struct {
	mu     sync.RWMutex
	lv     *maelstrom.KV
	leader string

	Node  *maelstrom.Node
	Store *store.KafkaStore
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
	// then contend for the election to become the leader
	if r.leader == "" {
		// Attempt to acquire the leadership if no leader has been assigned yet
		err := r.lv.CompareAndSwap(context.Background(), CLUSTER_LEADER, r.Node.ID(), r.Node.ID(), true)

		if err != nil {

			// If CompareAndSwap fails, it means someone else may have acquired leadership
			// Retrieve the current leader from the key-value store
			leader, err0 := r.lv.Read(context.Background(), CLUSTER_LEADER)
			if err0 != nil {
				_ = fmt.Errorf(err0.Error())
			} else {
				// If successful, update the local state with the retrieved leader information
				r.leader = leader.(string)
			}
		} else {
			r.leader = r.Node.ID()
		}
	}
}

func (r *Replicator) GetLeader() string {

	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.leader
}
