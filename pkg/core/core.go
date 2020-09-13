package core

import (
	"errors"
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

func (conf *Config) validateBackupRequest() error {
	if conf.aws.IAMAccessKey == "" {
		return errors.New("missing AWS IAM access key")
	}

	if conf.aws.IAMSecretKey == "" {
		return errors.New("missing AWS IAM secret key")
	}

	if conf.aws.Region == "" {
		return errors.New("missing AWS region")
	}

	if conf.aws.Bucket == "" {
		return errors.New("missing S3 bucket name")
	}

	if conf.dg.ExportPath == "" {
		return errors.New("missing Dgraph exports path")
	}

	if conf.dg.ExportFormat != "json" && conf.dg.ExportFormat != "rdf" {
		return errors.New("incorrect export format, requried json or rdf")
	}

	if conf.dg.Host == "" {
		return errors.New("missing dgraph host name")
	}

	if conf.dg.HostPort == "" {
		return errors.New("missing dgraph host port")
	}

	return nil
}

// BackupSequence - initiates backup sequence
// 1. Request Dgraph export
// 2. Compress export contents
// 3. Upload to S3 bucket
// 4. Clean export and archive files
func (conf *Config) BackupSequence() error {

	err := conf.validateBackupRequest()
	if err != nil {
		return err
	}

	err = conf.dg.Export()
	if err != nil {
		return err
	}

	archiveName := fmt.Sprintf("./%s-%s.zip",
		conf.dg.ExportFilePrefix,
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

	fmt.Println("Backup successful.")
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
