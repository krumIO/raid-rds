package strikes

import (
	"fmt"
	"log"

	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// ToDo: Features to implement
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

func (a *Strikes) SetLogger(loggerName string) {
	a.Log = raidengine.GetLogger(loggerName, false)
}

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

func connectToDb() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "The database host must be available and accepting connections",
		Function:    utils.CallerPath(0),
	}
	_, err := getConfig()
	if err != nil {
		result.Message = err.Error()
		return
	}
	result.Passed = true
	return
}

func createTable() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "A table can be created in the database",
		Function:    utils.CallerPath(0),
	}
	name, err := getConfig()
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
