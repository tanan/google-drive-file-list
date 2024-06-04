package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/api/drive/v3"
)

var (
	folderID = flag.String("f", "", "folder id")
	driveID  = flag.String("d", "", "drive id")
)

const (
	FolderType = "application/vnd.google-apps.folder"
	PageLimit  = 100
	FileFields = "nextPageToken, files(id, name, mimeType, parents)"
)

func getFileName(srv *drive.Service, id string) (string, error) {
	res, err := srv.Files.Get(id).SupportsAllDrives(true).Fields("name").Do()
	if err != nil {
		return "", err
	}
	return res.Name, nil
}

func addChildren(srv *drive.Service, folderID string, node *Node) error {
	// List files in the current folder
	var files []*drive.File
	query := fmt.Sprintf("'%s' in parents and trashed=false", folderID)
	var pageToken string
	for {
		res, err := srv.Files.List().
			Corpora("drive").
			IncludeItemsFromAllDrives(true).
			SupportsAllDrives(true).
			DriveId(*driveID).
			Q(query).
			Fields(FileFields).
			PageToken(pageToken).
			PageSize(PageLimit).
			Do()
		if err != nil {
			log.Printf("Unable to retrieve files: %v", err)
			return err
		}

		files = append(files, res.Files...)

		pageToken = res.NextPageToken
		if pageToken == "" {
			break
		}
	}

	// Create node and add as a child node
	for _, f := range files {
		child := createNodeFromFile(f)
		node.AddChild(child)
		if f.MimeType == FolderType {
			err := addChildren(srv, f.Id, child)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func createNodeFromFile(f *drive.File) *Node {
	if f.MimeType == FolderType {
		return NewNode(f.Id, f.Name, true)
	}
	return NewNode(f.Id, f.Name, false)
}

func printNode(node *Node, prefix string) {
	if !node.IsDir {
		fmt.Println(prefix + node.Name)
		return
	}
	for _, child := range node.Children {
		printNode(child, fmt.Sprintf("%s%s/", prefix, node.Name))
	}
}

func buildTree(srv *drive.Service, folderID string) (*Node, error) {
	name, err := getFileName(srv, folderID)
	if err != nil {
		return nil, err
	}
	root := NewNode(folderID, name, true)
	if err := addChildren(srv, folderID, root); err != nil {
		return nil, err
	}
	return root, nil
}

func main() {
	flag.Parse()
	if *folderID == "" {
		log.Fatal("-f (folderID) must be set")
	}

	// Google Application Default Credentials are used for authentication.
	ctx := context.Background()
	srv, err := drive.NewService(ctx)
	if err != nil {
		panic(err)
	}

	tree, err := buildTree(srv, *folderID)
	if err != nil {
		log.Fatalf("Unable to build file tree: %v", err)
	}
	printNode(tree, "/")
}
