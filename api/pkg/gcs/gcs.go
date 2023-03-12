package gcs

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type gcsClient struct {
	Client *storage.Client
	Config Config
}
type Config struct {
	Ctx context.Context
}
type GcsPackage interface {
	DeleteArtifact(url string)
}

func NewGcsClient(client *storage.Client, cfg Config) *gcsClient {
	return &gcsClient{
		Client: client,
		Config: cfg,
	}
}

func (gc *gcsClient) DeleteArtifact(url string) {
	// Sets the name for the new bucket.
	gcsBucket, gcsLocation := gc.RemoveAndSplit(url, "/")
	fmt.Println(gcsBucket)
	fmt.Println(gcsLocation)

	bucket := gc.Client.Bucket(gcsBucket)

	it := bucket.Objects(gc.Config.Ctx, &storage.Query{
		Prefix: gcsLocation,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		if err := bucket.Object(attrs.Name).Delete(gc.Config.Ctx); err != nil {
			// TODO: Handle error.
			fmt.Println(err)
			fmt.Println("Error Deleting Object")
		}
		fmt.Println(attrs.Name)

	}
}

func (gc *gcsClient) RemoveAndSplit(str, delimiter string) (string, string) {
	// Remove first 4 characters
	str = str[5:]

	// Split string using delimiter
	splitStr := strings.SplitN(str, delimiter, 2)

	return splitStr[0], splitStr[1]
}
