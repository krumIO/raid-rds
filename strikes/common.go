package strikes

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
	"github.com/spf13/viper"
)

var (
	AWS_REGIONS_ABBR = map[string]string{
		"us-east-1":      "use1",
		"us-east-2":      "use2",
		"us-west-1":      "usw1",
		"us-west-2":      "usw2",
		"ca-central-1":   "cac1",
		"eu-west-1":      "euw1",
		"eu-west-2":      "euw2",
		"eu-central-1":   "euc1",
		"eu-north-1":     "eun1",
		"ap-northeast-1": "apne1",
		"ap-northeast-2": "apne2",
		"ap-southeast-1": "apse1",
		"ap-southeast-2": "apse2",
		"ap-south-1":     "aps1",
		"sa-east-1":      "sae1",
	}
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
	if viper.IsSet("raids.RDS.aws.config.host") && viper.IsSet("raids.RDS.aws.config.database") {
		return "database_host_placeholder", nil
	}
	return "", errors.New("database url must be set in the config file")
}

func getDBInstanceIdentifier() (string, error) {
	if viper.IsSet("raids.RDS.aws.config.instance_identifier") {
		return viper.GetString("raids.RDS.aws.config.instance_identifier"), nil
	}
	return "", errors.New("database instance identifier must be set in the config file")
}

func getRDSRegion() string {
	return viper.GetString("raids.RDS.aws.config.region")
}

func getAWSConfig() (cfg aws.Config, err error) {
	if viper.IsSet("raids.RDS.aws.creds") &&
		viper.IsSet("raids.RDS.aws.creds.aws_access_key") &&
		viper.IsSet("raids.RDS.aws.creds.aws_secret_key") &&
		viper.IsSet("raids.RDS.aws.creds.aws_region") {

		access_key := viper.GetString("raids.RDS.aws.creds.aws_access_key")
		secret_key := viper.GetString("raids.RDS.aws.creds.aws_secret_key")
		session_key := viper.GetString("raids.RDS.aws.creds.aws_session_key")
		region := viper.GetString("raids.RDS.aws.creds.aws_region")

		creds := credentials.NewStaticCredentialsProvider(access_key, secret_key, session_key)
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds), config.WithRegion(region))
	}
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
		Description: "Check if the instance is available/exists",
		Function:    utils.CallerPath(0),
	}

	instance, err := getRDSInstance(cfg)
	if err != nil {
		// Handle error
		result.Message = err.Error()
		result.Passed = false
		return
	}
	result.Passed = len(instance.DBInstances) > 0
	return
}

func getRDSInstance(cfg aws.Config) (instance *rds.DescribeDBInstancesOutput, err error) {
	rdsClient := rds.NewFromConfig(cfg)
	identifier, _ := getDBInstanceIdentifier()

	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(identifier),
	}

	instance, err = rdsClient.DescribeDBInstances(context.TODO(), input)
	return
}
