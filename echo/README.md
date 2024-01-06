# Echo Challenge

## Code

1. Import the required packages:

   ```go
   package main

   import (
   	"encoding/json"
   	"log"

   	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
   )
   ```

2. Create a Maelstrom node:

   ```go
   func main() {
   	n := maelstrom.NewNode()
   ```
   
3. Define a message handler:

   ```go
   	n.Handle("echo", func(msg maelstrom.Message) error {
   		var body map[string]interface{}
   
   		if err := json.Unmarshal(msg.Body, &body); err != nil {
   			return err
   		}
   
   		body["type"] = "echo_ok"
   
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

- The code creates a Maelstrom node and sets up a message handler for the "echo" message type.
- The handler unmarshals the message body into a loosely-typed map, updates the message type, and replies with the modified message.
- The `Run()` method continuously reads messages from STDIN, triggers the associated handler, and returns an error if no handler exists for a message type.


## Test

```
m test -w echo --bin ./echo --node-count 1 --time-limit 10

```
