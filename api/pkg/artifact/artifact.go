package artifact

import (
	"context"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type Service interface {
	DeleteArtifact(ctx context.Context, url string) error
}

type GcsArtifactClient struct {
	API *storage.Client
}

func (gac *GcsArtifactClient) DeleteArtifact(ctx context.Context, url string) error {
	// Get bucket name and gcsPrefix
	// the [5:] is to remove the "gs://" on the artifact uri
	// ex : gs://bucketName/path → bucketName/path
	gcsBucket, gcsLocation := gac.getGcsBucketAndLocation(url[5:])

	// Sets the name for the bucket.
	bucket := gac.API.Bucket(gcsBucket)

	it := bucket.Objects(ctx, &storage.Query{
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
		if err := bucket.Object(attrs.Name).Delete(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (gac *GcsArtifactClient) getGcsBucketAndLocation(str string) (string, string) {
	// Split string using delimiter
	// ex : bucketName/path/path1/item → (bucketName , path/path1/item)
	splitStr := strings.SplitN(str, "/", 2)
	return splitStr[0], splitStr[1]
}

func NewGcsArtifactClient(api *storage.Client) Service {
	return &GcsArtifactClient{
		API: api,
	}
}

type NopArtifactClient struct{}

func (nac *NopArtifactClient) DeleteArtifact(ctx context.Context, url string) error {
	return nil
}

func NewNopArtifactClient() Service {
	return &NopArtifactClient{}
}
