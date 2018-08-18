package cmd

import (
	"fmt"
	"os"

	pocket "github.com/brimstone/go-pocket"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	p       *pocket.PocketClient
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pocket",
	Short: "CLI tool for Pocket",
	Long: `A CLI tool for http://getpocket.com

This tool allows for getting items already in pocket.`,
	PersistentPreRunE: setup,
}

func setup(cmd *cobra.Command, args []string) error {
	p = pocket.NewPocketClient(&pocket.PocketClientOptions{
		ConsumerKey: viper.GetString("key"),
		AccessToken: viper.GetString("token"),
	})
	/*
		status := make(chan string)
		go func() {
			url := <-status
			fmt.Printf("Please visit %s\n", url)
		}()
		err := p.Auth(status)
		fmt.Println(p.AccessToken)
		return err
	*/
	return nil
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pocket.yaml)")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".pocket" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pocket")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
