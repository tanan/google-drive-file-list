package main

import "time"

type Node struct {
	ID           string
	Name         string
	IsDir        bool
	CreatedTime  time.Time
	ModifiedTime time.Time
	Children     []*Node
}

func NewNode(id string, name string, isDir bool, createdTime time.Time, modifiedTime time.Time) *Node {
	return &Node{
		ID:           id,
		Name:         name,
		CreatedTime:  createdTime,
		ModifiedTime: modifiedTime,
		IsDir:        isDir,
		Children:     []*Node{},
	}
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}
