# Grow-Only Counter 

A classic application of the Conflict-Free Replicated Data Types (CRDTs) particularly state-based CRDTS (aka. cv-CRDTS - Convergent CRDTS). This  application demonstrates a distributed counter implementation using Maelstrom and a gossip protocol. The counter supports addition operations and can handle concurrent updates across multiple nodes in a distributed network.

## Overview

The counter application consists of a Maelstrom node with three main functionalities: addition (`add`), gossip-based synchronization (`gossip`), and reading the current counter value (`read`). The distributed nature of the application is facilitated through a gossip protocol that allows nodes to share and synchronize counter values.

## Components

### 1. `main.go`

The `main.go` file serves as the entry point of the counter application. It initializes a Maelstrom node, a global counter (`Counter`), a global key (`final_result`), and a gossip protocol. The Maelstrom node handles different types of messages, including addition, gossip, and read operations. Additionally, the application uses a sequential consistent key-value store (`maelstrom.SeqKV`) to manage the distributed counter.

**Key Highlights:**
- Addition (`add`) operation: Increments the counter value based on the received delta and performs a compare-and-swap operation for atomic updates in a concurrent environment.
- Gossip (`gossip`) operation: Synchronizes counter values received from peer nodes using the gossip protocol.
- Read (`read`) operation: Retrieves and returns the current counter value, considering updates from both local and distributed sources.

### 2. `p2p/gossip.go`

The `p2p/gossip.go` file contains the implementation of the gossip protocol. It defines a `Gossip` struct with methods for starting and stopping the gossip mechanism.

**Key Highlights:**
- `NewGossip`: Initializes a new gossip instance with a given interval.
- `Start`: Initiates the gossip mechanism to share counter values among nodes.
- `Stop`: Stops the gossip mechanism.


## Test

```
m test -w g-counter --bin ./counter --node-count 3 --rate 100 --time-limit 20 --nemesis partition
```

## Conclusion

The counter application showcases a distributed counter with concurrency control using Maelstrom and a gossip protocol. It provides insights into handling addition operations, gossip-based synchronization, and reading counter values in a distributed setting. The modular design allows for easy extension and customization to fit various distributed system scenarios.
