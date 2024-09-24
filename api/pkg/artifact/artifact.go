package artifact

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"

	"io"
	"net/url"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"google.golang.org/api/iterator"
)

const (
	gcsArtifactClientType = "gcs"
	gcsURLScheme          = "gs"
	s3ArtifactClientType  = "s3"
	s3URLScheme           = "s3"
)

var ErrObjectNotExist = errors.New("storage: object doesn't exist")

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

type URLScheme string

type SchemeInterface interface {
	GetURLScheme() string
	ParseURL(gsURL string) (*URL, error)
}

func (urlScheme URLScheme) GetURLScheme() string {
	return string(urlScheme)
}

// ParseURL parses an artifact storage string into a URL struct. The expected
// format of the string is [url-scheme]://[bucket-name]/[object-path]. If the provided
// URL is formatted incorrectly an error will be returned.
func (urlScheme URLScheme) ParseURL(gsURL string) (*URL, error) {
	u, err := url.Parse(gsURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != string(urlScheme) {
		return nil, fmt.Errorf("the scheme specified in the given URL is '%s' but the expected scheme is '%s'",
			u.Scheme, urlScheme)
	}

	bucket, object := u.Host, strings.TrimLeft(u.Path, "/")

	if bucket == "" {
		return nil, fmt.Errorf("the bucket in the given URL is an empty string")
	}

	if object == "" {
		return nil, fmt.Errorf("the object in the given URL is an empty string")
	}

	return &URL{
		Bucket: bucket,
		Object: object,
	}, nil
}

type Service interface {
	GetType() string
	GetURLScheme() string
	ParseURL(gsURL string) (*URL, error)
	ReadArtifact(ctx context.Context, url string) ([]byte, error)
	WriteArtifact(ctx context.Context, url string, content []byte) error
	DeleteArtifact(ctx context.Context, url string) error
}

type GcsArtifactClient struct {
	URLScheme
	api *storage.Client
}

func NewGcsArtifactClient() (Service, error) {
	api, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed initializing gcs for the artifact client with error: %s", err.Error())
	}
	return &GcsArtifactClient{
		URLScheme: gcsURLScheme,
		api:       api,
	}, nil
}

func (gac *GcsArtifactClient) GetType() string {
	return gcsArtifactClientType
}

func (gac *GcsArtifactClient) ReadArtifact(ctx context.Context, url string) ([]byte, error) {
	u, err := gac.ParseURL(url)
	if err != nil {
		return nil, err
	}

	reader, err := gac.api.Bucket(u.Bucket).Object(u.Object).NewReader(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil, ErrObjectNotExist
		}
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
	w := gac.api.Bucket(u.Bucket).Object(u.Object).NewWriter(ctx)

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
	bucket := gac.api.Bucket(u.Bucket)

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

type S3ArtifactClient struct {
	URLScheme
	client *s3.Client
}

func NewS3ArtifactClient() (Service, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s,failed loading s3 config for the artifact client", err.Error())
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("AWS_ENDPOINT_URL"))
	})
	return &S3ArtifactClient{
		URLScheme: s3URLScheme,
		client:    client,
	}, nil
}

func (s3c *S3ArtifactClient) GetType() string {
	return s3ArtifactClientType
}

func (s3c *S3ArtifactClient) ReadArtifact(ctx context.Context, url string) ([]byte, error) {
	u, err := s3c.ParseURL(url)
	if err != nil {
		return nil, err
	}

	reader, err := s3c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(u.Object),
	})
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return nil, ErrObjectNotExist
		}
		return nil, err
	}
	defer reader.Body.Close() //nolint:errcheck

	bytes, err := io.ReadAll(reader.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (s3c *S3ArtifactClient) WriteArtifact(ctx context.Context, url string, content []byte) error {
	u, err := s3c.ParseURL(url)
	if err != nil {
		return err
	}

	_, err = s3c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(u.Object),
		Body:   bytes.NewReader(content),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s3c *S3ArtifactClient) DeleteArtifact(ctx context.Context, url string) error {
	u, err := s3c.ParseURL(url)
	if err != nil {
		return err
	}

	// TODO: To confirm and refactor if versioning is enabled on S3 and to specify the versionId to be deleted
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(u.Object),
	}

	_, err = s3c.client.DeleteObject(ctx, input)
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return ErrObjectNotExist
		}
		return err
	}
	return nil
}

type NopArtifactClient struct{}

func NewNopArtifactClient() Service {
	return &NopArtifactClient{}
}

func (nac *NopArtifactClient) GetType() string {
	return ""
}

func (nac *NopArtifactClient) GetURLScheme() string {
	return ""
}

func (nac *NopArtifactClient) ParseURL(_ string) (*URL, error) {
	return nil, nil
}

func (nac *NopArtifactClient) ReadArtifact(_ context.Context, _ string) ([]byte, error) {
	return nil, nil
}

func (nac *NopArtifactClient) WriteArtifact(_ context.Context, _ string, _ []byte) error {
	return nil
}

func (nac *NopArtifactClient) DeleteArtifact(_ context.Context, _ string) error {
	return nil
}
