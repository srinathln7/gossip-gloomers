package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {

	// Create a maelstrom node to begin with

	n := maelstrom.NewNode()

	n.Handle("echo", func(msg maelstrom.Message) error {

		// In this handler, we’re unmarshaling to a generic map since we simply want to echo back the same message we received.
		// Unmarshal the maelstrom msg body as an loosely-typed map
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type to return back
		body["type"] = "echo_ok"

		// Echo the original msg back with the updated msg type
		// The Reply() method will automatically set the source and destination fields in the return message and it
		// will associate the message as a reply to the original one received

		return n.Reply(msg, body)
	})

	// Finally, we’ll delegate execution to the Node by calling its Run() method.
	// This method continuously reads messages from STDIN and fires off a goroutine for each one to the associated handler.
	// If no handler exists for a message type, Run() will return an error.

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}

}
