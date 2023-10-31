package strikes

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
	"github.com/spf13/viper"
)

type Strikes struct {
	Log hclog.Logger
}

type Movement struct {
	Strike string
}

func (a *Strikes) SetLogger(loggerName string) {
	a.Log = raidengine.GetLogger(loggerName, false)
}

func getDBConfig() (string, error) {
	err := checkConfigValues([]string{
		"raids.rds.config.host",
		"raids.rds.config.database",
	})
	if err != nil {
		return "", err
	}
	return "database_host_placeholder", nil
}

func getHostDBInstanceIdentifier() (string, error) {
	id := viper.GetString("raids.rds.config.instance_identifier")
	err := checkConfigValues([]string{
		"raids.rds.config.instance_identifier",
	})
	return id, err // id will be "" if not set, err will be nil if id is set
}

func getHostRDSRegion() (string, error) {
	region := viper.GetString("raids.rds.config.primary_region")
	err := checkConfigValues([]string{
		"raids.rds.config.primary_region",
	})
	return region, err // region will be "" if not set, err will be nil if region is set
}

func getAWSConfig() (cfg aws.Config, err error) {
	err = checkConfigValues([]string{
		"aws.access_key",
		"aws.secret_key",
		"aws.session_key",
		"aws.region",
	})
	if err != nil {
		return
	}

	access_key := viper.GetString("aws.access_key")
	secret_key := viper.GetString("aws.secret_key")
	session_key := viper.GetString("aws.session_key")
	region := viper.GetString("aws.region")

	creds := credentials.NewStaticCredentialsProvider(access_key, secret_key, session_key)
	cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds), config.WithRegion(region))
	return
}

func connectToDb() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "The database host must be available and accepting connections",
		Function:    utils.CallerPath(0),
	}
	_, err := getDBConfig()
	if err != nil {
		result.Message = err.Error()
		return
	}
	result.Passed = true
	return
}

func checkRDSInstanceMovement(cfg aws.Config) (result raidengine.MovementResult) {
	// check if the instance is available
	result = raidengine.MovementResult{
		Description: "Check whether the instance can be reached",
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
	result.Message = "Instance found"
	result.Passed = len(instance.DBInstances) > 0
	return
}

func getRDSInstanceFromIdentifier(cfg aws.Config, identifier string) (instance *rds.DescribeDBInstancesOutput, err error) {
	rdsClient := rds.NewFromConfig(cfg)

	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(identifier),
	}

	instance, err = rdsClient.DescribeDBInstances(context.TODO(), input)
	return
}

// TODO: This could be a good addition to the SDK for future raids to use
func checkConfigValues(config_values []string) (err error) {
	missing_values := []string{}
	for _, value := range config_values {
		if !viper.IsSet(value) {
			missing_values = append(missing_values, value)
		}
	}
	if len(missing_values) > 0 {
		err = errors.New("Missing config values: " + strings.Join(missing_values, ", "))
		return
	}
	return
}
