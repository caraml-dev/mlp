package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/caraml-dev/mlp/api/config"
	"github.com/caraml-dev/mlp/api/log"
	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
)

type BootstrapOptions struct {
	ProjectReaders []string
	MLPAdmins      []string
}

var (
	bootstrapOpts = &BootstrapOptions{}
	bootstrapCmd  = &cobra.Command{
		Use:   "bootstrap",
		Short: "Start bootstrap job to populate Keto",
		Run: func(cmd *cobra.Command, args []string) {
			err := startKetoBootstrap(globalConfig, bootstrapOpts)
			if err != nil {
				log.Panicf("unable to bootstrap keto: %v", err)
			}
		},
	}
)

func init() {
	bootstrapCmd.Flags().StringSliceVarP(&bootstrapOpts.ProjectReaders, "project-readers", "r",
		[]string{}, "Comma separated list of project readers")
	bootstrapCmd.Flags().StringSliceVar(&bootstrapOpts.MLPAdmins, "mlp-admins", []string{},
		"Comma separated list of MLP admins")
}

func startKetoBootstrap(globalCfg *config.Config, bootstrapOpts *BootstrapOptions) error {
	authEnforcer, err := enforcer.NewEnforcerBuilder().
		KetoEndpoints(globalCfg.Authorization.KetoRemoteRead, globalCfg.Authorization.KetoRemoteWrite).
		Build()
	if err != nil {
		return err
	}
	updateRequest := enforcer.NewAuthorizationUpdateRequest()
	updateRequest.SetRoleMembers("mlp.projects.reader", bootstrapOpts.ProjectReaders)
	updateRequest.SetRoleMembers("mlp.admin", bootstrapOpts.MLPAdmins)
	err = authEnforcer.UpdateAuthorization(context.Background(), updateRequest)
	if err != nil {
		return err
	}
	return nil
}
