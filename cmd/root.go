package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/privateerproj/privateer-sdk/command"
	"github.com/privateerproj/privateer-sdk/plugin"
	"github.com/privateerproj/privateer-sdk/raidengine"

	"github.com/krumIO/raid-rds/strikes"
)

var (
	// Build information is added by the Makefile at compile time
	buildVersion       string
	buildGitCommitHash string
	buildTime          string

	RaidName = "RDS" // TODO: Change this to the name of your Raid

	// runCmd represents the base command when called without any subcommands
	runCmd = &cobra.Command{
		Use:   RaidName,
		Short: "This Raid evaluates RDS",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			command.InitializeConfig()
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Serve plugin
			raid := &Raid{}
			serveOpts := &plugin.ServeOpts{
				Plugin: raid,
			}

			plugin.Serve(RaidName, serveOpts)
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the runCmd.
func Execute(version, commitHash, builtAt string) {
	buildVersion = version
	buildGitCommitHash = commitHash
	buildTime = builtAt

	err := runCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	command.SetBase(runCmd) // This initializes the base CLI functionality
	// Todo: Add any additional flags/options here
	viper.BindPFlag("raids.rds.tactic", runCmd.PersistentFlags().Lookup("tactic"))
}

// Raid meets the Privateer Service Pack interface
type Raid struct {
}

// cleanupFunc is called when the plugin is stopped
func cleanupFunc() error {
	// Todo: Cleanup seeded data and connections
	// Todo: add flag to optionally keep seeded data
	return nil
}

// Start is called from Privateer after the plugin is served
// At minimum, this should call raidengine.Run()
// Adding raidengine.SetupCloseHandler(cleanupFunc) will allow you to append custom cleanup behavior
func (r *Raid) Start() error {
	raidengine.SetupCloseHandler(cleanupFunc)
	return raidengine.Run(RaidName, getStrikes()) // Return errors from strike executions
}

// GetStrikes returns a list of probe objects
func getStrikes() []raidengine.Strike {
	logger := raidengine.GetLogger(RaidName, false)
	a := &strikes.Antijokes{
		Log: logger,
	}
	availableStrikes := map[string][]raidengine.Strike{
		"CCC-Taxonomy": {
			a.SQLFeatures,
			a.VerticalScaling,
			a.Replication,
			a.MultiRegion,
			a.AutomatedBackup,
			a.BackupRecovery,
			a.Encryption,
			a.RBAC,
			a.Logging,
			a.Monitoring,
			a.Alerting,
		},
		"CIS": {
			a.DNE,
		},
	}
	tactic := viper.GetString("raids.rds.tactic")
	strikes := availableStrikes[tactic]
	if len(strikes) == 0 {
		message := fmt.Sprintf("No strikes were found for the provided strike set: %s", tactic)
		logger.Error(message)
	}
	return strikes
}
