package main

import (
	"fmt"

	"github.com/nomadsingles/platform/app"

	"github.com/nomadsingles/platform/utils"
	"github.com/spf13/cobra"
)

func doLoadConfig(filename string) (err string) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Sprintf("%v", r)
		}
	}()
	utils.TranslationsPreInit()
	utils.EnableConfigFromEnviromentVars()
	utils.LoadConfig(filename)
	utils.InitializeConfigWatch()
	utils.EnableConfigWatch()
	return ""
}

func initDBCommandContextCobra(cmd *cobra.Command) error {
	config, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}
	initDBCommandContext(config)

	return nil
}

func initDBCommandContext(configFileLocation string) {
	if errstr := doLoadConfig(configFileLocation); errstr != "" {
		return
	}

	utils.ConfigureCmdLineLog()

	app.NewServer()
	app.InitStores()

}
