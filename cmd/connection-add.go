package cmd

import (
	"fmt"
	"log"

	"github.com/Hapaa16/sshub/pkg/db"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var cAdd = &cobra.Command{
	Use:   "add",
	Short: "Save new ssh connection config",
	// Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sh := SSHClient{}

		titlePrompt := promptui.Prompt{
			Label: "Connection name",
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("title cannot be empty")
				}
				return nil
			},
		}

		title, err := titlePrompt.Run()

		hostPrompt := promptui.Prompt{
			Label: "Host",
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("host name cannot be empty")
				}
				return nil
			},
		}

		if err != nil {
			log.Fatal(err)
		}

		host, err := hostPrompt.Run()

		if err != nil {
			log.Fatal(err)
		}

		usernamePrompt := promptui.Prompt{
			Label:   "Username",
			Default: "ubuntu",
		}

		username, err := usernamePrompt.Run()

		if err != nil {
			log.Fatal(err)
		}

		sh.Title = title
		sh.Host = host
		sh.Username = username

		connectionTypePrompt := promptui.Select{
			Label: "Connection type",
			Items: []string{"Password", "Private key"},
		}

		_, connectionType, err := connectionTypePrompt.Run()

		if err != nil {
			log.Fatal(err)
		}

		if connectionType == "Password" {
			sh.getPassword()
		} else {
			sh.getPrivateKeyPath()
		}

		insertSQL := `INSERT INTO connections (Title,Host,Username,IsPasswordConnection,Password,PemLocation) VALUES (?, ?, ?, ?, ?, ?)`

		_, err = db.DB.Exec(insertSQL, sh.Title, sh.Host, sh.Username, sh.IsPasswordConnection, sh.Password, sh.PrivateKeyPath)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Successfully added new connection ✅")
	},
}

func init() {
	connectionCmd.AddCommand(cAdd)

}
