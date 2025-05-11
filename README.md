# Dijkstra Demo

This project provides an implementation of **Dijkstra's shortest path algorithm** and **Depth-First Search (DFS)** for network routing. The implementation includes support for node blacklisting and is designed for flexibility with various topologies using Go's `sync.Map`.

## Features

- Dijkstra's algorithm for finding the shortest path
- DFS-based path enumeration
- Blacklist support for excluding specific nodes
- Address and Edge abstractions with Base58 encoding
- Extensive unit testing with different network topologies (circle, subnet, multi-path, diamond, etc.)

## Project Structure

```
.
├── common.go           # Shared types and utilities (Address, EdgeId, Base58, MarshalText, etc.)
├── dijkstra.go         # Dijkstra shortest path algorithm
├── dfs.go              # DFS traversal and topology utilities
├── route.go            # Common interface definitions and path sorting
├── dijkstra_test.go    # Unit tests for Dijkstra algorithm
├── dfs_test.go         # Unit tests for DFS
├── go.mod              # Go module file
├── go.sum              # Dependency checksums
```

## Getting Started

### Prerequisites

- Go 1.18 or later
- Git

### Installation

Clone the repository:

```
git clone https://github.com/smallyunet/dijkstra-demo
cd dijkstra-demo
```

Install dependencies:

```
go mod tidy
```

Run the tests:

```
go test -v
```

## Usage

The main structs and functions:

- `DFS` and `Dijkstra` implement the `route` interface
- `NewTopology(nodes, edges, blacklist)` initializes the graph
- `GetShortPathTree(from, to)` returns shortest paths

Example use:

```
route := &Dijkstra{}
route.NewTopology(nodes, edges, blacklist)
shortestPath := route.GetShortPathTree(fromAddress, toAddress)
```

## Address and EdgeId Format

- `Address` is a 20-byte identifier
- `EdgeId` is a 40-byte structure composed of two addresses
- Conversion to/from Base58 is supported via `ToBase58` and `FromBase58`

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

