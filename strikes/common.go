package strikes

import (
	"errors"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
)

type RDSParams struct {
	Log hclog.Logger
}

type Movement struct {
	Strike string
}

func getDBURL() (string, error) {
	if viper.IsSet("raids.RDS.config") && viper.IsSet("raids.RDS.config.database") {
		return viper.GetString("raids.RDS.config.database.host"), nil
	}
	return "", errors.New("Database URL must be set in the config file (raids.RDS.config.database.host)")
}
