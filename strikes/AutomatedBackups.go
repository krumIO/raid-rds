package strikes

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// Todo/Roadmap: Features to evaluate implementing
// AutomatedBackup.go - AWS CLI - check backup interval

// This creates a database table
func (a *Strikes) AutomatedBackups() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "AutomatedBackups"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Check for automated backups against the specified RDS instance",
		DocsURL:     "https://www.github.com/krumIO/raid-rds",
		ControlID:   "CCC-Taxonomy-1",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	// Get Configuration
	cfg, err := getAWSConfig()
	if err != nil {
		result.Message = err.Error()
		return
	}

	rdsInstanceMovement := checkRDSInstanceMovement(cfg)
	result.Movements["CheckForDBInstance"] = rdsInstanceMovement
	if !rdsInstanceMovement.Passed {
		result.Message = rdsInstanceMovement.Message
		return
	}

	automatedBackupsMovement := checkRDSAutomatedBackupMovement(cfg)
	result.Movements["CheckForDBInstanceAutomatedBackups"] = automatedBackupsMovement
	if !automatedBackupsMovement.Passed {
		result.Message = automatedBackupsMovement.Message
		return
	}

	result.Passed = true
	result.Message = "Automated Backups are enabled"
	return
}

func checkRDSAutomatedBackupMovement(cfg aws.Config) (result raidengine.MovementResult) {

	result = raidengine.MovementResult{
		Description: "Check whether the instance has automated backups enabled",
		Function:    utils.CallerPath(0),
	}

	rdsClient := rds.NewFromConfig(cfg)
	instanceIdentifier, _ := getHostDBInstanceIdentifier()

	input := &rds.DescribeDBInstanceAutomatedBackupsInput{
		DBInstanceIdentifier: aws.String(instanceIdentifier),
	}

	backups, err := rdsClient.DescribeDBInstanceAutomatedBackups(context.TODO(), input)
	if err != nil {
		result.Message = err.Error()
		result.Passed = false
		return
	}

	// Loop through the instances and print information
	backupCount := len(backups.DBInstanceAutomatedBackups)
	result.Message = fmt.Sprintf("%d Automated backups found", backupCount)
	result.Passed = backupCount > 0
	return
}
