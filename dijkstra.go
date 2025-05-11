// Package dijkstra provides an implementation of Dijkstra's shortest path algorithm
// with support for blacklisted nodes. It was extracted from the original route package.
package dijkstra

import (
	"math"
	"sync"
)

type Dijkstra struct {
	DirectedGraph *DirectedGraph
}

func (dijkstra *Dijkstra) NewTopology(nodes, edges *sync.Map, blacklist []Address) {
	directedGraph := NewDirectGraph(nodes, edges, blacklist)
	dijkstra.DirectedGraph = directedGraph
}

func (dijkstra *Dijkstra) GetShortPathTree(from, to Address) ShortPathTree {
	return dijkstra.DirectedGraph.GetPairShortestPath(from, to)
}

type DirectedGraph map[Address]Vertex

type Vertex struct {
	From map[Address]int64
	To   map[Address]int64
}

func NewDirectGraph(nodes *sync.Map, edges *sync.Map, blacklist []Address) *DirectedGraph {
	blackSet := make(map[Address]int)
	for k, v := range blacklist {
		blackSet[v] = k
	}
	dg := make(DirectedGraph)
	nodes.Range(func(key, value interface{}) bool {
		vernice := Vertex{
			From: make(map[Address]int64),
			To:   make(map[Address]int64),
		}
		if _, ok := blackSet[key.(Address)]; !ok {
			dg[key.(Address)] = vernice
		}
		return true
	})
	edges.Range(func(key, value interface{}) bool {
		addr1 := key.(EdgeId).GetAddr1()
		addr2 := key.(EdgeId).GetAddr2()
		_, ok1 := dg[addr1]
		_, ok2 := dg[addr2]
		if ok1 && ok2 {
			dg[addr1].From[addr2] = value.(int64)
			dg[addr2].To[addr1] = value.(int64)
		}
		return true
	})
	return &dg
}

func (g *DirectedGraph) GetPairShortestPath(from, to Address) ShortPathTree {
	dist := map[Address]int64{from: 0}
	prev := map[Address]int64{from: 0}
	pathList := map[Address][]Address{from: {from}}
	g.searchPathDijkstra(dist, prev, from, to, &pathList)
	path := pathList[to]
	// reverse path slice for interface format
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return ShortPathTree{path}
}

func (g *DirectedGraph) searchPathDijkstra(dist, prev map[Address]int64, from, to Address,
	pathList *map[Address][]Address) {
	// already find shortest path
	if from == to {
		return
	}
	// find all reachable node
	for node, distance := range (*g)[from].To {
		if _, ok := dist[node]; !ok {
			dist[node] = math.MaxInt64
		}
		if _, ok := prev[node]; !ok {
			dist[node] = func(a, b int64) int64 {
				if a < b {
					// dynamic update path which not in prev and have min distance
					temp := make([]Address, len((*pathList)[from]))
					copy(temp, (*pathList)[from])
					(*pathList)[node] = append(temp, node)
					return a
				}
				return b
			}(dist[from]+distance, dist[node])
		}
	}
	// find nearest node
	var nearestNode Address
	minDistance := int64(math.MaxInt64)
	for node, distance := range dist {
		if _, ok := prev[node]; !ok {
			if distance < minDistance {
				nearestNode = node
				minDistance = distance
			}
			// if more than one nearest node, select target node if exist
			if distance == minDistance && node == to {
				nearestNode = node
				minDistance = distance
			}
		}
	}
	// if no nearest node from this 'from'
	if minDistance == int64(math.MaxInt64) {
		return
	}
	// update prev set
	if _, ok := prev[nearestNode]; !ok {
		prev[nearestNode] = int64(len(prev))
	}
	g.searchPathDijkstra(dist, prev, nearestNode, to, pathList)
}
