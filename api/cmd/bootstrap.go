package cmd

import (
	"context"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"

	"github.com/spf13/cobra"

	"github.com/caraml-dev/mlp/api/log"
	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
)

type BootstrapConfig struct {
	KetoRemoteRead  string
	KetoRemoteWrite string
	ProjectReaders  []string
	MLPAdmins       []string
}

var (
	bootstrapConfigFile string
	bootstrapCmd        = &cobra.Command{
		Use:   "bootstrap",
		Short: "Start bootstrap job to populate Keto",
		Run: func(cmd *cobra.Command, args []string) {
			bootstrapConfig, err := loadBootstrapConfig(bootstrapConfigFile)
			if err != nil {
				log.Panicf("unable to load role members from input file: %v", err)
			}
			err = startKetoBootstrap(bootstrapConfig)
			if err != nil {
				log.Panicf("unable to bootstrap keto: %v", err)
			}
		},
	}
)

func init() {
	bootstrapCmd.Flags().StringVarP(&bootstrapConfigFile, "config", "c", "",
		"Path to keto bootstrap configuration")
	err := bootstrapCmd.MarkFlagRequired("config")
	if err != nil {
		log.Panicf("unable to mark flag as required: %v", err)
	}
}

func loadBootstrapConfig(path string) (*BootstrapConfig, error) {
	bootstrapCfg := &BootstrapConfig{
		ProjectReaders: []string{},
		MLPAdmins:      []string{},
	}
	k := koanf.New(".")
	err := k.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		return nil, err
	}
	err = k.Unmarshal("", bootstrapCfg)
	if err != nil {
		return nil, err
	}
	return bootstrapCfg, nil
}

func startKetoBootstrap(bootstrapCfg *BootstrapConfig) error {
	authEnforcer, err := enforcer.NewEnforcerBuilder().
		KetoEndpoints(bootstrapCfg.KetoRemoteRead, bootstrapCfg.KetoRemoteWrite).
		Build()
	if err != nil {
		return err
	}

	defaultPermissions := []string{"mlp.projects.post"}

	updateRequest := enforcer.NewAuthorizationUpdateRequest()
	updateRequest.SetRoleMembers(enforcer.MLPProjectsReaderRole, bootstrapCfg.ProjectReaders)
	updateRequest.SetRoleMembers(enforcer.MLPAdminRole, bootstrapCfg.MLPAdmins)
	updateRequest.AddRolePermissions(enforcer.MLPAdminRole, defaultPermissions)
	return authEnforcer.UpdateAuthorization(context.Background(), updateRequest)
}
