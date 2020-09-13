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

	"github.com/AugustDev/dgbr/pkg/awsx"
	"github.com/AugustDev/dgbr/pkg/core"
	"github.com/AugustDev/dgbr/pkg/dgraph"
	"github.com/spf13/cobra"
)

var restoreName string

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore Dgraph backup from S3 bucket",
	Long: `This function allows you to restore a backup from S3 bucket
by downloding it, extracting and importing to Dgraph using live loader.

Example use:

dgbr restore --AWS_ACCESS_KEY=X --AWS_SECRET_KEY=Y --bucket=my-dgraph-backups --region=eu-west-1 --name=dg_date.zip`,
	Run: func(cmd *cobra.Command, args []string) {

		awsc := awsx.Config{
			Region:       region,
			Bucket:       bucket,
			IAMAccessKey: awsAccessKey,
			IAMSecretKey: awsSecretKey,
		}

		dg := dgraph.Config{
			AlphaHost: alphaHost,
			AlphaPort: alphaPort,
			ZeroHost:  zeroHost,
			ZeroPort:  zeroPort,
		}

		conf := core.New(dg, awsc)

		err := conf.RestoreSequence(restoreName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {
	restoreCmd.Flags().StringVar(&restoreName, "name", "", "backup name for the restore in S3 bucket (required)")
	restoreCmd.Flags().StringVar(&alphaHost, "alphaHost", "localhost", "alpha server host name")
	restoreCmd.Flags().StringVar(&alphaPort, "alphaPort", "9080", "alpha server port")
	restoreCmd.Flags().StringVar(&zeroHost, "zeroHost", "localhost", "zero server host name")
	restoreCmd.Flags().StringVar(&zeroPort, "zeroPort", "5080", "zero server port")
	rootCmd.AddCommand(restoreCmd)
}
