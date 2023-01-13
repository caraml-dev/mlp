package cluster

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"k8s.io/client-go/rest"
	clientcmdapiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

func TestK8sClusterCredsToRestConfig(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		input     K8sConfig
		output    rest.Config
	}{
		{
			name:      "Test basic server, ca cert",
			wantError: false,
			input: K8sConfig{
				Name: "dummy-cluster",
				Cluster: &clientcmdapiv1.Cluster{
					Server:                "https://some_ip_address",
					InsecureSkipTLSVerify: true,
				},
				AuthInfo: &clientcmdapiv1.AuthInfo{
					ClientCertificateData: []byte(`ABCDEF`),
					ClientKeyData:         []byte(`12345`),
				},
			},
			output: rest.Config{
				Host: "https://some_ip_address",
				TLSClientConfig: rest.TLSClientConfig{
					Insecure: true,
					CertData: []byte(`ABCDEF`),
					KeyData:  []byte(`12345`),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credsM := NewK8sClusterCreds(&tt.input)
			res, err := credsM.ToRestConfig()
			if err != nil && !tt.wantError {
				t.Errorf("Error not expected but occurred: %s", err.Error())
			}
			if diff := cmp.Diff(res, &tt.output); diff != "" {
				t.Errorf("diff is not empty %s", diff)
			}
		})
	}
}

func TestGenerateKubeConfig(t *testing.T) {

	tests := []struct {
		name      string
		wantError bool
		input     K8sConfig
		output    *clientcmdapiv1.Config
	}{
		{
			name:      "Test basic server, ca cert",
			wantError: false,
			input: K8sConfig{
				Name: "dummy-cluster",
				Cluster: &clientcmdapiv1.Cluster{
					Server:                "https://some_ip_address",
					InsecureSkipTLSVerify: true,
				},
				AuthInfo: &clientcmdapiv1.AuthInfo{
					ClientCertificateData: []byte(`ABCDEF`),
					ClientKeyData:         []byte(`12345`),
				},
			},
			output: &clientcmdapiv1.Config{
				Clusters: []clientcmdapiv1.NamedCluster{
					{
						Name: "dummy-cluster",
						Cluster: clientcmdapiv1.Cluster{
							Server:                "https://some_ip_address",
							InsecureSkipTLSVerify: true,
						},
					},
				},
				AuthInfos: []clientcmdapiv1.NamedAuthInfo{
					{
						Name: k8sUser,
						AuthInfo: clientcmdapiv1.AuthInfo{
							ClientCertificateData: []byte(`ABCDEF`),
							ClientKeyData:         []byte(`12345`),
						},
					},
				},
				Contexts: []clientcmdapiv1.NamedContext{
					{
						Name: "dummy-cluster-user",
						Context: clientcmdapiv1.Context{
							Cluster:  "dummy-cluster",
							AuthInfo: k8sUser,
						},
					},
				},
				CurrentContext: "dummy-cluster-user",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := generateKubeConfig(&tt.input)
			if diff := cmp.Diff(res, tt.output); diff != "" {
				t.Errorf("diff is not empty %s", diff)
			}
		})
	}
}
