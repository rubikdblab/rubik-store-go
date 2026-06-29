package skiplist

import (
	"bytes"
	"math/rand"
)

type EntryType uint8

const (
	EntryPut    EntryType = 0
	EntryDelete EntryType = 1
)

type Entry struct {
	Key       []byte
	Value     []byte
	SeqNum    uint64
	EntryType EntryType
}

type node struct {
	entry   *Entry
	forward []*node
}

type Skiplist struct {
	sentinel    *node
	maxLevel    int
	level       int
	length      int
	probability int
}

func New(maxLevel int) *Skiplist {
	return &Skiplist{
		sentinel: &node{
			forward: make([]*node, maxLevel),
		},
		maxLevel:    maxLevel,
		level:       0,
		length:      0,
		probability: 50,
	}
}

func (s *Skiplist) Search(key []byte) *Entry {
	// Empty Skiplist
	if s.length == 0 {
		return nil
	}

	n := s.sentinel
	for l := s.level; l >= 0; l-- {
		for n.forward[l] != nil && bytes.Compare(n.forward[l].entry.Key, key) < 0 {
			n = n.forward[l]
		}
	}

	candidate := n.forward[0]
	if candidate != nil && bytes.Equal(key, candidate.entry.Key) {
		return candidate.entry
	}
	return nil
}

func (s *Skiplist) Insert(entry *Entry) {
	toUpdate := make([]*node, s.level+1)

	n := s.sentinel
	for l := s.level; l >= 0; l-- {
		for n.forward[l] != nil && bytes.Compare(n.forward[l].entry.Key, entry.Key) < 0 {
			n = n.forward[l]
		}
		toUpdate[l] = n
	}

	if n.forward[0] != nil && bytes.Equal(entry.Key, n.forward[0].entry.Key) {
		n.forward[0].entry = entry
		return
	}

	nn := &node{
		entry:   entry,
		forward: make([]*node, s.maxLevel),
	}

	l := 0
	for ; l < len(toUpdate); l++ {
		tu := toUpdate[l]
		nn.forward[l] = tu.forward[l]
		tu.forward[l] = nn
		if rand.Intn(100) < s.probability {
			break
		}
	}
	if l == s.level+1 && s.level+1 != s.maxLevel && rand.Intn(100) < s.probability {
		s.level++
		s.sentinel.forward[s.level] = nn
	}
	s.length++
}

func (s *Skiplist) Remove(key []byte) {
	toUpdate := make([]*node, s.level+1)

	n := s.sentinel
	for l := s.level; l >= 0; l-- {
		for n.forward[l] != nil && bytes.Compare(n.forward[l].entry.Key, key) < 0 {
			n = n.forward[l]
		}
		toUpdate[l] = n
	}

	if n.forward[0] == nil || !bytes.Equal(n.forward[0].entry.Key, key) {
		return
	}

	toDelete := n.forward[0]
	n.forward[0] = toDelete.forward[0]

	for l := 1; l <= s.level; l++ {
		if toUpdate[l].forward[l] == toDelete {
			toUpdate[l].forward[l] = toDelete.forward[l]
		}
	}
	s.length--
}
