package main

import (
	"fmt"

	"github.com/AugustDev/dgraph-backup-restore/pkg/core"
	"github.com/AugustDev/dgraph-backup-restore/pkg/dgraph"
)

func main() {

	dg := dgraph.Config{
		Hostname:         "localhost",
		HostPort:         "8080",
		ExportPath:       "/Users/amiras/exports/",
		ExportFormat:     "rdf",
		ExportFilePrefix: "aax",
		AlphaHost:        "localhost",
		AlphaPort:        "9080",
		ZeroHost:         "localhost",
		ZeroPort:         "5080",
	}

	conf := core.New(dg, awsc)

	// conf.BackupSequence()

	err := conf.RestoreSequence()
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(filepath)
}
