package cmd

import (
	"fmt"
	"log"

	"github.com/Hapaa16/sshub/pkg/db"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var deleteConnection = &cobra.Command{
	Use:   "delete",
	Short: "Delete existing connection",
	Run: func(c *cobra.Command, args []string) {
		options := make([]string, 0)

		rows, err := db.DB.Query("SELECT Title FROM connections;")

		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var rs string

			rows.Scan(&rs)

			options = append(options, rs)
		}

		p := promptui.Select{
			Label: "Connection",
			Items: options,
		}

		_, val, err := p.Run()

		if err != nil {
			log.Fatal(err)
		}

		deleteQuery := "DELETE FROM connections WHERE Title = ?"

		_, err = db.DB.Exec(deleteQuery, val)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Successfully delete connection %s\n", val)
	},
}

func init() {
	connectionCmd.AddCommand(deleteConnection)
}
