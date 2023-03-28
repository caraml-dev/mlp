package artifact

import (
	"context"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type gcsClient struct {
	Api    *storage.Client
	Config Config
}
type Config struct {
	Ctx context.Context
}

type ArtifactService interface {
	DeleteArtifact(url string) error
}

func NewGcsClient(api *storage.Client, cfg Config) ArtifactService {
	return &gcsClient{
		Api:    api,
		Config: cfg,
	}
}

func (gc *gcsClient) DeleteArtifact(url string) error {
	// Get bucket name and gcsPrefix
	// the [5:] is to remove the "gs://" on the artifact uri
	// ex : gs://bucketName/path → bucketName/path
	gcsBucket, gcsLocation := gc.getGcsBucketAndLocation(url[5:])

	// Sets the name for the bucket.
	bucket := gc.Api.Bucket(gcsBucket)

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

func (gc *gcsClient) getGcsBucketAndLocation(str string) (string, string) {
	// Split string using delimiter
	// ex : bucketName/path/path1/item → (bucketName , path/path1/item)
	splitStr := strings.SplitN(str, "/", 2)
	return splitStr[0], splitStr[1]
}
