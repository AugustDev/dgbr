package dgraph

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AugustDev/dgraph-backup-restore/pkg/dgraph/internal/types"
)

// Config - collects all parameters related to interacting with Dgraph
type Config struct {
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
func (dg *Config) Export() error {
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
