# Grow-Only Counter 

A classic application of the Conflict-Free Replicated Data Types (CRDTs) particularly state-based CRDTS (aka. cv-CRDTS - Convergent CRDTS). This  application demonstrates a distributed counter implementation using Maelstrom and a gossip protocol. The counter supports addition operations and can handle concurrent updates across multiple nodes in a distributed network. Recommend to watch this [video](https://www.youtube.com/watch?v=Fm8iUFM2iWU&t=1s) and read [this](https://jepsen.io/consistency) on consistency models before proceding further.

## Overview

The counter application consists of a Maelstrom node with two main functionalities: addition (add) and reading (read) operations. The application maintains a global counter (Counter) and utilizes a sequential consistent key-value store (maelstrom.SeqKV) to ensure atomic and consistent updates in a distributed environment.

### Components
1. main.go
The main.go file serves as the entry point of the counter application. It initializes a Maelstrom node, a global counter (Counter), a global key (final_result), and a sequential consistent key-value store (maelstrom.SeqKV). The Maelstrom node handles different types of messages, including addition and reading operations.

Key Highlights:

- Addition (add) operation: Increments the counter value based on the received delta and performs a compare-and-swap operation for atomic updates in a concurrent environment.
- Read (read) operation: Retrieves and returns the current counter value from the distributed key-value store.

2. internal/counter/counter.go
The counter.go file defines the Counter struct, representing the global counter. It includes a read-write mutex to manage concurrent access to the counter value.

Key Highlights:

- Represents the global counter.
- Uses a read-write mutex for concurrent access to the counter value.


## Test

```
m test -w g-counter --bin ./counter --node-count 3 --rate 100 --time-limit 20 --nemesis partition
```

## Conclusion

The counter application showcases a distributed counter with concurrency control using Maelstrom and a sequential consistent key-value store. It provides insights into handling addition operations and reading counter values in a distributed setting. The modular design allows for easy extension and customization to fit various distributed system scenarios.