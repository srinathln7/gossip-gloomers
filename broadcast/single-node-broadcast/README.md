# Single Broadcast and Read

## Usage

1. Import the required packages:

   ```go
   package main

   import (
   	"encoding/json"
   	"log"

   	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
   )
   ```

2. Create a Maelstrom node and initialize a cache for storing broadcasted values:

   ```go
   var cache map[string][]int = make(map[string][]int)
   ```

3. Define message handlers for "broadcast," "read," and "topology" message types:

   ```go
   n.Handle("broadcast", func(msg maelstrom.Message) error {
   	var body map[string]interface{}

   	if err := json.Unmarshal(msg.Body, &body); err != nil {
   		return err
   	}

   	// Store the broadcasted values locally in each node against its unique identifier in-memory
   	cache[n.ID()] = append(cache[n.ID()], int(body["message"].(float64)))

   	return n.Reply(msg, map[string]interface{}{
   		"type": "broadcast_ok",
   	})
   })

   n.Handle("read", func(msg maelstrom.Message) error {
   	var res map[string]interface{}

   	if err := json.Unmarshal(msg.Body, &res); err != nil {
   		return err
   	}

   	res["type"] = "read_ok"
   	res["messages"] = cache[n.ID()]

   	return n.Reply(msg, res)
   })

   n.Handle("topology", func(msg maelstrom.Message) error {
   	return n.Reply(msg, map[string]interface{}{
   		"type": "topology_ok",
   	})
   })
   ```

4. Run the Maelstrom node:

   ```go
   if err := n.Run(); err != nil {
   	log.Fatal(err)
   }
   ```

## Explanation and Comments in the Code

- The code initializes a cache to store broadcasted values locally for each node.
- The "broadcast" handler stores broadcasted values in the cache against the unique identifier of each node.
- The "read" handler retrieves stored values from the cache for the requesting node.
- The "topology" handler replies with a simple acknowledgment message.
- The `Run()` method continuously reads messages from STDIN, triggers the associated handler, and returns an error if no handler exists for a message type.


## Test

```
m test -w broadcast --bin broadcast --node-count 1 --time-limit 20 --rate 10
```