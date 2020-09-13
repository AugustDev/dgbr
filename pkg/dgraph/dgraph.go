package dgraph

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/AugustDev/dgbr/pkg/dgraph/internal/types"
)

// TempDataFolder - holds directory to extract backups before importing
var TempDataFolder = "data/"

// Config - collects all parameters related to interacting with Dgraph
type Config struct {
	Host             string
	HostPort         string
	ExportPath       string
	ExportFormat     string
	ExportFilePrefix string
	AlphaHost        string
	AlphaPort        string
	ZeroHost         string
	ZeroPort         string
}

// Export - initiates exporting Dgraph
func (dg *Config) Export() error {
	exportURL := fmt.Sprintf("http://%s:%s/admin/export?format=%s", dg.Host, dg.HostPort, dg.ExportFormat)

	res, err := http.Get(exportURL)
	if err != nil {
		return err
	}

	response := types.ExportResponse{}
	json.NewDecoder(res.Body).Decode(&response)

	if len(response.Errors) > 0 {
		return errors.New(response.Errors[0].Extensions.Code)
	}

	return nil
}

// Restore - initiates importing backup to Dgraph
func (dg *Config) Restore(restorePath string, schemaPath string) error {
	alphaAddr := fmt.Sprintf("%s:%s", dg.AlphaHost, dg.AlphaPort)
	zeroAddr := fmt.Sprintf("%s:%s", dg.ZeroHost, dg.ZeroPort)

	log.Println("COMMAND", "dgraph", "live", "-f", TempDataFolder, "-a", alphaAddr, "-z", zeroAddr, "-s", schemaPath)
	cmd := exec.Command("dgraph", "live", "-f", TempDataFolder, "-a", alphaAddr, "-z", zeroAddr, "-s", schemaPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		return err
	}
	io.WriteString(stdin, "4\n")
	cmd.Wait()

	return nil
}
