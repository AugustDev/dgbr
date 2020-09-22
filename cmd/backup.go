package cmd

import (
	"fmt"
	"os"

	"github.com/AugustDev/dgbr/pkg/awsx"
	"github.com/AugustDev/dgbr/pkg/core"
	"github.com/AugustDev/dgbr/pkg/dgraph"
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
