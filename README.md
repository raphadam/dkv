![Logo](docs/dkvlogo.png)

## Introduction

Simple yet powerful distributed key-value store written in Go. It uses the Raft consensus algorithm to ensure data consistency across nodes, Serf for node discovery and Prometheus for logging and monitoring. Exposes a RESTful JSON API for client interactions.

## Features

- **Distributed Consensus:** Utilizes the Raft algorithm to ensure consistency and reliability across nodes.
- **Dynamic Node Discovery:** Implements Serf for efficient peer-to-peer node discovery and communication.
- **Concurrency:** Nodes handle operations concurrently, ensuring scalability and performance.
- **RESTful API:** Exposes a simple and intuitive REST/JSON API for client interactions.
- **Monitoring and Logging:** Integrated with Prometheus for real-time monitoring and logging of metrics.

## Usage

### Dynamic Node Management

The cluster can dynamically scale, meaning nodes can join or leave at any time. This flexibility allows for easy scaling and fault tolerance.

- **Read Requests:** Clients can send GET requests to any node in the cluster, and it will respond with the correct value, even if it’s not the leader.
- **Write Requests:** Write operations must be directed to the leader node. If a client attempts to write to a follower node, it will be receive the address of the current leader.

### API Reference

- **Endpoint:** /
- **Method:** POST

**Request Body:**

```json
{
    "cmd": "set|get|delete",
    "key": "your key",
    "val": "your val"
}
```

**Response:**

```json
{
    "val": "val corresponding to the key",
}
```

## Future Work

This project is still under development with several features planned for future releases:

- **Sharding:** Implementing sharding to distribute data across multiple nodes, enhancing scalability and performance.
- **Consistent Hashing:** Introducing consistent hashing to ensure an even distribution of data and reduce the impact of node changes.
- **Improved Load Balancing:** Developing strategies to balance read and write loads more effectively across the cluster.

## Installation

### Prerequisites

- Go 1.22+
- Git

### Steps

Clone the repo:

    git cline https://github.com/raphadam/dkv.git
    cd dkv

Run the command:

    go run ./cmd


