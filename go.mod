module github.com/krumIO/raid-rds

go 1.14

require (
	github.com/aws/aws-sdk-go-v2 v1.21.2
	github.com/aws/aws-sdk-go-v2/config v1.19.0
	github.com/aws/aws-sdk-go-v2/credentials v1.13.43
	github.com/aws/aws-sdk-go-v2/service/rds v1.57.0
	github.com/hashicorp/go-hclog v1.2.0
	github.com/privateerproj/privateer-sdk v0.0.6
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.15.0
)

// For Development Only
// replace github.com/privateerproj/privateer-sdk => ../privateer-sdk
