package dijkstra

import (
	"sync"
)

type ShortPathTree [][]Address

func (s ShortPathTree) Len() int {
	return len(s)
}
func (s ShortPathTree) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func (s ShortPathTree) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type route interface {
	NewTopology(nodes, edges *sync.Map, blacklist []Address)
	// GetShortPathTree return spt should contains path with format: [target ... media ... self]
	GetShortPathTree(from, to Address) ShortPathTree
}
