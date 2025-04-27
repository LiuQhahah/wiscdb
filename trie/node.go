package trie

type node struct {
	children map[byte]*node
	ignore   *node
	ids      []uint64
}
