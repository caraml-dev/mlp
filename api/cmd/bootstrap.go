package cmd

import (
	"context"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"

	"github.com/spf13/cobra"

	"github.com/caraml-dev/mlp/api/config"
	"github.com/caraml-dev/mlp/api/log"
	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
)

type BootstrapRoleMembers struct {
	ProjectReaders []string
	MLPAdmins      []string
}

var (
	bootstrapRoleMembersInputFile string
	bootstrapCmd                  = &cobra.Command{
		Use:   "bootstrap",
		Short: "Start bootstrap job to populate Keto",
		Run: func(cmd *cobra.Command, args []string) {
			bootstrapRoleMembers, err := loadRoleMemberFromInputFile(bootstrapRoleMembersInputFile)
			if err != nil {
				log.Panicf("unable to load role members from input file: %v", err)
			}
			err = startKetoBootstrap(globalConfig, bootstrapRoleMembers)
			if err != nil {
				log.Panicf("unable to bootstrap keto: %v", err)
			}
		},
	}
)

func init() {
	bootstrapCmd.Flags().StringVarP(&bootstrapRoleMembersInputFile, "role-members", "r", "",
		"Path to an input file that map roles to members")
	err := bootstrapCmd.MarkFlagRequired("role-members")
	if err != nil {
		log.Panicf("unable to mark flag as required: %v", err)
	}
}

func loadRoleMemberFromInputFile(path string) (*BootstrapRoleMembers, error) {
	bootstrapRoleMembers := &BootstrapRoleMembers{
		ProjectReaders: []string{},
		MLPAdmins:      []string{},
	}
	k := koanf.New(".")
	err := k.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		return nil, err
	}
	err = k.Unmarshal("", bootstrapRoleMembers)
	if err != nil {
		return nil, err
	}
	return bootstrapRoleMembers, nil
}

func startKetoBootstrap(globalCfg *config.Config, bootstrapOpts *BootstrapRoleMembers) error {
	authEnforcer, err := enforcer.NewEnforcerBuilder().
		KetoEndpoints(globalCfg.Authorization.KetoRemoteRead, globalCfg.Authorization.KetoRemoteWrite).
		Build()
	if err != nil {
		return err
	}
	updateRequest := enforcer.NewAuthorizationUpdateRequest()
	updateRequest.SetRoleMembers(enforcer.MLPProjectsReaderRole, bootstrapOpts.ProjectReaders)
	updateRequest.SetRoleMembers(enforcer.MLPAdminRole, bootstrapOpts.MLPAdmins)
	err = authEnforcer.UpdateAuthorization(context.Background(), updateRequest)
	if err != nil {
		return err
	}
	return nil
}
