package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"

	"cloud.google.com/go/resourcemanager/apiv3"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	resourcemanagerpb "google.golang.org/genproto/googleapis/cloud/resourcemanager/v3"
)

type Folder struct {
	Name     string                           `json:"name"`
	ID       string                           `json:"id"`
	Children []*Folder                        `json:"children,omitempty"`
	Projects []*resourcemanagerpb.Project `json:"projects,omitempty"`
}

func collectProjects(ctx context.Context, client *resourcemanager.ProjectsClient, folderID string, folder *Folder) {
	req := &resourcemanagerpb.ListProjectsRequest{Parent: "folders/" + folderID}
	it := client.ListProjects(ctx, req)
	for {
		proj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to list projects: %v", err)
		}
		folder.Projects = append(folder.Projects, proj)
	}
}

func buildTreeConcurrently(ctx context.Context, client *resourcemanager.FoldersClient, projectClient *resourcemanager.ProjectsClient, folderID string, folderChan chan<- *Folder, wg *sync.WaitGroup) {
	defer wg.Done()

	folder := &Folder{}
	folder.ID = folderID
	resp, err := client.GetFolder(ctx, &resourcemanagerpb.GetFolderRequest{Name: "folders/" + folderID})
	if err != nil {
		log.Fatalf("Failed to get folder: %v", err)
	}
	folder.Name = resp.Name[len("folders/"):]

	req := &resourcemanagerpb.ListFoldersRequest{Parent: "folders/" + folderID}
	it := client.ListFolders(ctx, req)

	var childWg sync.WaitGroup
	for {
		subFolder, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to list sub-folders: %v", err)
		}
		childChan := make(chan *Folder, 1)
		childWg.Add(1)
		go buildTreeConcurrently(ctx, client, projectClient, subFolder.Name[len("folders/"):], childChan, &childWg)
		child := <-childChan
		folder.Children = append(folder.Children, child)
	}
	childWg.Wait()
	collectProjects(ctx, projectClient, folderID, folder)
	folderChan <- folder
}

func main() {
	keyFile := flag.String("key-file", "", "path to GCP ServiceAccount JSON file")
	folderIDs := flag.String("folders", "", "comma-separated list of folder IDs")
	flag.Parse()

	if *keyFile == "" {
		log.Fatal("Missing required argument --key-file")
	}
	if *folderIDs == "" {
		log.Fatal("Missing required argument --folders")
	}

	folders := strings.Split(*folderIDs, ",")

	ctx := context.Background()
	folderClient, err := resourcemanager.NewFoldersClient(ctx, option.WithCredentialsFile(*keyFile))
	if err != nil {
		log.Fatalf("Failed to create folder client: %v", err)
	}
	defer folderClient.Close()

	projectClient, err := resourcemanager.NewProjectsClient(ctx, option.WithCredentialsFile(*keyFile))
	if err != nil {
		log.Fatalf("Failed to create project client: %v", err)
	}
	defer projectClient.Close()

	folderChan := make(chan *Folder, len(folders))
	var wg sync.WaitGroup
	for _, folderID := range folders {
		wg.Add(1)
		go buildTreeConcurrently(ctx, folderClient, projectClient, folderID, folderChan, &wg)
	}
	wg.Wait()
	close(folderChan)

	var rootFolders []*Folder
	for folder := range folderChan {
		rootFolders = append(rootFolders, folder)
	}

	jsonResult, err := json.MarshalIndent(rootFolders, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON result: %v", err)
	}

	fmt.Println(string(jsonResult))
}
