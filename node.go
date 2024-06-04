package main

type Node struct {
	ID       string
	Name     string
	IsDir    bool
	Children []*Node
}

func NewNode(id string, name string, isDir bool) *Node {
	return &Node{ID: id, Name: name, IsDir: isDir, Children: []*Node{}}
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}
