package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/api/drive/v3"
)

var (
	folderID = flag.String("f", "", "folder id")
	driveID  = flag.String("d", "", "drive id")
)

const (
	FolderType = "application/vnd.google-apps.folder"
	PageLimit  = 100
	FileFields = "nextPageToken, files(id, name, mimeType, parents, createdTime, modifiedTime)"
	TimeFormat = "2006-01-02T15:04"
)

func addChildren(client *GoogleDriveClient, driveID string, folderID string, node *Node) error {
	// List files in the current folder
	var files []*drive.File
	query := fmt.Sprintf("'%s' in parents and trashed=false", folderID)
	var pageToken string
	for {
		opts := NewOptions(DriveID(driveID), Query(query), Fields(FileFields), PageToken(pageToken), PageSize(PageLimit))
		executor := Retry(client.ListFiles, 3, 10*time.Second)
		res, err := executor(context.Background(), opts)
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
			err := addChildren(client, driveID, f.Id, child)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func strToTime(t string) time.Time {
	parsedTime, _ := time.Parse(time.RFC3339, t)
	return parsedTime
}

func createNodeFromFile(f *drive.File) *Node {
	if f.MimeType == FolderType {
		return NewNode(f.Id, f.Name, true, strToTime(f.CreatedTime), strToTime(f.ModifiedTime))
	}
	return NewNode(f.Id, f.Name, false, strToTime(f.CreatedTime), strToTime(f.ModifiedTime))
}

func printNode(w *csv.Writer, node *Node, prefix string) {
	if !node.IsDir {
		w.Write([]string{fmt.Sprint(prefix + node.Name), node.CreatedTime.Format(TimeFormat), node.ModifiedTime.Format(TimeFormat)})
		return
	}
	for _, child := range node.Children {
		printNode(w, child, fmt.Sprintf("%s%s/", prefix, node.Name))
	}
}

func buildTree(client *GoogleDriveClient, driveID string, folderID string) (*Node, error) {
	opts := NewOptions(DriveID(driveID), FolderID(folderID))
	executor := Retry(client.GetFileName, 3, 5*time.Second)
	name, err := executor(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	root := NewNode(folderID, name, true, time.Now(), time.Now())
	if err := addChildren(client, driveID, folderID, root); err != nil {
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
	client := NewGoogleDriveClient(context.Background())

	tree, err := buildTree(client, *driveID, *folderID)
	if err != nil {
		log.Fatalf("Unable to build file tree: %v", err)
	}
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	printNode(w, tree, "/")
}
