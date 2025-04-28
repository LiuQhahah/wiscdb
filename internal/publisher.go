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
	if p.noOfSubscribers() != 0 {
		reqs.IncrRef()
		p.pubCh <- reqs
	}

}

func (p *publisher) noOfSubscribers() int {
	p.Lock()
	defer p.Unlock()
	return len(p.subscribers)
}
