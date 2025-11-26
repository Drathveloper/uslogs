package logutils

import (
	"container/list"
	"sync"
)

// MaskPattern represents a pattern to mask
type MaskPattern struct {
	Start    string
	Mask     byte
	DelimMap [256]bool
}

// NewMaskPattern creates a new MaskPattern
func NewMaskPattern(start string, mask byte, delims ...byte) MaskPattern {
	var p MaskPattern
	p.Start = start
	p.Mask = mask
	for _, d := range delims {
		p.DelimMap[d] = true
	}
	return p
}

// A node in the trie structure used to implement Aho-Corasick
type node struct {
	root bool // true if this is the root

	b []byte // The blice at this node

	output bool // True means this node represents a blice that should
	// be output when matching
	index int // index into original dictionary if output is true

	counter uint64 // Set to the value of the Masker.counter when a
	// match is output to prevent duplicate output
	// The use of fixed size arrays is space-inefficient but fast for
	// lookups.

	child [256]*node // A non-nil entry in this array means that the
	// index represents a byte value which can be
	// appended to the current node. Blices in the
	// trie are built up byte by byte through these
	// child node pointers.

	fails [256]*node // Where to fail to (by following the fail
	// pointers) for each possible byte

	suffix *node // Pointer to the longest possible strict suffix of
	// this node

	fail *node // Pointer to the next node which is in the dictionary
	// which can be reached from here following suffixes. Called fail
	// because it is used to fallback in the trie when a match fails.
}

// Masker is returned by NewMatcher and contains a list of blices to
// match against
type Masker struct {
	counter uint64 // Counts the number of matches done, and is used to
	// prevent output of multiple matches of the same string
	trie []node // preallocated block of memory containing all the
	// nodes
	extent   int              // offset into trie that is currently free
	root     *node            // Points to trie[0]
	suffixes map[int][][]byte // The suffixes of the dictionary
	heap     sync.Pool        // a pool of haystacks to de-duplicate results in
	// a thread-safe manner
}

// findBlice looks for a blice in the trie starting from the root and
// returns a pointer to the node representing the end of the blice. If
// the blice is not found it returns nil.
func (m *Masker) findBlice(b []byte) *node {
	n := &m.trie[0]

	for n != nil && len(b) > 0 {
		n = n.child[int(b[0])]
		b = b[1:]
	}

	return n
}

// getFreeNode: gets a free node structure from the Masker's trie
// pool and updates the extent to point to the next free node.
func (m *Masker) getFreeNode() *node {
	m.extent += 1

	if m.extent == 1 {
		m.root = &m.trie[0]
		m.root.root = true
	}

	return &m.trie[m.extent-1]
}

// buildTrie builds the fundamental trie structure from a set of
// blices.
func (m *Masker) buildTrie(dictionary [][]byte) {

	// Work out the maximum size for the trie (all dictionary entries
	// are distinct plus the root). This is used to preallocate memory
	// for it.

	maxSize := 1
	for _, blice := range dictionary {
		maxSize += len(blice)
	}
	m.trie = make([]node, maxSize)

	// Calling this an ignoring its argument simply allocated
	// m.trie[0] which will be the root element

	m.getFreeNode()

	// This loop builds the nodes in the trie by following through
	// each dictionary entry building the children pointers.

	for i, blice := range dictionary {
		n := m.root
		var path []byte
		for _, b := range blice {
			path = append(path, b)

			c := n.child[int(b)]

			if c == nil {
				c = m.getFreeNode()
				n.child[int(b)] = c
				c.b = make([]byte, len(path))
				copy(c.b, path)

				// Nodes directly under the root node will have the
				// root as their fail point as there are no suffixes
				// possible.

				if len(path) == 1 {
					c.fail = m.root
				}

				c.suffix = m.root
			}

			n = c
		}

		// The last value of n points to the node representing a
		// dictionary entry

		n.output = true
		n.index = i
	}

	l := new(list.List)
	l.PushBack(m.root)

	for l.Len() > 0 {
		n := l.Remove(l.Front()).(*node)

		for i := 0; i < 256; i++ {
			c := n.child[i]
			if c != nil {
				l.PushBack(c)

				for j := 1; j < len(c.b); j++ {
					c.fail = m.findBlice(c.b[j:])
					if c.fail != nil {
						break
					}
				}

				if c.fail == nil {
					c.fail = m.root
				}

				for j := 1; j < len(c.b); j++ {
					s := m.findBlice(c.b[j:])
					if s != nil && s.output {
						c.suffix = s
						break
					}
				}
			}
		}
	}

	for i := 0; i < m.extent; i++ {
		for c := 0; c < 256; c++ {
			n := &m.trie[i]
			for n.child[c] == nil && !n.root {
				n = n.fail
			}

			m.trie[i].fails[c] = n
		}
	}

	m.trie = m.trie[:m.extent]
}

// NewMasker creates a new Masker used to mask blices based on a set of
// patterns
func NewMasker(dictionary ...string) *Masker {
	m := new(Masker)

	var d [][]byte
	for _, s := range dictionary {
		d = append(d, []byte(s))
	}

	m.buildTrie(d)

	return m
}

// Mask applies a mask to a blice based on a set of patterns
func (m *Masker) Mask(in []byte, patterns []MaskPattern) []byte {
	for i, b := range in {
		c := int(b)

		if !m.root.root && m.root.child[c] == nil {
			m.root = m.root.fails[c]
		}

		if m.root.child[c] != nil {
			f := m.root.child[c]
			m.root = f

			if f.output {
				m.applyMask(in, i, f.index, patterns)
			}

			for !f.suffix.root {
				f = f.suffix
				m.applyMask(in, i, f.index, patterns)
			}
		}
	}
	return in
}

func (m *Masker) applyMask(in []byte, pos int, index int, patterns []MaskPattern) {
	if index >= len(patterns) {
		return
	}

	p := patterns[index]

	i := pos + 1
	start := i

	for i < len(in) && !p.DelimMap[in[i]] {
		i++
	}

	if i >= len(in) {
		return
	}

	for j := start; j < i; j++ {
		in[j] = p.Mask
	}
}
