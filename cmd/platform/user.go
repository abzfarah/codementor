// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.

package main

import (
	"errors"


	"github.com/nomadsingles/platform/app"

	"github.com/nomadsingles/platform/model"
	"github.com/nomadsingles/platform/utils"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Management of users",
}

var userActivateCmd = &cobra.Command{
	Use:   "activate [emails, usernames, userIds]",
	Short: "Activate users",
	Long:  "Activate users that have been deactivated.",
	Example: `  user activate user@example.com
  user activate username`,
	RunE: userActivateCmdF,
}

var userDeactivateCmd = &cobra.Command{
	Use:   "deactivate [emails, usernames, userIds]",
	Short: "Deactivate users",
	Long:  "Deactivate users. Deactivated users are immediately logged out of all sessions and are unable to log back in.",
	Example: `  user deactivate user@example.com
  user deactivate username`,
	RunE: userDeactivateCmdF,
}

var userCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a user",
	Long:  "Create a user",
	Example: `  user create --email user@example.com --username userexample --password Password1 
  user create --firstname Joe --system_admin --email joe@example.com --username joe --password Password1`,
	RunE: userCreateCmdF,
}

var userInviteCmd = &cobra.Command{
	Use:   "invite [email] [teams]",
	Short: "Send user an email invite to a team.",
	Long: `Send user an email invite to a team.
You can invite a user to multiple teams by listing them.
You can specify teams by name or ID.`,
	Example: `  user invite user@example.com myteam
  user invite user@example.com myteam1 myteam2`,
	RunE: userInviteCmdF,
}

var resetUserPasswordCmd = &cobra.Command{
	Use:     "password [user] [password]",
	Short:   "Set a user's password",
	Long:    "Set a user's password",
	Example: "  user password user@example.com Password1",
	RunE:    resetUserPasswordCmdF,
}

var resetUserMfaCmd = &cobra.Command{
	Use:   "resetmfa [users]",
	Short: "Turn off MFA",
	Long: `Turn off multi-factor authentication for a user. 
If MFA enforcement is enabled, the user will be forced to re-enable MFA as soon as they login.`,
	Example: "  user resetmfa user@example.com",
	RunE:    resetUserMfaCmdF,
}

var deleteUserCmd = &cobra.Command{
	Use:     "delete [users]",
	Short:   "Delete users and all posts",
	Long:    "Permanently delete user and all related information including posts.",
	Example: "  user delete user@example.com",
	RunE:    deleteUserCmdF,
}

var deleteAllUsersCmd = &cobra.Command{
	Use:     "deleteall",
	Short:   "Delete all users and all posts",
	Long:    "Permanently delete all users and all related information including posts.",
	Example: "  user deleteall",
	RunE:    deleteUserCmdF,
}

var migrateAuthCmd = &cobra.Command{
	Use:   "migrate_auth [from_auth] [to_auth] [match_field]",
	Short: "Mass migrate user accounts authentication type",
	Long: `Migrates accounts from one authentication provider to another. For example, you can upgrade your authentication provider from email to ldap.

from_auth: 
	The authentication service to migrate users accounts from.
	Supported options: email, gitlab, saml. 

to_auth:
	The authentication service to migrate users to.
	Supported options: ldap. 

match_field:
	The field that is guaranteed to be the same in both authentication services. For example, if the users emails are consistent set to email.
	Supported options: email, username.

Will display any accounts that are not migrated successfully.`,
	Example: "  user migrate_auth email ladp email",
	RunE:    migrateAuthCmdF,
}

var verifyUserCmd = &cobra.Command{
	Use:     "verify [users]",
	Short:   "Verify email of users",
	Long:    "Verify the emails of some users.",
	Example: "  user verify user1",
	RunE:    verifyUserCmdF,
}

