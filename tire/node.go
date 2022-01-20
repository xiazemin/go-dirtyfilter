package tire

import (
	"fmt"
)

type node struct {
	end   bool
	child map[string]*node
}

func newNode() *node {
	return &node{
		child: make(map[string]*node),
	}
}

func (n *node) print(depth int) {
	for k, v := range n.child {
		v.print(depth + 1)
		fmt.Print(depth, "--", k, ":", v.end, ",")
	}
}
