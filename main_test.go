package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_printNode(t *testing.T) {
	root1 := NewNode("root1", "root1", true)
	node1 := NewNode("node1", "node1.txt", false)

	root1.AddChild(node1)

	root2 := NewNode("root2", "root2", true)
	node2 := NewNode("node2", "node2.pdf", false)
	node3 := NewNode("node3", "node3", true)
	node4 := NewNode("node4", "node4.pdf", false)

	node3.AddChild(node4)
	root2.AddChild(node2)
	root2.AddChild(node3)

	tests := []struct {
		name string
		node *Node
	}{
		{name: "printNode_test1", node: root1},
		{name: "printNode_test2", node: root2},
	}
	for _, tt := range tests {
		testfile := filepath.Join("testdata", fmt.Sprintf("%s.txt", tt.name))
		want, err := readFile(testfile)
		if err != nil {
			t.Fatalf("error readFile(): testfile: %s, error: %s", testfile, err.Error())
		}
		w := &bytes.Buffer{}
		t.Run(tt.name, func(t *testing.T) {
			printNode(w, tt.node, "/")
			got := w.String()
			if diff := cmp.Diff(want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
