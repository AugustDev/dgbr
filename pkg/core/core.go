package core

import (
	"fmt"
	"time"

	"github.com/AugustDev/dgraph-backup-restore/pkg/awsx"
	"github.com/AugustDev/dgraph-backup-restore/pkg/dgraph"
	"github.com/AugustDev/dgraph-backup-restore/pkg/utils"
)

// Config - collects all configuration variables into unified struct
type Config struct {
	dg  dgraph.Config
	aws awsx.Config
}

// New - returns initializer containing unified config variables
func New(dg dgraph.Config, aws awsx.Config) Config {
	return Config{
		dg:  dg,
		aws: aws,
	}
}

// BackupSequence - initiates backup sequence
// 1. Request Dgraph export
// 2. Compress export contents
// 3. Upload to S3 bucket
// 4. Clean export and archive files
func (conf *Config) BackupSequence() error {

	err := conf.dg.Export()
	if err != nil {
		return err
	}

	archiveName := fmt.Sprintf("./%s-%s-%s.zip",
		conf.dg.ExportFilePrefix,
		conf.dg.Hostname,
		time.Now().Format(time.RFC3339),
	)

	err = utils.Archive(conf.dg.ExportPath, archiveName)
	if err != nil {
		return err
	}

	err = conf.aws.UploadToS3(archiveName)
	if err != nil {
		return err
	}

	err = utils.Clean(conf.dg.ExportPath, archiveName)
	if err != nil {
		return err
	}

	return nil
}

// RestoreSequence - initiates restore sequence
// 1. Download backup from S3 bucket
// 2. Unarchive to temporary folder
// 3. Obtain schema
// 4. Perform dgraph import using live loader
// 5. Cleans temporary folder
func (conf *Config) RestoreSequence() error {

	filepath, err := conf.aws.DownloadFromS3("aa-localhost-2020-09-13T08:52:59+01:00.zip")
	if err != nil {
		fmt.Println(err)
	}

	restorePath, err := utils.Unarchive(filepath)
	if err != nil {
		return err
	}

	schemaPath, err := utils.GetSchemaPath(dgraph.TempDataFolder)
	if err != nil {
		return err
	}

	err = conf.dg.Restore(restorePath, schemaPath)
	if err != nil {
		return err
	}

	err = utils.CleanAfterRestore()
	if err != nil {
		return err
	}

	return nil
}
