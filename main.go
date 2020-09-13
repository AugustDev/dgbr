package main

import (
	"fmt"

	"github.com/AugustDev/dgraph-backup-restore/pkg/dgraph"
	"github.com/AugustDev/dgraph-backup-restore/pkg/utils"
)

func main() {
	fmt.Println("Welcome")

	dg := dgraph.Dgraph{
		Hostname:         "localhost",
		HostPort:         "8080",
		ExportPath:       "/Users/August/exports",
		ExportFormat:     "rdf",
		ExportFilePrefix: "am",
	}

	// err := dg.Export()
	// if err != nil {
	// 	fmt.Println("Error exporting")
	// }

	// path, err := dg.Archive()
	// if err != nil {
	// 	fmt.Println("Error archiving")
	// }

	utils.Archive(dg.ExportPath, "output.zip")

	// fmt.Println(path)
}
