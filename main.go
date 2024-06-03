package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/drive/v3"
)

const (
	folderID   = "xxxxx"
	folderType = "application/vnd.google-apps.folder"
)

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

func addChildren(srv *drive.Service, folderID string, node *Node) error {
	// List files in the current folder
	var files []*drive.File
	query := fmt.Sprintf("'%s' in parents and trashed=false", folderID)
	var pageToken string
	for {
		res, err := srv.Files.List().
			Q(query).
			PageToken(pageToken).
			PageSize(100).
			Fields("nextPageToken, files(id, name, mimeType, parents)").
			Do()
		if err != nil {
			log.Fatalf("Unable to retrieve files: %v", err)
		}

		files = append(files, res.Files...)

		pageToken = res.NextPageToken
		if pageToken == "" {
			break
		}
	}

	for _, f := range files {
		child := createNodeFromFile(f)
		node.AddChild(child)
		if f.MimeType == folderType {
			err := addChildren(srv, f.Id, child)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func createNodeFromFile(f *drive.File) *Node {
	if f.MimeType == folderType {
		return NewNode(f.Id, f.Name, true)
	}
	return NewNode(f.Id, f.Name, false)
}

func printNode(node *Node, prefix string) {
	fmt.Println(prefix + node.Name)
	if node.IsDir {
		for _, child := range node.Children {
			printNode(child, prefix+"  ")
		}
	}
}

func buildTree(srv *drive.Service, folderID string) (*Node, error) {
	root := NewNode(folderID, "root", true)
	if err := addChildren(srv, folderID, root); err != nil {
		return nil, err
	}
	return root, nil
}

func main() {
	ctx := context.Background()

	// Google Application Default Credentials are used for authentication.
	srv, err := drive.NewService(ctx)
	if err != nil {
		panic(err)
	}

	tree, _ := buildTree(srv, folderID)
	printNode(tree, "")
}
