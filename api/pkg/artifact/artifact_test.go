package artifact

import (
	"reflect"
	"testing"

	"cloud.google.com/go/storage"
)

func TestGcsArtifactClient_ParseURL(t *testing.T) {
	type fields struct {
		API *storage.Client
	}
	type args struct {
		gsURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *URL
		wantErr bool
	}{
		{
			name: "valid short url",
			fields: fields{
				API: nil,
			},
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
			name: "valid url",
			fields: fields{
				API: nil,
			},
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
			gac := &GcsArtifactClient{
				API: tt.fields.API,
			}
			got, err := gac.ParseURL(tt.args.gsURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GcsArtifactClient.ParseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GcsArtifactClient.ParseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
