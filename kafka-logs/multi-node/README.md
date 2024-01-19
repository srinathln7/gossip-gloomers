# Distributed Log Replicator

This repository contains code for a distributed log replicator that utilizes a linearizable key-value store for leader election and coordination. Note in this example,
all replicas redirect all read/write requests to the leader.

## Main Application (main.go)

The `main.go` file represents the entry point of the distributed log replicator application. It is structured as follows:

- **Node Initialization:**
  The application initializes a local instance of the replicator, including a Node, a Kafka store, and a linearizable key-value store.

- **Message Handling:**
  The Node handles various types of messages, including "send," "poll," "commit_offsets," and "list_committed_offsets." These messages are part of the communication protocol for replicating log entries.

- **Leader Election:**
  Before processing requests, the application runs a leader election algorithm to determine the leader node in the cluster.

- **Request Processing:**
  Based on leadership status, the node either processes write requests locally or forwards them to the leader. Read requests are handled by retrieving data from the local store.

- **Application Execution:**
  The application runs the Node, handling incoming messages and coordinating with other nodes in the cluster.

### Replicator Implementation (replicator.go)

The `replicator.go` file contains the implementation of the distributed log replicator. Key functionalities include:

- **Initialization:**
  The creation of a new Replicator instance, initializing the Node, Kafka store, and linearizable key-value store.

- **Leader Election:**
  The `RunLeaderElection` function implements the leader election algorithm. It attempts to acquire leadership by performing a CompareAndSwap operation on the "cluster_leader" key in the key-value store.

- **Local State Management:**
  The `Replicator` struct maintains essential information, including the Node, Kafka store, and the current leader. A Mutex (`mu`) is used for safe concurrent access to shared resources.

### Kafka Store Implementation (store.go)

The `store.go` file represents the implementation of the Kafka store used by the distributed log replicator. It is mentioned as an internal detail and presumed to be similar to a single-node Kafka log. 

## Test

```
m test -w kafka --bin ./logs --node-count 2 --concurrency 2n --time-limit 20 --rate 1000
```