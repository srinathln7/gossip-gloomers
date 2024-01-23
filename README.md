# gossip-gloomers

This repository presents the solutions to a series of challenges in distributed systems, originating from [Fly.io](https://fly.io/dist-sys/) and devised by [Kyle Kingsbury](https://aphyr.com/about), the creator of [Jepsen](https://jepsen.io/). These challenges leverage Maelstrom, a messaging framework, and are built on top of it. Maelstrom, itself constructed on Jepsen, serves as a platform for communication between distributed systems. With Maelstrom, you can construct a "node" within your distributed system, and the framework seamlessly manages the routing of messages between these nodes. This capability allows Maelstrom to introduce failures and conduct verification checks based on the specified consistency guarantees for each challenge.

## Getting Started

### Prerequisites

- Go installed on your machine

### Installation

Run the following command in your terminal

```
go get github.com/jepsen-io/maelstrom/demo/go
go install .
```

For more information about Maelstrom, refer to the [official documentation](https://github.com/jepsen-io/maelstrom). Once you download the maelstrom binaries, set the following in your `~/.bashrc` file before getting started with the challenge

```
alias m=maelstrom
```
