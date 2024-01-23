# Read-committed

Refer to the challenge [here](https://fly.io/dist-sys/6c/). The approach is pretty similar to the approach we took in 6b. The main difference is that we start broadcasting only when all the transactions are over unlike in read-uncommitted state where are broadcast immediately after we write a value to the store.


### `main.go`

This file defines the main functionality of the distributed system. It uses the Maelstrom framework to create a node that handles transactions, gossip communication, and error handling. Transactions involve read (`r`) and write (`w`) operations on a key-value store (`kvstore`). The system broadcasts updates (`gossip`) to synchronize the key-value store across all nodes. It also handles conflicts and aborts transactions if necessary.

### `kvstore.go`

This file defines the key-value store (`kvstore`) with functions for setting, getting, and synchronizing values. The store maintains a local snapshot and supports synchronization through gossip communication. The `Broadcast` function sends gossip messages to other nodes for synchronization, ensuring consistency across the distributed system. Together, these files create a distributed system where nodes can perform transactions on a shared key-value store, and updates are propagated through gossip communication to maintain consistency among the nodes.

## Test

```
m test -w txn-rw-register --bin ./kvstore --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-committed --availability total --nemesis partition
```