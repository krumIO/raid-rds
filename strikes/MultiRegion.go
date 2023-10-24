package strikes

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

func (a *Strikes) MultiRegion() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "MultiRegion"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Check if AWS RDS instance has multi region. This strike only checks for a read replica in a seperate region",
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

	multiRegionMovement := checkRDSMultiRegionMovement(cfg)
	result.Movements["CheckForMultiRegionDBInstances"] = multiRegionMovement
	if !multiRegionMovement.Passed {
		result.Message = multiRegionMovement.Message
		return
	}

	result.Passed = true
	result.Message = "Completed Successfully"

	return
}

func checkRDSMultiRegionMovement(cfg aws.Config) (result raidengine.MovementResult) {

	result = raidengine.MovementResult{
		Description: "Check if the instance has multi region enabled",
		Function:    utils.CallerPath(0),
	}
	instanceIdentifier, _ := getHostDBInstanceIdentifier()

	instance, _ := getRDSInstanceFromIdentifier(cfg, instanceIdentifier)

	// get read replicas from the instance
	readReplicas := instance.DBInstances[0].ReadReplicaDBInstanceIdentifiers

	if len(readReplicas) == 0 {
		result.Passed = false
		result.Message = "Multi Region instances not found"
		return
	}

	hostRDSRegion, _ := getHostRDSRegion()

	// loop over the read replicas and check if they are in a different region
	for _, replica := range readReplicas {
		// we are getting the instance identifier the read replicas
		// get instance from the replica identifier
		replicaInstance, err := getRDSInstanceFromIdentifier(cfg, replica)

		if err != nil {
			result.Passed = false
			result.Message = err.Error()
			return
		}

		if len(replicaInstance.DBInstances) == 0 {
			result.Passed = false
			result.Message = "Cannot access the replica instance " + replica
			return
		}

		// check if replica region matches the host region
		az := *replicaInstance.DBInstances[0].AvailabilityZone
		// db instance doesnt contain the region so we need to remove the last character from the az
		if az[:len(az)-1] == hostRDSRegion {
			result.Passed = false
			result.Message = "Multi Region instances not found"
			return
		}
	}

	result.Passed = true
	return

}
