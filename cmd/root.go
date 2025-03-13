package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "sshub",
	Short: "SSH connection hub",
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
