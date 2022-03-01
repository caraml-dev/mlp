package impl

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/flyteorg/flyteidl/clients/go/admin"
	adminpb "github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"
	"github.com/flyteorg/flytestdlib/config"
)

type Flyte struct {
	adminServiceClient service.AdminServiceClient
}

func NewFlyte(ctx context.Context) (*Flyte, error) {
	adminURL, err := url.Parse("flyte.pipelines.d.ai.golabs.io:80")
	if err != nil {
		return nil, fmt.Errorf("failed to parse admin URL: %s", err)
	}

	adminCfg := &admin.Config{
		Endpoint:              config.URL{URL: *adminURL},
		UseInsecureConnection: true,
		AuthType:              admin.AuthTypeClientSecret,
		ClientID:              "flytepropeller",
		ClientSecretLocation:  "/etc/flyte/client_secret",
		Scopes:                []string{"all"},
	}
	adminClient := admin.InitializeAdminClient(ctx, adminCfg)

	return &Flyte{
		adminServiceClient: adminClient,
	}, err
}

func (f *Flyte) ListPipelines(ctx context.Context) (interface{}, error) {
	in := &adminpb.ResourceListRequest{
		Id: &adminpb.NamedEntityIdentifier{
			Project: "flyteexamples",
			Domain:  "development",
		},
		Limit: 10,
	}

	workflowList, err := f.adminServiceClient.ListWorkflows(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %s", err)
	}
	log.Println("workflowList", workflowList)

	return workflowList, nil
}
