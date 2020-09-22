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
