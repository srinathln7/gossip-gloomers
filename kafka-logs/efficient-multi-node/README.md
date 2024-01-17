# Distributed Log Replicator

## Components

### `main.go`

The `main.go` file contains the main entry point for each node in the distributed system. It utilizes the Maelstrom framework for distributed systems and orchestrates the interaction between nodes. Key functionalities include:

- **Leader Election:** Implements the leader election algorithm to determine the leader node responsible for handling write requests.

- **Write Propagation:** Propagates write updates to followers, ensuring consistency across the distributed system.

- **Commit Propagation:** Propagates commit updates to followers, ensuring that committed offsets are consistent.

- **Polling:** Handles polling requests to retrieve log entries and respond with the latest committed offsets.

- **Usage:**
  - Compile and run the code on each node in the distributed system.

### `replicator.go`

The `replicator.go` file defines the `Replicator` struct, responsible for managing the state and operations related to log replication. Key components and functionalities include:

- **Leader Election:** Implements the leader election algorithm based on the Compare-and-Swap (CAS) operation using Maelstrom's KV store.

- **Get Leader:** Retrieves the current leader, initiating leader election if no leader is set.

- **Node Initialization:** Initializes the Maelstrom node, Kafka store, and linearizable key-value store.

- **Run Leader Election:** Executes the leader election algorithm to acquire or retrieve the current leader.

- **Usage:**
  - Create an instance of `Replicator` for each node in the distributed system.
  - Utilize the `RunLeaderElection` and `GetLeader` methods to manage leadership.

## Test

```
m test -w kafka --bin ./logs --node-count 2 --concurrency 2n --time-limit 20 --rate 1000
```