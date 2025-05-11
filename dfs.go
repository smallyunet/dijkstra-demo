package dijkstra

import (
	"sort"
	"sync"
)

type DFS struct {
	Topology *Topology
}

func (dfs *DFS) NewTopology(nodes, edges *sync.Map, blacklist []Address) {
	topology := NewTopology(nodes, edges, blacklist)
	dfs.Topology = topology
}

func (dns *DFS) GetShortPathTree(from, to Address) ShortPathTree {
	return dns.Topology.GetPairPathSorted(from, to)
}

// Topology represents a network topology
type Topology struct {
	nodes map[Address]int64
	edges map[Address]map[Address]int64
	sp    ShortPathTree
}

// Edge represents a directed edge in a graph
type Edge struct {
	NodeA    Address
	NodeB    Address
	Distance int64
}

// NewTopology creates a new topology
func NewTopology(nodes *sync.Map, edges *sync.Map, previousAddrs []Address) *Topology {
	t := &Topology{
		nodes: make(map[Address]int64),
		edges: make(map[Address]map[Address]int64),
	}

	nodes.Range(func(key, value interface{}) bool {
		tmpAddr := key.(Address)
		if !AddressContains(previousAddrs, tmpAddr) {
			t.nodes[tmpAddr] = value.(int64)
		}
		return true
	})

	edges.Range(func(key, value interface{}) bool {
		addr1 := key.(EdgeId).GetAddr1()
		addr2 := key.(EdgeId).GetAddr2()
		if !AddressContains(previousAddrs, addr1) && !AddressContains(previousAddrs, addr2) {
			if _, ok := t.edges[addr1]; !ok {
				t.edges[addr1] = make(map[Address]int64)
			}
			t.edges[addr1][addr2] = value.(int64)
		}
		return true
	})
	return t
}

func (self *Topology) GetAllPath(from Address) ShortPathTree {
	if 0 == len(self.nodes) || 0 == len(self.edges) {
		return [][]Address{}
	}
	var path []Address
	self.searchPathDFS(from, path)
	return self.sp
}

func (self *Topology) GetAllPathSorted(from Address) ShortPathTree {
	self.GetAllPath(from)
	sort.Sort(self.sp)
	return self.sp
}

func (self *Topology) searchPathDFS(from Address, path []Address) {
	path = append(path, from)
	// path may repeat when recurse return to upper layer
	for _, v := range self.sp {
		if sliceEqual(v, path) {
			return
		}
	}
	pathTemp := make([]Address, len(path))
	copy(pathTemp, path)
	self.sp = append(self.sp, pathTemp)
	for n := range self.edges[from] {
		// skip nodes already in path record
		walked := false
		for _, v := range path {
			if n == v {
				walked = true
				break
			}
		}
		if walked {
			continue
		}
		self.searchPathDFS(n, path)
	}
	return
}

func (self *Topology) GetPairPath(from, to Address) ShortPathTree {
	if 0 == len(self.nodes) || 0 == len(self.edges) {
		return [][]Address{}
	}
	var path []Address
	self.searchPathDFSTarget(from, to, path)
	return self.sp
}

func (self *Topology) GetPairPathSorted(from, to Address) ShortPathTree {
	self.GetPairPath(from, to)
	sort.Sort(self.sp)
	return self.sp
}

func (self *Topology) searchPathDFSTarget(from, to Address, path []Address) {
	path = append(path, to)
	if to == from {
		pathTemp := make([]Address, len(path))
		copy(pathTemp, path)
		self.sp = append(self.sp, pathTemp)
		return
	}
	for _, v := range self.sp {
		if sliceEqual(v, path) {
			return
		}
	}
	for n := range self.edges[to] {
		walked := false
		for _, v := range path {
			if n == v {
				walked = true
				break
			}
		}
		if walked {
			continue
		}
		self.searchPathDFSTarget(from, n, path)
	}
}

func sliceEqual(a, b []Address) bool {
	if len(a) != len(b) {
		return false
	}
	if (a == nil) != (b == nil) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
