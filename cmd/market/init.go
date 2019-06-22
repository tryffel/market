package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/akamensky/argparse"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
)

type Initializer struct {
	command *argparse.Command
}

func NewInitializer(parser *argparse.Parser) Initializer {
	i := Initializer{
		command: parser.NewCommand("init", "Initialize service by creating admin user after prompt."),
	}
	return i
}

func (i *Initializer) Happened() bool {
	return i.command.Happened()
}

func GetAdminUser() (string, string, error) {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	user, _ := reader.ReadString('\n')
	user = strings.ReplaceAll(user, "\n", "")

	fmt.Println("Enter new password")
	passw, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	fmt.Println("Enter password again")
	passwsecond, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	pass1 := string(passw)
	pass2 := string(passwsecond)

	fmt.Printf("Creating user '%s'\n", user)

	if pass1 != pass2 {
		fmt.Println("Passwords not matching. Aborting.")
		return "", "", errors.New("invalid password")
	}

	return user, pass1, err
}
