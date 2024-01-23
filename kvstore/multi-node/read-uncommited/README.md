# Read-Uncommitted

Refer to the challenge [here](https://fly.io/dist-sys/6b/)

### `main.go`

- The `main.go` file is the entry point of the application and defines the main functionality using the Maelstrom framework.
- It creates a Maelstrom node using `maelstrom.NewNode()` and a KVStore instance using `kvstore.NewKVStore()`.
- The `node.Handle` function is used to define handlers for different message types. Two message types are handled:
  - **"txn"**: Handles transactions by setting or reading values in the KVStore.
  - **"gossip"**: Handles gossip messages for synchronization between nodes.
- The `kvs.SetAndBroadcast` function is called to set a value in the KVStore and broadcast the change to other nodes.
- The `kvs.SyncLocal` function is called to synchronize local state with incoming gossip messages.
- The application then runs the Maelstrom node using `node.Run()`.

### `kvstore.go`

- The `kvstore.go` file contains the implementation of the key-value store (KVStore) and synchronization methods.
- It defines a `value` struct representing a key-value pair with a version.
- The `kvStore` struct represents the actual KVStore with methods for setting values, handling synchronization, and getting values.
- The `NewKVStore` function creates a new instance of the KVStore.
- The `SetAndBroadcast` method sets a value in the KVStore, increments the version, and broadcasts the change to other nodes using the `broadcast` function.
- The `SyncLocal` method synchronizes the local state with incoming gossip messages.
- The `set` method sets a value in the KVStore and returns the new version.
- The `Get` method retrieves a value from the KVStore based on the key


## Test

```
m test -w txn-rw-register --bin ./kvstore --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-uncommitted

m test -w txn-rw-register --bin ./kvstore --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-uncommitted --availability total --nemesis partition
```