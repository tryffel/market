package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/market/cmd"
	"github.com/tryffel/market/config"
	"github.com/tryffel/market/modules/Error"
	"github.com/tryffel/market/storage/models"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logFormat := &prefixed.TextFormatter{
		ForceColors:    true,
		FullTimestamp:  true,
		QuoteCharacter: "'",
	}
	logFormat.ForceFormatting = true
	logrus.SetFormatter(logFormat)

	parser := argparse.NewParser("Market", "Market document storage for never "+
		"losing files and documents again. For more help see --help. "+
		"For basic usage, run market -c <config-file>.")
	configFile := parser.String("c", "Config",
		&argparse.Options{Required: false, Help: "Configuration file"})
	createConfig := parser.Flag("n", "new",
		&argparse.Options{Required: false, Help: "Create new configuration file"})

	migrator := NewMigrator(parser)
	init := NewInitializer(parser)

	err := parser.Parse(os.Args)
	if err != nil {
		logrus.Error("Failed to parse input flags")
		os.Exit(1)
	}

	conf := config.Config{}
	if *createConfig == true {
		err = conf.SaveFile(*configFile)
		if err != nil {
			Error.Log(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	c, err := config.ReadConfig(*configFile)
	c.AddDefaults()
	c.SaveFile(*configFile)
	service, err := cmd.NewService(c)
	if err != nil {
		e := Error.Wrap(&err, "failed to start market server")
		Error.Log(e)
		os.Exit(1)
	}

	if init.Happened() {
		user, pass, err := GetAdminUser()

		u := models.User{
			Name:     user,
			IsActive: true,
		}
		err = service.Store.User.Create(&u, pass)
		if err != nil {
			logrus.Error("Failed to create new user", err)
			os.Exit(1)
		} else {

			err := service.Store.Group.AddUserToGroup(u.Nid, 1)
			if err != nil {
				logrus.Error(err)
			}

			fmt.Println("New user created successfully")
		}

		os.Exit(0)
	}

	err = migrator.RunMigrations(service.Store.GetDbEngine())
	if err != nil {
		Error.Log(err)
		os.Exit(1)
	}
	service.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT)
	<-ch
	service.Stop()

	c.Logging.Directory = "2"

}
