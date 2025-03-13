package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Hapaa16/sshub/pkg/db"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

const green = "\033[32m"

const reset = "\033[0m" // Reset color after printing

type SSHClient struct {
	ID                   int
	Title                string
	Host                 string
	Username             string
	IsPasswordConnection bool
	Password             string
	PrivateKeyPath       string
}

func connectToSsh(cd *cobra.Command, args []string) {

	d := make(map[string]SSHClient)

	names := make([]string, 0)

	rows, err := db.DB.Query("SELECT id, Title, Host, Username, IsPasswordConnection, Password, PemLocation FROM connections;")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		connection := SSHClient{}

		rows.Scan(&connection.ID, &connection.Title, &connection.Host, &connection.Username, &connection.IsPasswordConnection, &connection.Password, &connection.PrivateKeyPath)

		d[connection.Title] = connection

		names = append(names, connection.Title)

	}

	connectionTypePrompt := promptui.Select{
		Label: "Connections",
		Items: names,
	}
	_, cnames, err := connectionTypePrompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	pem, err := os.Open(d[cnames].PrivateKeyPath)

	if err != nil {
		log.Fatal(err)
	}

	defer pem.Close()

	content, err := io.ReadAll(pem)

	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	signer, err := ssh.ParsePrivateKey(content)

	if err != nil {
		log.Fatalf("Unable to parse private key: %v", err)
	}
	c := d[cnames]

	config := &ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{

			ssh.PublicKeys(signer),
		},

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// Ignore host key verification (NOT recommended for production use)
	}

	serverAddress := fmt.Sprintf("%s:%s", c.Host, "22")

	fmt.Println(green+"\nConnecting to", c.Title, "..."+reset)

	client, err := ssh.Dial("tcp", serverAddress, config)

	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer client.Close()

	fmt.Println(green+"✅ Successfully connected to", c.Host+reset)

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("linux", 80, 40, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}

	//set input and output
	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
	}

	err = session.Wait()
	if err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
}

func (s *SSHClient) getPrivateKeyPath() {
	prompt := promptui.Prompt{
		Label: "Enter file path",
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("file path cannot be empty")
			}
			return nil
		},
	}

	filePath, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v\n", err)
	}
	s.PrivateKeyPath = filePath
}

func (s *SSHClient) getPassword() {
	prompt := promptui.Prompt{
		Label: "Enter a password",

		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("password cannot be empty")
			}
			return nil
		},
		Mask: '*',
	}

	pw, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v\n", err)
	}
	s.IsPasswordConnection = true
	s.Password = pw
}

var connectionCmd = &cobra.Command{
	Use:   "connect",
	Short: "SSH connection",
	// Args:  cobra.MinimumNArgs(1),
	Run: connectToSsh,
}

func init() {

	rootCmd.AddCommand(connectionCmd)
}
