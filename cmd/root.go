package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vs4vijay/lazykubectl/pkg/k8s"
	"github.com/vs4vijay/lazykubectl/pkg/tui"

	"os"

	"github.com/spf13/viper"
)

var (
	cfgFile        string
	kubeConfigFile string
	dryrun         bool
	home           = k8s.Home()
)

var rootCmd = &cobra.Command{
	Use:   "lazykubectl",
	Short: "A Kubernetes Client",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("lazykubectl")

		configData, err := ioutil.ReadFile(kubeConfigFile)
		if err != nil {
			return fmt.Errorf("Error in Reading Config: %v", err)
			os.Exit(1)
		}

		kubeConfig := k8s.KubeConfig{
			Type:     "MANIFEST",
			Manifest: string(configData),
		}

		kubeapi, err := k8s.NewKubeAPI(kubeConfig)
		if err != nil {
			return err
		}

		app, err := tui.NewApp(kubeapi)
		if err != nil {
			return err
		}

		if dryrun {
			kubeapi.DryRun()
		} else {
			app.Start()
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lazykubectl.yaml)")

	flags := rootCmd.Flags()
	flags.StringVarP(&kubeConfigFile, "kubeconfig", "c", filepath.Join(home, ".kube", "config"), "")
	flags.BoolVar(&dryrun, "dryrun", false, "")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(home)
		viper.SetConfigName(".lazykubectl")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
