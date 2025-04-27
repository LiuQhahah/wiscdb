package internal

import (
	"sync"
	"wiscdb/trie"
)

type publisher struct {
	sync.Mutex
	pubCh       chan requests
	subscribers map[uint64]subscriber
	nextID      uint64
	indexer     *trie.Trie
}

func (p *publisher) sendUpdates(reqs requests) {

}
