package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCommand(rt *Runtime) *cobra.Command {
	rt.ApplyEnvFallbacks()
	var configFile string

	root := &cobra.Command{
		Use:           "proxctl",
		Short:         "Lightweight proxy deployment and operations CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       Version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if configFile != "" {
				viper.SetConfigFile(configFile)
			} else {
				viper.SetConfigName("defaults")
				viper.SetConfigType("yaml")
				viper.AddConfigPath("/etc/proxctl")
				viper.AddConfigPath("./configs")
			}
			viper.SetEnvPrefix("PROXCTL")
			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
			viper.AutomaticEnv()
			_ = viper.ReadInConfig()
			bindRuntime(rt)
			return nil
		},
	}

	root.PersistentFlags().StringVar(&configFile, "config-file", "", "proxctl defaults file")
	root.PersistentFlags().String("provider", rt.Config.Provider, "default provider")
	root.PersistentFlags().String("xray-config", rt.Config.ConfigPath, "Xray config path")
	root.PersistentFlags().String("state-dir", rt.Config.StateDir, "state directory")
	root.PersistentFlags().String("xray-bin", rt.Config.XrayBin, "Xray binary")
	root.PersistentFlags().String("service", rt.Config.Service, "Xray systemd service")
	root.PersistentFlags().Bool("yes", false, "accept low-risk defaults")
	root.PersistentFlags().Bool("json", false, "emit JSON where supported")
	root.PersistentFlags().Bool("no-color", false, "disable color")
	_ = viper.BindPFlag("provider", root.PersistentFlags().Lookup("provider"))
	_ = viper.BindPFlag("xray.config_path", root.PersistentFlags().Lookup("xray-config"))
	_ = viper.BindPFlag("operations.state_dir", root.PersistentFlags().Lookup("state-dir"))
	_ = viper.BindPFlag("xray.binary", root.PersistentFlags().Lookup("xray-bin"))
	_ = viper.BindPFlag("xray.service", root.PersistentFlags().Lookup("service"))
	_ = viper.BindPFlag("yes", root.PersistentFlags().Lookup("yes"))
	_ = viper.BindPFlag("json", root.PersistentFlags().Lookup("json"))
	_ = viper.BindPFlag("no_color", root.PersistentFlags().Lookup("no-color"))

	root.AddCommand(
		newStatusCommand(rt),
		newAuditCommand(rt),
		newHealthCommand(rt),
		newDoctorCommand(rt),
		newConfigCommand(rt),
		newClientCommand(rt),
		newBackupCommand(rt),
		newRestoreCommand(rt),
		newAdoptCommand(rt),
		newInstallCommand(rt),
		newInitCommand(rt),
		newPlanCommand(rt),
		newApplyCommand(rt),
		newSSHCommand(rt),
		newFirewallCommand(rt),
		newBootCheckCommand(rt),
		newDocsCommand(rt),
		newWizardCommand(rt),
	)

	return root
}

func Execute(ctx context.Context, args []string, rt *Runtime) error {
	cmd := NewRootCommand(rt)
	cmd.SetArgs(args)
	cmd.SetOut(rt.Out)
	cmd.SetErr(rt.Err)
	return cmd.ExecuteContext(ctx)
}

func bindRuntime(rt *Runtime) {
	rt.Config.Provider = viper.GetString("provider")
	rt.Config.ConfigPath = viper.GetString("xray.config_path")
	rt.Config.StateDir = viper.GetString("operations.state_dir")
	rt.Config.XrayBin = viper.GetString("xray.binary")
	rt.Config.Service = viper.GetString("xray.service")
	rt.Config.Yes = viper.GetBool("yes")
	rt.Config.JSON = viper.GetBool("json")
	rt.Config.NoColor = viper.GetBool("no_color")
	if rt.Config.Provider == "" {
		rt.Config.Provider = "xray"
	}
	if rt.Config.ConfigPath == "" {
		rt.Config.ConfigPath = DefaultConfig().ConfigPath
	}
	if rt.Config.StateDir == "" {
		rt.Config.StateDir = DefaultConfig().StateDir
	}
	if rt.Config.XrayBin == "" {
		rt.Config.XrayBin = DefaultConfig().XrayBin
	}
	if rt.Config.Service == "" {
		rt.Config.Service = DefaultConfig().Service
	}
}

func requireProvider(provider string) error {
	if provider != "xray" {
		return fmt.Errorf("unsupported provider: %s", provider)
	}
	return nil
}

func dieIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
