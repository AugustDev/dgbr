package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/AugustDev/dgraph-backup-restore/pkg/dgraph"
	"github.com/mholt/archiver"
)

// Archive - compresses the contents of the directory
func Archive(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = strings.TrimPrefix(path, source)
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

// Unarchive - unarchives compressed backup
func Unarchive(source string) (string, error) {

	err := archiver.Unarchive(source, dgraph.TempDataFolder)
	if err != nil {
		return "", err
	}

	restoreDir := fmt.Sprintf("%s/%s", dgraph.TempDataFolder, strings.TrimSuffix(source, filepath.Ext(source)))

	return restoreDir, nil
}

// DeleteFolderContents - removes exports folder conents after successful backup
// We don't want to upload old backups again
func deleteFolderContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// Clean - deletes contents of exports directory and created archive after S3 upload
func Clean(exportsPath, archiveName string) error {
	err := deleteFolderContents(exportsPath)
	if err != nil {
		return err
	}

	err = os.Remove(archiveName)
	if err != nil {
		return err
	}

	return nil
}

// GetSchemaPath - returns path to the schema file
// NOTE: assuming there is only one .schema.gz file
func GetSchemaPath(dir string) (schemaPath string, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".schema.gz") {
			schemaPath = path
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return schemaPath, nil
}

// CleanAfterRestore - deletes tempoerary directory after backup restore
func CleanAfterRestore() error {
	err := deleteFolderContents(dgraph.TempDataFolder)
	if err != nil {
		return err
	}

	return nil
}
