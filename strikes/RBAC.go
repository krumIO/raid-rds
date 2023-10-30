package strikes

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// Todo/Roadmap: Features to evaluate implementing
// RBAC.go - AWS CLI

// This creates a database table
func (a *Strikes) RBAC() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "RBAC"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Check if database IAM authentication is enabled on the specified RDS instance",
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

	iamDatabaseAuthMovement := checkForIAMDatabaseAuthMovement(cfg)
	result.Movements["CheckForIAMDatabaseAuth"] = iamDatabaseAuthMovement
	if !iamDatabaseAuthMovement.Passed {
		result.Message = iamDatabaseAuthMovement.Message
		return
	}

	result.Passed = true
	result.Message = "Completed Successfully"
	return
}

func checkForIAMDatabaseAuthMovement(cfg aws.Config) (result raidengine.MovementResult) {

	result = raidengine.MovementResult{
		Description: "Check if the instance has IAM Database Authentication enabled",
		Function:    utils.CallerPath(0),
	}

	instanceIdentifier, _ := getHostDBInstanceIdentifier()

	instance, err := getRDSInstanceFromIdentifier(cfg, instanceIdentifier)
	if err != nil {
		// Handle error
		result.Message = err.Error()
		result.Passed = false
		return
	}

	if !instance.DBInstances[0].IAMDatabaseAuthenticationEnabled {
		result.Message = "IAM Database Authentication is not enabled"
		result.Passed = false
		return
	}

	// Loop through the instances and print information
	result.Passed = true
	return
}
