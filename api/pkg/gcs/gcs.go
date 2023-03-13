package gcs

import (
	"context"
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
	DeleteArtifact(url string) error
}

func NewGcsClient(client *storage.Client, cfg Config) *gcsClient {
	return &gcsClient{
		Client: client,
		Config: cfg,
	}
}

func (gc *gcsClient) DeleteArtifact(url string) error {
	// Get bucket name and gcsPrefix
	gcsBucket, gcsLocation := gc.RemoveAndSplit(url, "/")

	// Sets the name for the bucket.
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
			return err
		}
		if err := bucket.Object(attrs.Name).Delete(gc.Config.Ctx); err != nil {
			return err
		}
	}
	return nil
}

func (gc *gcsClient) RemoveAndSplit(str, delimiter string) (string, string) {
	// Split string using delimiter
	splitStr := strings.SplitN(str, delimiter, 2)

	return splitStr[0], splitStr[1]
}
