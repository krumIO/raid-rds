package strikes

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
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

	rdsClient := rds.NewFromConfig(cfg)
	identifier, _ := getDBInstanceIdentifier()

	input := &rds.DescribeDBInstanceAutomatedBackupsInput{
		DBInstanceIdentifier: aws.String(identifier),
	}

	backups, err := rdsClient.DescribeDBInstanceAutomatedBackups(context.TODO(), input)
	if err != nil {
		result.Message = "Failed to fetch automated backups for instance " + identifier
		result.Passed = false
		return
	}

	var regions []string
	for _, backup := range backups.DBInstanceAutomatedBackups {
		regions = append(regions, *backup.Region)
	}

	// This checks if theres a read replica in a different region
	if len(regions) > 0 {
		hostDBRegion := getRDSRegion()
		for _, region := range regions {
			// region from the instances are in the form of "use2"
			abbreviation, exists := AWS_REGIONS_ABBR[hostDBRegion]
			if exists {
				if region != abbreviation {
					result.Passed = true
					result.Message = "Completed Successfully"
					return
				}
			}

		}
	}

	result.Passed = false
	result.Message = "Multi Region instances not found"
	return

}
