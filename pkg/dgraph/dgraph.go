package dgraph

import (
	"compress/flate"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AugustDev/dgraph-backup-restore/pkg/dgraph/internal/types"
	"github.com/mholt/archiver"
)

// Dgraph - collects all parameters related to interacting with Dgraph
type Dgraph struct {
	Hostname         string
	HostPort         string
	ExportPath       string
	ExportFormat     string
	ExportFilePrefix string
}

// func (dg *Dgraph) validate() error {
// 	if dg.Hostname == nil {
// 		return errors.New("dgraph hostname is missing")
// 	}

// 	if dg.ExportFormat == nil {
// 		return errors.New("dgraph export path is missing")
// 	}

// 	if dg.ExportFormat == nil {
// 		return errors.New("dgraph export format is missing")
// 	}

// 	return nil
// }

// Export - initiates exporting Dgraph
func (dg *Dgraph) Export() error {
	exportURL := fmt.Sprintf("http://%s:%s/admin/export?format=%s", dg.Hostname, dg.HostPort, dg.ExportFormat)

	res, err := http.Get(exportURL)
	if err != nil {
		return err
	}

	response := types.ExportResponse{}
	json.NewDecoder(res.Body).Decode(&response)

	if response.Code == "Success" {
		return nil
	}

	return errors.New(response.Message)
}

func (dg *Dgraph) Archive() (filePath string, err error) {
	z := archiver.Zip{
		CompressionLevel: flate.DefaultCompression,
	}

	archiveName := fmt.Sprintf("./%s-%s-%s.zip",
		dg.ExportFilePrefix,
		dg.Hostname,
		time.Now().Format(time.RFC3339),
	)

	fmt.Println(archiveName)

	err = z.Archive([]string{dg.ExportPath}, archiveName)
	if err != nil {
		log.Fatal("err Zipping", err)
		return "", err
	}
	return archiveName, nil
}
