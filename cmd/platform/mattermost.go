// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


package main

import (

	"os"


	"github.com/spf13/cobra"



	// Enterprise Deps
	_ "github.com/dgryski/dgoogauth"
	_ "github.com/go-ldap/ldap"
	_ "github.com/mattermost/rsc/qr"
)

//ENTERPRISE_IMPORTS

func main() {
	var rootCmd = &cobra.Command{
		Use:   "platform",
		Short: "Dating website",
		Long:  `Nomad Singles is a Somali Dating Website`,
		RunE:  runServerCmd,
	}
	rootCmd.PersistentFlags().StringP("config", "c", "config.json", "Configuration file to use.")


	rootCmd.AddCommand(serverCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

