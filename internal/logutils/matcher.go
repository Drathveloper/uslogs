package logutils

import (
	"container/list"
)

// MaskPattern represents a pattern to mask.
type MaskPattern struct {
	Start    string
	Mask     byte
	DelimMap [256]bool
}

// NewMaskPattern creates a new MaskPattern.
func NewMaskPattern(start string, mask byte, delimiters ...byte) MaskPattern {
	var pattern MaskPattern
	pattern.Start = start
	pattern.Mask = mask
	for _, delimiter := range delimiters {
		pattern.DelimMap[delimiter] = true
	}
	return pattern
}

// Masker is returned by NewMatcher and contains a list of blices to
// match against.
type Masker struct {
	root   *node
	trie   []node
	extent int
}

// NewMasker creates a new Masker used to mask blices based on a set of
// patterns.
func NewMasker(dictionary ...string) *Masker {
	masker := new(Masker)

	d := make([][]byte, 0)
	for _, s := range dictionary {
		d = append(d, []byte(s))
	}

	masker.buildTrie(d)

	return masker
}

// Mask applies a mask to a blice based on a set of patterns.
func (m *Masker) Mask(input []byte, patterns []MaskPattern) []byte {
	for index, item := range input {
		intItem := int(item)

		if !m.root.root && m.root.child[intItem] == nil {
			m.root = m.root.fails[intItem]
		}

		if m.root.child[intItem] != nil {
			childNode := m.root.child[intItem]
			m.root = childNode

			if childNode.output {
				m.applyMask(input, index, childNode.index, patterns)
			}

			for !childNode.suffix.root {
				childNode = childNode.suffix
				m.applyMask(input, index, childNode.index, patterns)
			}
		}
	}
	return input
}

// A node in the trie structure used to implement Aho-Corasick.
type node struct {
	child  [256]*node
	fails  [256]*node
	suffix *node
	fail   *node
	blice  []byte
	index  int
	root   bool
	output bool
}

// findBlice looks for a blice in the trie starting from the root and
// returns a pointer to the node representing the end of the blice. If
// the blice is not found it returns nil.
func (m *Masker) findBlice(blice []byte) *node {
	currNode := &m.trie[0]

	for currNode != nil && len(blice) > 0 {
		currNode = currNode.child[int(blice[0])]
		blice = blice[1:]
	}

	return currNode
}

// getFreeNode: gets a free node structure from the Masker's trie
// pool and updates the extent to point to the next free node.
func (m *Masker) getFreeNode() *node {
	m.extent++

	if m.extent == 1 {
		m.root = &m.trie[0]
		m.root.root = true
	}

	return &m.trie[m.extent-1]
}

// buildTrie builds the fundamental trie structure from a set of
// blices.
//
//nolint:gocognit,cyclop,funlen
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

	for bliceIndex, blice := range dictionary {
		currNode := m.root
		var path []byte
		for _, item := range blice {
			path = append(path, item)

			child := currNode.child[int(item)]

			if child == nil {
				child = m.getFreeNode()
				currNode.child[int(item)] = child
				child.blice = make([]byte, len(path))
				copy(child.blice, path)

				// Nodes directly under the root node will have the
				// root as their fail point as there are no suffixes
				// possible.

				if len(path) == 1 {
					child.fail = m.root
				}

				child.suffix = m.root
			}

			currNode = child
		}

		// The last value of n points to the node representing a
		// dictionary entry

		currNode.output = true
		currNode.index = bliceIndex
	}

	linkedList := new(list.List)
	linkedList.PushBack(m.root)

	for linkedList.Len() > 0 {
		n := linkedList.Remove(linkedList.Front()).(*node) //nolint:forcetypeassert

		for i := range 256 {
			child := n.child[i]
			if child != nil {
				linkedList.PushBack(child)

				for j := 1; j < len(child.blice); j++ {
					child.fail = m.findBlice(child.blice[j:])
					if child.fail != nil {
						break
					}
				}

				if child.fail == nil {
					child.fail = m.root
				}

				for j := 1; j < len(child.blice); j++ {
					s := m.findBlice(child.blice[j:])
					if s != nil && s.output {
						child.suffix = s
						break
					}
				}
			}
		}
	}

	for idx := range m.extent {
		for char := range 256 {
			trieItem := &m.trie[idx]
			for trieItem.child[char] == nil && !trieItem.root {
				trieItem = trieItem.fail
			}

			m.trie[idx].fails[char] = trieItem
		}
	}

	m.trie = m.trie[:m.extent]
}

func (m *Masker) applyMask(input []byte, endPos int, index int, patterns []MaskPattern) {
	if index >= len(patterns) {
		return
	}

	pattern := patterns[index]

	maskIndex := endPos + 1
	start := maskIndex

	for maskIndex < len(input) && !pattern.DelimMap[input[maskIndex]] {
		maskIndex++
	}

	if maskIndex >= len(input) {
		return
	}

	for idx := start; idx < maskIndex; idx++ {
		input[idx] = pattern.Mask
	}
}
