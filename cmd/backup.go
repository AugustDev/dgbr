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

	"github.com/AugustDev/dgraph-backup-restore/pkg/awsx"
	"github.com/AugustDev/dgraph-backup-restore/pkg/core"
	"github.com/AugustDev/dgraph-backup-restore/pkg/dgraph"
	"github.com/spf13/cobra"
)

var exportFormat string
var exportPath string
var exportFilePrefix string
var host string
var hostPort string
var alphaHost string
var alphaPort string
var zeroHost string
var zeroPort string

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Initiates backup to AWS S3 bucket",
	Long: `
Backup command initiates backup to AWS S3 bucket. First tool requests
dgraph to export database and schema to export folder, archives the data, pushes
to S3 bucket and cleans up afterwards.

Example use:

dgbr backup --AWS_ACCESS_KEY=X --AWS_SECRET_KEY=Y --bucket=my-dgraph-backups --region=eu-west-1 --export=/exports`,
	Run: func(cmd *cobra.Command, args []string) {

		awsc := awsx.Config{
			Region:       region,
			Bucket:       bucket,
			IAMAccessKey: awsAccessKey,
			IAMSecretKey: awsSecretKey,
		}

		dg := dgraph.Config{
			Host:             host,
			HostPort:         hostPort,
			ExportPath:       exportPath,
			ExportFormat:     exportFormat,
			ExportFilePrefix: exportFilePrefix,
			AlphaHost:        alphaHost,
			AlphaPort:        alphaPort,
			ZeroHost:         zeroHost,
			ZeroPort:         zeroPort,
		}

		conf := core.New(dg, awsc)

		err := conf.BackupSequence()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {
	backupCmd.Flags().StringVar(&exportFormat, "format", "rdf", "Export format rdf or json")
	backupCmd.Flags().StringVar(&exportPath, "export", "/exports", "Location where dgraph alpha is exporting data (required)")
	backupCmd.Flags().StringVar(&exportFilePrefix, "prefix", "dg", "Prefix appended to backup file name")
	backupCmd.Flags().StringVar(&host, "host", "localhost", "Hostname of running instance of Dgraph server")
	backupCmd.Flags().StringVar(&hostPort, "port", "8080", "Hostname port")
	backupCmd.Flags().StringVar(&alphaHost, "alphaHost", "localhost", "alpha server host name")
	backupCmd.Flags().StringVar(&alphaPort, "alphaPort", "9080", "alpha server port")
	backupCmd.Flags().StringVar(&alphaHost, "zeroHost", "localhost", "zero server host name")
	backupCmd.Flags().StringVar(&alphaPort, "zeroPort", "5080", "zero server port")

	backupCmd.MarkFlagRequired("export")

	rootCmd.AddCommand(backupCmd)
}
