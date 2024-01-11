# Fault-Tolerant Broadcast System

This repository contains a fault-tolerant broadcast system built on top of [Maelstrom](https://github.com/jepsen-io/maelstrom). The system utilizes a simple yet effective approach for communication between nodes in a distributed network. This README provides an overview of the components and functionalities of the system.

## Challenge

Refer to the challenge details [here](https://fly.io/dist-sys/3c/).

## Broadcast System

The broadcast system consists of three main components: `main.go`, `node.go`, and `gossip.go`. Each component plays a specific role in the overall functionality of the system.

### `main.go`

The `main.go` file serves as the entry point of the system. It performs the following tasks:

- Creates a new `Node` using the Maelstrom library.
- Defines handlers for different message types, including "broadcast," "gossip," "read," and "topology."
- Runs the Maelstrom node and initializes a gossip mechanism for broadcasting messages.

**Highlighted Features:**
- Handles multiple message types using Maelstrom.
- Utilizes a KV store for each node to manage locally stored data.
- Implements a gossip mechanism for broadcasting messages to peers.

### `node.go`

The `node.go` file defines the `Node` struct, representing a node in the distributed system. Key features include:

- Represents a node in the distributed system, maintaining information such as the Maelstrom node, KV store, and topology.
- Provides methods for setting and getting the network topology.
- Utilizes the Maelstrom library for communication in a distributed system.

**Highlighted Features:**
- Node structure with its own KV store and network topology.
- Methods for managing and retrieving network topology.

### `gossip.go`

The `gossip.go` file contains the implementation of the gossip mechanism used for broadcasting messages among nodes. Important features include:

- Implements a gossip mechanism for broadcasting messages to peers.
- Uses a ticker to initiate gossip periodically.
- Starts and stops the gossip mechanism, allowing flexibility in its usage.

**Highlighted Features:**
- Gossip mechanism for broadcasting messages to peers.
- Ticker-based initiation for periodic gossip.
- Start and stop methods for managing the gossip mechanism.

## System Operation

1. **Initialization:**
   - Each node is initialized with its Maelstrom node, KV store, and topology.
   - Handlers for different message types are defined, enabling the node to process broadcast, gossip, read, and topology messages.

2. **Message Broadcasting:**
   - The system uses the Maelstrom library to handle different message types.
   - The "broadcast" message type triggers the broadcast of messages to neighbors through the gossip mechanism.

3. **Gossip Mechanism:**
   - The gossip mechanism periodically sends messages to random peers in the network, facilitating message propagation.

4. **Topology Management:**
   - The system allows nodes to set and get their network topology, providing flexibility in managing communication links.

## Test

To run the test, execute the following command:

```bash
m test -w broadcast --bin ./broadcast --node-count 25 --time-limit 20 --rate 100 --latency 100 --nemesis partition

m test -w broadcast --bin ./broadcast --node-count 5 --time-limit 20 --rate 10 --nemesis partition
```

## Conclusion

The fault-tolerant broadcast system offers a straightforward and efficient way for nodes in a distributed network to communicate and share information. The combination of Maelstrom and a gossip protocol provides a modular and customizable design suitable for various distributed system scenarios.