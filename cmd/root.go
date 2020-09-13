/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var awsAccessKey string
var awsSecretKey string
var region string
var bucket string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dgbr",
	Short: "Dgraph community backup restore tool",
	Long: `Dgraph community backup restore tool:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dgraph-backup-restore.yaml)")
	rootCmd.PersistentFlags().StringVar(&awsAccessKey, "AWS_ACCESS_KEY", "", "AWS access key for your IAM user (required)")
	rootCmd.PersistentFlags().StringVar(&awsSecretKey, "AWS_SECRET_KEY", "", "AWS secret key for your IAM user (requried)")
	rootCmd.PersistentFlags().StringVar(&region, "region", "", "Your AWS S3 bucket region (requried)")
	rootCmd.PersistentFlags().StringVar(&bucket, "bucket", "", "Name of your AWS S3 bucket (requried)")

	// does not issue error for some reason :shrug:
	rootCmd.MarkFlagRequired("AWS_ACCESS_KEY")
	rootCmd.MarkFlagRequired("AWS_SECRET_KEY")
	rootCmd.MarkFlagRequired("region")
	rootCmd.MarkFlagRequired("bucket")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".dgraph-backup-restore" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dgraph-backup-restore")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
