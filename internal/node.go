package internal

import (
	"sort"

	"github.com/onsi/ginkgo/types"
	"sync"
)

var _global_node_id_counter = uint(0)
var _global_id_mutex = &sync.Mutex{}

func UniqueNodeID() uint {
	//There's a reace in the internal integration tests if we don't make
	//accessing _global_node_id_counter safe across goroutines.
	_global_id_mutex.Lock()
	defer _global_id_mutex.Unlock()
	_global_node_id_counter += 1
	return _global_node_id_counter
}

/* Node */
type Node struct {
	ID       uint
	NodeType types.NodeType

	Text         string
	Body         func()
	CodeLocation types.CodeLocation
	NestingLevel int

	MarkedFocus   bool
	MarkedPending bool
}

func NewNode(nodeType types.NodeType, text string, body func(), codeLocation types.CodeLocation, markedFocus bool, markedPending bool) Node {
	return Node{
		ID:            UniqueNodeID(),
		NodeType:      nodeType,
		Text:          text,
		Body:          body,
		CodeLocation:  codeLocation,
		MarkedFocus:   markedFocus,
		MarkedPending: markedPending,
		NestingLevel:  -1,
	}
}

func (n Node) IsZero() bool {
	return n.ID == 0
}

/* Nodes */
type Nodes []Node

func (n Nodes) CopyAppend(nodes ...Node) Nodes {
	out := Nodes{}
	for _, node := range n {
		out = append(out, node)
	}
	for _, node := range nodes {
		out = append(out, node)
	}
	return out
}

func (n Nodes) SplitAround(pivot Node) (Nodes, Nodes) {
	left := Nodes{}
	right := Nodes{}
	found := false
	for _, node := range n {
		if node.ID == pivot.ID {
			found = true
			continue
		}
		if found {
			right = append(right, node)
		} else {
			left = append(left, node)
		}
	}

	return left, right
}

func (n Nodes) FirstNodeWithType(nodeTypes ...types.NodeType) Node {
	for _, node := range n {
		if node.NodeType.Is(nodeTypes...) {
			return node
		}
	}
	return Node{}
}

func (n Nodes) WithType(nodeTypes ...types.NodeType) Nodes {
	out := Nodes{}
	for _, node := range n {
		if node.NodeType.Is(nodeTypes...) {
			out = append(out, node)
		}
	}
	return out
}

func (n Nodes) WithoutType(nodeTypes ...types.NodeType) Nodes {
	out := Nodes{}
	for _, node := range n {
		if !node.NodeType.Is(nodeTypes...) {
			out = append(out, node)
		}
	}
	return out
}

func (n Nodes) WithinNestingLevel(deepestNestingLevel int) Nodes {
	out := Nodes{}
	for _, node := range n {
		if node.NestingLevel <= deepestNestingLevel {
			out = append(out, node)
		}
	}
	return out
}

func (n Nodes) SortedByDescendingNestingLevel() Nodes {
	out := Nodes{}
	for _, node := range n {
		out = append(out, node)
	}
	sort.SliceStable(out, func(i int, j int) bool {
		return out[i].NestingLevel > out[j].NestingLevel
	})

	return out
}

func (n Nodes) Texts() []string {
	out := []string{}
	for _, node := range n {
		out = append(out, node.Text)
	}
	return out
}

func (n Nodes) CodeLocations() []types.CodeLocation {
	out := []types.CodeLocation{}
	for _, node := range n {
		out = append(out, node.CodeLocation)
	}
	return out
}

func (n Nodes) BestTextFor(node Node) string {
	if node.Text != "" {
		return node.Text
	}
	parentNestingLevel := node.NestingLevel - 1
	for _, node := range n {
		if node.Text != "" && node.NestingLevel == parentNestingLevel {
			return node.Text
		}
	}

	return ""
}

func (n Nodes) HasNodeMarkedPending() bool {
	for _, node := range n {
		if node.MarkedPending {
			return true
		}
	}
	return false
}

func (n Nodes) HasNodeMarkedFocus() bool {
	for _, node := range n {
		if node.MarkedFocus {
			return true
		}
	}
	return false
}