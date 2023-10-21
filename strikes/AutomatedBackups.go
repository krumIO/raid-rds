package strikes

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
	"github.com/spf13/viper"
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

	// Movement
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

	autmatedBackupsMovement := checkRDSAutomatedBackupMovement(cfg)
	result.Movements["CheckForDBInstanceAutomatedBackups"] = autmatedBackupsMovement
	if !autmatedBackupsMovement.Passed {
		result.Message = autmatedBackupsMovement.Message
		return
	}

	result.Passed = true
	result.Message = "Completed Successfully"
	return
}

func checkRDSInstanceMovement(cfg aws.Config) (result raidengine.MovementResult) {
	// check if the instance is available
	result = raidengine.MovementResult{
		Description: "Check if the instance is available/exists",
		Function:    utils.CallerPath(0),
	}

	rdsClient := rds.NewFromConfig(cfg)
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(viper.GetString("raids.RDS.config.aws_db_instance_identifier")),
	}

	instances, err := rdsClient.DescribeDBInstances(context.TODO(), input)
	if err != nil {
		// Handle error
		result.Message = err.Error()
		result.Passed = false
		return
	}
	result.Passed = len(instances.DBInstances) > 0
	return
}

func checkRDSAutomatedBackupMovement(cfg aws.Config) (result raidengine.MovementResult) {

	result = raidengine.MovementResult{
		Description: "Check if the instance has automated backups enabled",
		Function:    utils.CallerPath(0),
	}

	rdsClient := rds.NewFromConfig(cfg)
	input := &rds.DescribeDBInstanceAutomatedBackupsInput{
		DBInstanceIdentifier: aws.String(viper.GetString("raids.RDS.config.aws_db_instance_identifier")),
	}

	backups, err := rdsClient.DescribeDBInstanceAutomatedBackups(context.TODO(), input)
	if err != nil {
		// Handle error
		result.Message = err.Error()
		result.Passed = false
		return
	}

	// Loop through the instances and print information
	result.Passed = len(backups.DBInstanceAutomatedBackups) > 0
	return
}
