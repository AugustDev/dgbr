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
	dg   dgraph.Config
	awsx awsx.Config
}

// New - returns initializer containing unified config variables
func New(dg dgraph.Config, awsx awsx.Config) Config {
	return Config{
		dg:   dg,
		awsx: awsx,
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

	err = conf.awsx.UploadToS3(archiveName)
	if err != nil {
		return err
	}

	err = utils.Clean(conf.dg.ExportPath, archiveName)
	if err != nil {
		return err
	}

	return nil
}
