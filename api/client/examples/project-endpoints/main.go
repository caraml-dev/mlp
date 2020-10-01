package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"

	"github.com/gojek/mlp/api/client"
)

func main() {
	ctx := context.Background()
	basePath := "http://mlp.dev/api/v1"
	if os.Getenv("MLP_API_BASEPATH") != "" {
		basePath = os.Getenv("MLP_API_BASEPATH")
	}

	// Create an HTTP client with Google default credential
	httpClient := http.DefaultClient
	googleClient, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/userinfo.email")
	if err == nil {
		httpClient = googleClient
	} else {
		log.Println("Google default credential not found. Fallback to HTTP default client")
	}

	cfg := client.NewConfiguration()
	cfg.BasePath = basePath
	cfg.HTTPClient = httpClient

	apiClient := client.NewAPIClient(cfg)

	// Get all projects
	projects, _, err := apiClient.ProjectApi.ProjectsGet(ctx, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Projects:", projects)

	for _, project := range projects {
		log.Println()
		log.Println("---")
		log.Println()

		log.Println("Project:", project.Name)

		// Update project
		updatedProject, _, err := apiClient.ProjectApi.ProjectsProjectIdPut(ctx, project.Id, client.Project{
			Team:   "dsp-new",
			Stream: "dsp-new",
			Name:   project.Name,
		})
		if err != nil {
			panic(err)
		}
		if updatedProject.Team != "dsp-new" {
			panic("Team should be changed to dsp-new")
		}
		if updatedProject.Stream != "dsp-new" {
			panic("Stream should be changed to dsp-new")
		}
		log.Printf("Project %s updated\n", project.Name)
	}
}
