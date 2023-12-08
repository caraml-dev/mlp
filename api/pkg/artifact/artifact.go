package artifact

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// URL contains the information needed to identify the location of an object
// located in Google Cloud Storage.
type URL struct {
	// Bucket is the name of the Google Cloud Storage bucket where the object
	// is located.
	Bucket string

	// Object is the name and or path of the object stored in the bucket. It
	// should not start with a forward slash.
	Object string
}

type Service interface {
	ParseURL(gsURL string) (*URL, error)

	ReadArtifact(ctx context.Context, url string) ([]byte, error)
	WriteArtifact(ctx context.Context, url string, content []byte) error
	DeleteArtifact(ctx context.Context, url string) error
}

type GcsArtifactClient struct {
	API *storage.Client
}

func NewGcsArtifactClient(api *storage.Client) Service {
	return &GcsArtifactClient{
		API: api,
	}
}

// Parse parses a Google Cloud Storage string into a URL struct. The expected
// format of the string is gs://[bucket-name]/[object-path]. If the provided
// URL is formatted incorrectly an error will be returned.
func (gac *GcsArtifactClient) ParseURL(gsURL string) (*URL, error) {
	u, err := url.Parse(gsURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "gs" {
		return nil, err
	}

	bucket, object := u.Host, strings.TrimLeft(u.Path, "/")

	if bucket == "" {
		return nil, err
	}

	if object == "" {
		return nil, err
	}

	return &URL{
		Bucket: bucket,
		Object: object,
	}, nil
}

func (gac *GcsArtifactClient) ReadArtifact(ctx context.Context, url string) ([]byte, error) {
	u, err := gac.ParseURL(url)
	if err != nil {
		return nil, err
	}

	reader, err := gac.API.Bucket(u.Bucket).Object(u.Object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close() //nolint:errcheck

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (gac *GcsArtifactClient) WriteArtifact(ctx context.Context, url string, content []byte) error {
	u, err := gac.ParseURL(url)
	if err != nil {
		return err
	}
	w := gac.API.Bucket(u.Bucket).Object(u.Object).NewWriter(ctx)

	if _, err := fmt.Fprintf(w, "%s", content); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func (gac *GcsArtifactClient) DeleteArtifact(ctx context.Context, url string) error {
	u, err := gac.ParseURL(url)
	if err != nil {
		return err
	}

	// Sets the name for the bucket.
	bucket := gac.API.Bucket(u.Bucket)

	it := bucket.Objects(ctx, &storage.Query{
		Prefix: u.Object,
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

type NopArtifactClient struct{}

func NewNopArtifactClient() Service {
	return &NopArtifactClient{}
}

func (nac *NopArtifactClient) ParseURL(gsURL string) (*URL, error) {
	return nil, nil
}

func (nac *NopArtifactClient) ReadArtifact(ctx context.Context, url string) ([]byte, error) {
	return nil, nil
}

func (nac *NopArtifactClient) WriteArtifact(ctx context.Context, url string, content []byte) error {
	return nil
}

func (nac *NopArtifactClient) DeleteArtifact(ctx context.Context, url string) error {
	return nil
}
