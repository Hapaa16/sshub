package main

import (
	"github.com/Hapaa16/sshub/cmd"
	"github.com/Hapaa16/sshub/pkg/db"
)

func main() {

	db.NewSQlclient()
	cmd.Execute()

}
