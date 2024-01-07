# Unique ID Generator

```


### Usage

1. Import the required packages:

   ```
   
   package main

   import (
   	"encoding/json"
   	"log"

   	"github.com/google/uuid"
   	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
   )
   ```

2. Create a Maelstrom node:

   ```go
   func main() {
   	n := maelstrom.NewNode()
   ```
   
3. Define a message handler to generate UUIDs:

   ```go
   	n.Handle("generate", func(msg maelstrom.Message) error {
   		var body map[string]interface{}
   
   		if err := json.Unmarshal(msg.Body, &body); err != nil {
   			return err
   		}
   
   		body["type"] = "generate_ok"
   		body["id"] = uuid.New().String()
   
   		return n.Reply(msg, body)
   	})
   ```

4. Run the Maelstrom node:

   ```go
   	if err := n.Run(); err != nil {
   		log.Fatal(err)
   	}
   ```

## Explanation

- The code creates a Maelstrom node and sets up a message handler for the "generate" message type.
- The handler unmarshals the message body into a loosely-typed map, updates the message type, generates a new UUID, and replies with the modified message.
- The `Run()` method continuously reads messages from STDIN, triggers the associated handler, and returns an error if no handler exists for a message type.


## Test

The test will run a 3-node cluster for 30 seconds and request new IDs at the rate of 1000 requests per second. It checks for total availability and will induce network partitions during the test. It will also verify that all IDs are unique.

```
m test -w unique-ids --bin ./uid --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
```