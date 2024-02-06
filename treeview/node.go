package treeview

import (
	"io"
)

type Node struct {
	Nodes []Node

	// will be displayed as node title
	Text string
}

func render(w io.Writer, node Node, indent, prefix string) {
	var title string
	if len(node.Nodes) == 0 {
		title += "* " // decorate node with no child nodes
	} else {
		title += "+ " // decorate node with at least one child node
	}

	if node.Text == "" {
		title += "<nil>"
	} else {
		title += node.Text
	}
	line := string([]rune(indent)[:len([]rune(indent))-len([]rune(prefix))]) + prefix + title + "\n"
	io.WriteString(w, line)

	if len(node.Nodes) == 0 {
		return
	}

	for i := 0; i < len(node.Nodes)-1; i++ {
		n := node.Nodes[i]
		render(w, n, indent+"│   ", "├── ")
	}

	n := node.Nodes[len(node.Nodes)-1]
	render(w, n, indent+"    ", "└── ")
}

func Render(w io.Writer, root Node) {
	render(w, root, "", "")
}
