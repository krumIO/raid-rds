package strikes

import (
	"fmt"
	"log"

	"github.com/privateerproj/privateer-sdk/raidengine"
)

// ToDo: Features to implement
// SQLFeatures.go
// VerticalScaling.go
// Replication.go
// MultiRegion.go
// AutomatedBackup.go
// BackupRecovery.go
// Encryption.go
// RBAC.go
// Logging.go
// Monitoring.go
// Alerting.go

// This creates a database table
func (a *RDSParams) CreateTable() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "Create Table"
	result = raidengine.StrikeResult{
		Passed:  false,
		Message: "",
		DocsURL: "",
		// Future:
		// Movements: []raidengine.Movement{},
	}
	// Movement
	log.Printf("Create a Table")
	name, err := getDBURL()
	if err != nil {
		return
	}
	log.Printf("Creating table [table_name] in database [database_name] on host [host_name]")
	log.Printf(fmt.Sprintf("Preparing %s table", name))

	table_name := "jokes"
	// Demo the log timestamp
	for i := 1; i < 5000000; i++ {
		if i%500000 == 0 {
			log.Printf("Executing create table %s query", table_name)
		}
	}
	log.Printf("Table %s created", table_name)
	result.Passed = true
	return
}

// Todo: Movement for a simple sql read query
// Todo: Movement for a simple sql write query
// Todo: Movement for a simple table delete query
