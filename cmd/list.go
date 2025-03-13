package cmd

import (
	"fmt"
	"log"

	"github.com/Hapaa16/sshub/pkg/db"
	"github.com/spf13/cobra"
)

var connectionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available connections",
	// Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rows, err := db.DB.Query("SELECT id, Title FROM connections;")

		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		for rows.Next() {

			connection := SSHClient{}

			rows.Scan(&connection.ID, &connection.Title)

			message := "ID: %d | Name: %s\n"

			fmt.Printf(message, connection.ID, connection.Title)

		}
	},
}

func init() {
	connectionCmd.AddCommand(connectionListCmd)
}
