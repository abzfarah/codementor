// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


package main

import (
	"os"
	"os/signal"
	"syscall"

	l4g "github.com/alecthomas/log4go"
	"github.com/nomadsingles/platform/api"

	"github.com/nomadsingles/platform/app"


	"github.com/nomadsingles/platform/model"
	"github.com/nomadsingles/platform/utils"
	"github.com/nomadsingles/platform/web"

	"github.com/spf13/cobra"

)

var MaxNotificationsPerChannelDefault int64 = 1000000

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the Nomad server",
	RunE:  runServerCmd,
}

func runServerCmd(cmd *cobra.Command, args []string) error {
	config, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	runServer(config)
	return nil
}

func runServer(configFileLocation string) {
	if errstr := doLoadConfig(configFileLocation); errstr != "" {
		l4g.Exit("Unable to load nomadsingles configuration file: ", errstr)
		return
	}

	utils.InitTranslations(utils.Cfg.LocalizationSettings)


	pwd, _ := os.Getwd()
	l4g.Info(utils.T("nomadsingles.current_version"), model.CurrentVersion, model.BuildNumber, model.BuildDate, model.BuildHash)
	l4g.Info(utils.T("nomadsingles.entreprise_enabled"), model.BuildEnterpriseReady)
	l4g.Info(utils.T("nomadsingles.working_dir"), pwd)
	l4g.Info(utils.T("nomadsingles.config_file"), utils.FindConfigFile(configFileLocation))

	// Enable developer settings if this is a "dev" build
	if model.BuildNumber == "dev" {
		*utils.Cfg.ServiceSettings.EnableDeveloper = true
	}



	app.NewServer()
	app.InitStores()
	api.InitRouter()

	api.InitApi()

	web.InitWeb()



	app.StartServer()



	// wait for kill signal before attempting to gracefully shutdown
	// the running service
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c


	app.StopServer()
}
