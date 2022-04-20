/*
Copyright © 2022 Nethermind hello.nethermind.io

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/NethermindEth/1click/configs"
	"github.com/NethermindEth/1click/internal/pkg/generate"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "1click",
	Short: "A brief description of your application",
	Long:  `A tool to allow deploying validators with ease.`,
	// TODO: Start the TUI engine in this callback. Default behavior
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Disable completion default cmd
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Persistent flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.1click.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".1click" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(configs.ConfigFileName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
		fmt.Fprintln(os.Stderr, "Config file not found on the path provided nor in the home directory")

		// Generate config file
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		fmt.Printf("Generating config file in the %s directory\n", home)

		err = generate.GenerateConfig(home)
		cobra.CheckErr(err)

		viper.ReadInConfig()
		cobra.CheckErr(err)
	}

	initLogging()
}

/*
initLogging :
This function is responsible for :-
initializing the logging configurations
params :-
none
returns :-
none
*/
func initLogging() {
	var config configs.LogConfig

	log.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		FieldsOrder:     []string{configs.Component},
		TimestampFormat: "2006-01-02 15:04:05 --",
	})

	err := viper.UnmarshalKey("logs", &config)
	if err != nil {
		log.WithField(configs.Component, "Logger Init").Errorf("Unable to decode into struct, %v", err)
		return
	}
	log.WithField(configs.Component, "Logger Init").Infof("Logging configuration: %+v", config)

	level, err := log.ParseLevel(strings.ToLower(config.Level))
	if err != nil {
		log.WithField(configs.Component, "Logger Init").Error(err)
		return
	}
	log.SetLevel(level)
}
