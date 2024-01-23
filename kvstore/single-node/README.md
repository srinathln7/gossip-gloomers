# KVStore 

## Overview

The code consists of a single Go application (`main.go`) that utilizes the Maelstrom framework for distributed systems. It creates a Maelstrom node and a KVStore instance. The node handles transactions and interacts with the KVStore based on the received transaction requests.

## Components

### 1. Maelstrom Node

The Maelstrom node is created using the `maelstrom.NewNode()` function. It is responsible for handling incoming messages and orchestrating the execution of distributed transactions.

### 2. KVStore

The KVStore (`kvstore/internal/kvs`) is a basic in-memory key-value store with the following features:

- **Set Operation (`w`)**: Updates the value associated with a given key.
- **Get Operation (`r`)**: Retrieves the value associated with a given key.

### 3. Transaction Handling

The node handles transactions of the type `"txn"`. The transactions are received as JSON messages, and the body is unmarshaled to extract the transaction details. The application supports two types of transaction requests:

- **Write Operation (`w`)**: Updates the value for a specified key in the KVStore.
- **Read Operation (`r`)**: Retrieves the value for a specified key from the KVStore.

The responses to the transaction requests include the transaction type, key, and the result of the operation. All transaction responses are sent back with the type `"txn_ok"`.

## Test

```
m test -w txn-rw-register --bin ./kvstore --node-count 1 --time-limit 20 --rate 1000 --concurrency 2n --consistency-models read-uncommitted --availability total
```
