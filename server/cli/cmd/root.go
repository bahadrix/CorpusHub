package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var DataPath string

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "CorpusHub CLI",
	Long: `Command Line Interface for CorpusHub`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initDataPath)


	rootCmd.PersistentFlags().StringVar(&DataPath, "data-path", "", "Directory path for data files, (default is $HOME/.corpushub) Automatically created if not exists")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $DATA_PATH/ch.yaml)")

}

// initDataPath reads in config file and ENV variables if set.
func initDataPath() {


	if DataPath == "" {

		homePath, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		DataPath = path.Join(homePath, ".corpushub")
	}

	err := os.MkdirAll(DataPath, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	if cfgFile != ""  {
		viper.SetConfigFile(cfgFile)
	} else {
		cfgFile = path.Join(DataPath, "ch.yaml")

		viper.AddConfigPath(DataPath)
		viper.SetConfigName("ch")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
