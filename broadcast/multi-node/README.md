# Multi-Node Broadcast

Refer to the challenge [here](https://fly.io/dist-sys/3b/)

## Overview

The main components of the system include:

- **Node Struct**: Represents a node in the distributed system, encapsulating the Maelstrom node, waitgroups for go routine co-ordination, and a map to track received broadcast messages.

- **TopologyMsg Struct**: Defines the structure for the JSON message used to communicate network topology.

- **NwTopology Map**: Stores the network topology shared across all nodes.


## Components

### Node Struct

The `Node` struct represents a node in the distributed system and includes the following key components:


### TopologyMsg Struct

The `TopologyMsg` struct defines the structure of the JSON message used to communicate network topology. It includes a `Type` field and a `Topology` field, which represents the network topology as a map.

### NwTopology Map

The `NwTopology` map stores the network topology shared across all nodes.

## Message Handlers

- **hBroadcast**: Handles the "broadcast" message type. It stores locally broadcasted values and broadcasts them to neighbors.
  
- **hRead**: Handles the "read" message type. Responds with a list of received messages.

- **hTopology**: Handles the "topology" message type. Stores the network topology.

## Test

```
m test -w broadcast --bin ./broadcast --node-count 5 --time-limit 20 --rate 10
```