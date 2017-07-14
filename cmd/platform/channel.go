// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.

package main

import (
	"errors"


	"github.com/nomadsingles/platform/utils"
	"github.com/spf13/cobra"
)

var channelCmd = &cobra.Command{
	Use:   "channel",
	Short: "Management of channels",
}

var channelCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a channel",
	Long:  `Create a channel.`,
	Example: `  channel create --team myteam --name mynewchannel --display_name "My New Channel"
  channel create --team myteam --name mynewprivatechannel --display_name "My New Private Channel" --private`,
	RunE: createChannelCmdF,
}

var removeChannelUsersCmd = &cobra.Command{
	Use:     "remove [channel] [users]",
	Short:   "Remove users from channel",
	Long:    "Remove some users from channel",
	Example: "  channel remove mychannel user@example.com username",
	RunE:    removeChannelUsersCmdF,
}

var addChannelUsersCmd = &cobra.Command{
	Use:     "add [channel] [users]",
	Short:   "Add users to channel",
	Long:    "Add some users to channel",
	Example: "  channel add mychannel user@example.com username",
	RunE:    addChannelUsersCmdF,
}

var deleteChannelsCmd = &cobra.Command{
	Use:   "delete [channels]",
	Short: "Delete channels",
	Long: `Permanently delete some channels.
Permanently deletes a channel along with all related information including posts from the database.
Channels can be specified by [team]:[channel]. ie. myteam:mychannel or by channel ID.`,
	Example: "  channel delete myteam:mychannel",
	RunE:    deleteChannelsCmdF,
}

var listChannelsCmd = &cobra.Command{
	Use:   "list [teams]",
	Short: "List all channels on specified teams.",
	Long: `List all channels on specified teams.
Archived channels are appended with ' (archived)'.`,
	Example: "  channel list myteam",
	RunE:    listChannelsCmdF,
}

var restoreChannelsCmd = &cobra.Command{
	Use:   "restore [channels]",
	Short: "Restore some channels",
	Long: `Restore a previously deleted channel
Channels can be specified by [team]:[channel]. ie. myteam:mychannel or by channel ID.`,
	Example: "  channel restore myteam:mychannel",
	RunE:    restoreChannelsCmdF,
}

func init() {
	channelCreateCmd.Flags().String("name", "", "Channel Name")
	channelCreateCmd.Flags().String("display_name", "", "Channel Display Name")
	channelCreateCmd.Flags().String("team", "", "Team name or ID")
	channelCreateCmd.Flags().String("header", "", "Channel header")
	channelCreateCmd.Flags().String("purpose", "", "Channel purpose")
	channelCreateCmd.Flags().Bool("private", false, "Create a private channel.")

	channelCmd.AddCommand(
		channelCreateCmd,
		removeChannelUsersCmd,
		addChannelUsersCmd,
		deleteChannelsCmd,
		listChannelsCmd,
		restoreChannelsCmd,
	)
}

func createChannelCmdF(cmd *cobra.Command, args []string) error {
	initDBCommandContextCobra(cmd)

	if !utils.IsLicensed {
		return errors.New(utils.T("cli.license.critical"))
	}

	name, errn := cmd.Flags().GetString("name")
	if errn != nil || name == "" {
		return errors.New("Name is required")
	}
	displayname, errdn := cmd.Flags().GetString("display_name")
	if errdn != nil || displayname == "" {
		return errors.New("Display Name is required")
	}
	teamArg, errteam := cmd.Flags().GetString("team")
	if errteam != nil || teamArg == "" {
		return errors.New("Team is required")
	}

	return nil
}

func removeChannelUsersCmdF(cmd *cobra.Command, args []string) error {
	initDBCommandContextCobra(cmd)

	if !utils.IsLicensed {
		return errors.New(utils.T("cli.license.critical"))
	}

	if len(args) < 2 {
		return errors.New("Not enough arguments.")
	}




	return nil
}



func addChannelUsersCmdF(cmd *cobra.Command, args []string) error {
	initDBCommandContextCobra(cmd)

	if !utils.IsLicensed {
		return errors.New(utils.T("cli.license.critical"))
	}

	if len(args) < 2 {
		return errors.New("Not enough arguments.")
	}






	return nil
}



func deleteChannelsCmdF(cmd *cobra.Command, args []string) error {
	initDBCommandContextCobra(cmd)

	if len(args) < 1 {
		return errors.New("Enter at least one channel to delete.")
	}



	return nil
}

func listChannelsCmdF(cmd *cobra.Command, args []string) error {
	initDBCommandContextCobra(cmd)

	if !utils.IsLicensed {
		return errors.New(utils.T("cli.license.critical"))
	}

	if len(args) < 1 {
		return errors.New("Enter at least one team.")
	}



	return nil
}

func restoreChannelsCmdF(cmd *cobra.Command, args []string) error {
	initDBCommandContextCobra(cmd)

	if !utils.IsLicensed {
		return errors.New(utils.T("cli.license.critical"))
	}

	if len(args) < 1 {
		return errors.New("Enter at least one channel.")
	}



	return nil
}
