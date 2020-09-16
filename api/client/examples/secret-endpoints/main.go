package main

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"

	"github.com/gojek/mlp/client"
)

func main() {
	ctx := context.Background()
	basePath := "http://mlp.dev/api/v1"
	if os.Getenv("MLP_API_BASEPATH") != "" {
		basePath = os.Getenv("MLP_API_BASEPATH")
	}

	// Create an HTTP client with Google default credential
	googleClient, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/userinfo.email")
	if err != nil {
		panic(err)
	}

	cfg := client.NewConfiguration()
	cfg.BasePath = basePath
	cfg.HTTPClient = googleClient

	apiClient := client.NewAPIClient(cfg)

	apiClient.ProjectApi.ProjectsGet(ctx, nil)

	// Get all projects
	projects, _, err := apiClient.ProjectApi.ProjectsGet(ctx, nil)
	if err != nil {
		panic(err)
	}

	for _, project := range projects {
		fmt.Println()
		fmt.Println("---")
		fmt.Println()

		fmt.Println("Project:", project.Name)

		_, _, err := apiClient.SecretApi.ProjectsProjectIdSecretsPost(ctx, project.Id, client.Secret{
			Name: project.Name,
			Data: `{"data": "encrypted"}`,
		})
		if err != nil {
			panic(err)
		}

		// Get all secrets for projects
		secrets, _, err := apiClient.SecretApi.ProjectsProjectIdSecretsGet(ctx, project.Id)
		if err != nil {
			panic(err)
		}

		for _, secret := range secrets {
			fmt.Println("Secret name:", secret.Name)
			_, err := apiClient.SecretApi.ProjectsProjectIdSecretsSecretIdDelete(ctx, project.Id, secret.Id)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Secret %s: deleted", secret.Name)

		}

	}
}
