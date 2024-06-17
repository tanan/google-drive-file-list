package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_printNode(t *testing.T) {
	root1 := NewNode("root1", "root1", true, NewDate(2024, 1, 30, 11, 0), NewDate(2024, 1, 30, 11, 0))
	node1 := NewNode("node1", "node1.txt", false, NewDate(2024, 1, 30, 11, 0), NewDate(2024, 1, 31, 12, 0))

	root1.AddChild(node1)

	root2 := NewNode("root2", "root2", true, NewDate(2024, 1, 28, 9, 0), NewDate(2024, 1, 31, 12, 0))
	node2 := NewNode("node2", "node2.pdf", false, NewDate(2024, 1, 28, 9, 0), NewDate(2024, 1, 29, 10, 0))
	node3 := NewNode("node3", "node3", true, NewDate(2024, 1, 30, 11, 0), NewDate(2024, 1, 31, 12, 0))
	node4 := NewNode("node4", "node4.pdf", false, NewDate(2024, 1, 30, 11, 0), NewDate(2024, 1, 31, 12, 0))

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
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := csv.NewWriter(buf)
			printNode(w, tt.node, "/")
			w.Flush()
			got := buf.String()
			if diff := cmp.Diff(want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
