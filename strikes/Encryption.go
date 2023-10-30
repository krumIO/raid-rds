package strikes

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// Todo/Roadmap: Features to evaluate implementing
// Encryption.go - AWS CLI

// This creates a database table
func (a *Strikes) Encryption() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "Encryption"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Check if storage is encrypted on the specified RDS instance",
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

	storageEncryptedMovement := checkIfStorageIsEncryptedMovement(cfg)
	result.Movements["CheckForStorageEncryption"] = storageEncryptedMovement
	if !storageEncryptedMovement.Passed {
		result.Message = storageEncryptedMovement.Message
		return
	}

	result.Passed = true
	result.Message = "Completed Successfully"
	return
}

func checkIfStorageIsEncryptedMovement(cfg aws.Config) (result raidengine.MovementResult) {

	result = raidengine.MovementResult{
		Description: "Check if the instance has storage encryption enabled",
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

	if !instance.DBInstances[0].StorageEncrypted {
		result.Message = "Storage encryption is not enabled"
		result.Passed = false
		return
	}

	// Loop through the instances and print information
	result.Passed = true
	return
}
