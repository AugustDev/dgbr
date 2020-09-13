package core

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/AugustDev/dgbr/pkg/awsx"
	"github.com/AugustDev/dgbr/pkg/dgraph"
	"github.com/AugustDev/dgbr/pkg/utils"
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

	fmt.Println("validating request...")
	err := conf.validateBackupRequest()
	if err != nil {
		return err
	}

	fmt.Println("exporting data...")
	err = conf.dg.Export()
	if err != nil {
		return err
	}

	archiveName := fmt.Sprintf("./%s-%s.zip",
		conf.dg.ExportFilePrefix,
		time.Now().Format(time.RFC3339),
	)

	fmt.Println("archiving backup...")
	err = utils.Archive(conf.dg.ExportPath, archiveName)
	if err != nil {
		return err
	}

	fmt.Println("uploading backup...")
	err = conf.aws.UploadToS3(archiveName)
	if err != nil {
		return err
	}

	fmt.Println("cleaning up...")
	err = utils.Clean(conf.dg.ExportPath, archiveName)
	if err != nil {
		return err
	}

	fmt.Println("Backup successful.")
	return nil
}

func (conf *Config) validateRestoreRequest() error {
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

	if conf.dg.AlphaHost == "" {
		return errors.New("missing Dgraph alpha host")
	}

	if conf.dg.AlphaHost == "" {
		return errors.New("missing Dgraph alpha port")
	}

	if conf.dg.ZeroHost == "" {
		return errors.New("missing Dgraph zero host")
	}

	if conf.dg.ZeroPort == "" {
		return errors.New("missing Dgraph zero port")
	}

	return nil
}

// RestoreSequence - initiates restore sequence
// 1. Download backup from S3 bucket
// 2. Unarchive to temporary folder
// 3. Obtain schema
// 4. Perform dgraph import using live loader
// 5. Cleans temporary folder
func (conf *Config) RestoreSequence(filename string) error {

	fmt.Println("Validating config...")
	err := conf.validateRestoreRequest()
	if err != nil {
		return err
	}

	fmt.Println("Downloading backup from S3...")
	filepath, err := conf.aws.DownloadFromS3(filename)
	if err != nil {
		return err
	}

	fmt.Println("Unarchiving backup...")
	restorePath, err := utils.Unarchive(filepath)
	if err != nil {
		return err
	}

	fmt.Println("Finding schema...")
	schemaPath, err := utils.GetSchemaPath(dgraph.TempDataFolder)
	if err != nil {
		return err
	}

	fmt.Println("Initiating Dgrpah restore...")
	err = conf.dg.Restore(restorePath, schemaPath)
	if err != nil {
		return err
	}

	fmt.Println("Cleaning ...")
	err = utils.CleanAfterRestore()
	if err != nil {
		return err
	}

	_ = os.Remove(filepath)

	fmt.Println("Restore successful.")
	return nil
}
