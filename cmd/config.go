package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func setConfig(cmd *cobra.Command, args []string) {

	fmt.Println("Set config")

}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Add list ssh connections",
	Run:   setConfig,
}

func init() {
	configCmd.Flags().String("list", "l", "List ssh connections")
	rootCmd.AddCommand(configCmd)
}
