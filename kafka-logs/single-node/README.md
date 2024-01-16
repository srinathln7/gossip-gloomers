# Replicator

This application demonstrates a simple message replication system using Maelstrom and Kafka-like stores. The system includes three main components: `main.go`, `kafka-store.go`, and `replicator.go`.

## Overview

The Replicator application replicates messages, managing offsets and committed offsets in a Kafka-like store. It exposes operations for sending messages, polling messages, committing offsets, and listing committed offsets.

## Components

### 1. `main.go`

`main.go` is the entry point of the application. It initializes a Replicator, handles various message types, and interacts with the Kafka-like store.

**Highlighted Features:**
- Handles "send," "poll," "commit_offsets," and "list_committed_offsets" message types.
- Uses a Kafka-like store for storing and retrieving data.
- Responds to message requests and updates the Kafka store accordingly.

### 2. `kafka-store.go`

`kafka-store.go` defines the KafkaStore, a simple in-memory store for handling data and offsets. It supports setting data, getting data based on offsets, committing offsets, and listing committed offsets.

**Highlighted Features:**
- Implements an in-memory Kafka-like store.
- Provides methods for setting data, getting data based on offsets, committing offsets, and listing committed offsets.

### 3. `replicator.go`

`replicator.go` initializes the Replicator structure, which includes a Maelstrom node and a Kafka-like store. It provides a simple way to create a new Replicator instance.

**Highlighted Features:**
- Initializes a Maelstrom node and a Kafka-like store for the Replicator.

## Test

```
m test -w kafka --bin ./logs --node-count 1 --concurrency 2n --time-limit 20 --rate 1000
```

## Conclusion

The Replicator application showcases a basic message replication system with Kafka-like store functionalities. It leverages Maelstrom for message handling and includes a simple in-memory store for data storage and retrieval. The modular design allows for easy extension and customization to fit various distributed system scenarios.


