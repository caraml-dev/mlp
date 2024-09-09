package artifact

import (
	"reflect"
	"testing"
)

func TestArtifactClient_GetURLScheme(t *testing.T) {
	tests := []struct {
		name           string
		artifactClient SchemeInterface
		want           string
	}{
		{
			name: "gcs client",
			artifactClient: &GcsArtifactClient{
				URLScheme: "gs",
				API:       nil,
			},
			want: "gs",
		},
		{
			name: "s3 client",
			artifactClient: &S3ArtifactClient{
				URLScheme: "gs",
			},
			want: "gs",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.artifactClient.GetURLScheme(); got != tt.want {
				t.Errorf("GcsArtifactClient.GetScheme() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURLScheme_ParseURL(t *testing.T) {
	type args struct {
		gsURL string
	}
	tests := []struct {
		name      string
		urlScheme URLScheme
		args      args
		want      *URL
		wantErr   bool
	}{
		{
			name:      "valid short url",
			urlScheme: "gs",
			args: args{
				gsURL: "gs://bucket-name/object-path",
			},
			want: &URL{
				Bucket: "bucket-name",
				Object: "object-path",
			},
			wantErr: false,
		},
		{
			name:      "valid url",
			urlScheme: "gs",
			args: args{
				gsURL: "gs://bucket-name/object-path/object-path-2/object-path-3/file-1.txt",
			},
			want: &URL{
				Bucket: "bucket-name",
				Object: "object-path/object-path-2/object-path-3/file-1.txt",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.urlScheme.ParseURL(tt.args.gsURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("URLScheme.ParseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("URLScheme.ParseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
