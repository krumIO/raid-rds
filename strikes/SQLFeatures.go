package strikes

import (
	"fmt"
	"log"

	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// Todo/Roadmap: Features to evaluate implementing
// SQLFeatures.go
// VerticalScaling.go: not sure if possible
// Replication.go - future: check for master agent replication
// 	write to agent (fail)
// 	write to master (pass)
// 	write to master and read from agent (pass)
// MultiRegion.go - not sure if possible
// AutomatedBackup.go - AWS CLI - check backup interval
// BackupRecovery.go - AWS CLI - check for point in time recovery
// Encryption.go - AWS CLI - check for encryption and for connection with certificate
// RBAC.go - requires precreated users with varying roles
// Logging.go - check for enabled, req API/CLI
// Monitoring.go - check for enabled, req API/CLI
// Alerting.go - check for enabled, req API/CLI

// This creates a database table
func (a *Strikes) SQLFeatures() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "SQLFeatures"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "This will perform various SQL queries",
		DocsURL:     "https://www.github.com/krumIO/raid-rds",
		ControlID:   "CCC-Taxonomy-1",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	// Movement
	connectToDbMovement := connectToDb()
	result.Movements["ConnectToDB"] = connectToDbMovement
	log.Printf("Connect to database")
	if !connectToDbMovement.Passed {
		result.Message = connectToDbMovement.Message
		return
	}

	createTableMovement := createTable()
	result.Movements["CreateTable"] = createTableMovement
	log.Printf("Create a Table")
	if !createTableMovement.Passed {
		result.Message = createTableMovement.Message
		return
	}

	result.Passed = true
	result.Message = "Completed Successfully"
	return
}

func createTable() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "A table can be created in the database",
		Function:    utils.CallerPath(0),
	}
	name, err := getDBConfig()
	if err != nil {
		result.Message = err.Error()
		return
	}

	log.Printf("Creating table [table_name] in database [database_name] on host [host_name]")
	log.Printf(fmt.Sprintf("Preparing table on %s", name))

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
