package main

import (
	"fmt"

	"github.com/AugustDev/dgraph-backup-restore/pkg/dgraph"
)

func main() {
	fmt.Println("Welcome")

	dg := dgraph.Dgraph{
		Hostname:     "localhost",
		HostPort:     "8080",
		ExportPath:   "/exports",
		ExportFormat: "rdf",
	}

	err := dg.Export()
	if err != nil {
		fmt.Println("Error exporting")
	}
}
