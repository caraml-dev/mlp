package gcs

import (
	"fmt"
	"testing"

	"github.com/gojek/mlp/api/pkg/gcs/mocks"
)

func Test_gcsClient_DeleteArtifact(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "Valid URL",
			url:         "BucketName/path/to/item",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "Not a Valid Bucket",
			url:         "NotAValidBucketName/path/to/item",
			wantErr:     true,
			expectedErr: fmt.Errorf("storage: bucket doesn't exist"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gc := mocks.GcsService{}
			gc.On("DeleteArtifact", tc.url).Return(tc.expectedErr)
			if err := gc.DeleteArtifact(tc.url); (err != nil) != tc.wantErr {
				t.Errorf("gcsClient.DeleteArtifact() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