func init() {
	userCreateCmd.Flags().String("username", "", "Username")
	userCreateCmd.Flags().String("email", "", "Email")
	userCreateCmd.Flags().String("password", "", "Password")
	userCreateCmd.Flags().String("nickname", "", "Nickname")
	userCreateCmd.Flags().String("firstname", "", "First Name")
	userCreateCmd.Flags().String("lastname", "", "Last Name")
	userCreateCmd.Flags().String("locale", "", "Locale (ex: en, fr)")
	userCreateCmd.Flags().Bool("system_admin", false, "Make the user a system administrator")

	deleteUserCmd.Flags().Bool("confirm", false, "Confirm you really want to delete the user and a DB backup has been performed.")

	deleteAllUsersCmd.Flags().Bool("confirm", false, "Confirm you really want to delete the user and a DB backup has been performed.")

	userCmd.AddCommand(
		userActivateCmd,
		userDeactivateCmd,
		userCreateCmd,
		userInviteCmd,
		resetUserPasswordCmd,
		resetUserMfaCmd,
		deleteUserCmd,
		deleteAllUsersCmd,
		migrateAuthCmd,
		verifyUserCmd,
	)
}

func userActivateCmdF(cmd *cobra.Command, args []string) error {
	initDBCommandContextCobra(cmd)

	if len(args) < 1 {
		return errors.New("Enter user(s) to activate.")
	}

	changeUsersActiveStatus(args, true)
	return nil
}

func changeUsersActiveStatus(userArgs []string, active bool) {

}

func changeUserActiveStatus(user *model.User, userArg string, activate bool) {

}

func userDeactivateCmdF(cmd *cobra.Command, args []string) error {
	initDBCommandContextCobra(cmd)

	if len(args) < 1 {
		return errors.New("Enter user(s) to deactivate.")
	}

	changeUsersActiveStatus(args, false)
	return nil
}

func userCreateCmdF(cmd *cobra.Command, args []string) error {
	initDBCommandContextCobra(cmd)
	username, erru := cmd.Flags().GetString("username")
	if erru != nil || username == "" {
		return errors.New("Username is required")
	}
	email, erre := cmd.Flags().GetString("email")
	if erre != nil || email == "" {
		return errors.New("Email is required")
	}
	password, errp := cmd.Flags().GetString("password")
	if errp != nil || password == "" {
		return errors.New("Password is required")
	}
	nickname, _ := cmd.Flags().GetString("nickname")
	firstname, _ := cmd.Flags().GetString("firstname")
	lastname, _ := cmd.Flags().GetString("lastname")
	locale, _ := cmd.Flags().GetString("locale")
	system_admin, _ := cmd.Flags().GetBool("system_admin")

	user := &model.User{
		Username:  username,
		Email:     email,
		Password:  password,
		Nickname:  nickname,
		FirstName: firstname,
		LastName:  lastname,
		Locale:    locale,
	}

	ruser, err := app.CreateUser(user)
	if err != nil {
		return errors.New("Unable to create user. Error: " + err.Error())
	}

	if system_admin {
		app.UpdateUserRoles(ruser.Id, "system_user system_admin")
	}


	return nil
}

func userInviteCmdF(cmd *cobra.Command, args []string) error {
	initDBCommandContextCobra(cmd)
	utils.InitHTML()

	if len(args) < 2 {
		return errors.New("Not enough arguments.")
	}

	email := args[0]
	if !model.IsValidEmail(email) {
		return errors.New("Invalid email")
	}


	return nil
}


func resetUserPasswordCmdF(cmd *cobra.Command, args []string) error {

	return nil
}

func resetUserMfaCmdF(cmd *cobra.Command, args []string) error {


	return nil
}

func deleteUserCmdF(cmd *cobra.Command, args []string) error {

	return nil
}

func deleteAllUsersCommandF(cmd *cobra.Command, args []string) error {
	initDBCommandContextCobra(cmd)

	return nil
}

func migrateAuthCmdF(cmd *cobra.Command, args []string) error {



	return nil
}

func verifyUserCmdF(cmd *cobra.Command, args []string) error {
	initDBCommandContextCobra(cmd)
	if len(args) < 1 {
		return errors.New("Enter at least one user.")
	}



	return nil
}
