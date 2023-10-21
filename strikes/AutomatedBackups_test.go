package strikes

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func TestAutomatedBackup(t *testing.T) {
	viper.AddConfigPath("../")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println("Config file not found...")
		return
	}

	strikes := Strikes{}
	strikeName, result := strikes.AutomatedBackups()

	fmt.Println(strikeName)
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(b))
	fmt.Println()
}
