// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


package model

type Permission struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Role struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}


// General permission that encompases all system admin functions
// in the future this could be broken up to allow access to some
// admin functions but not others
var PERMISSION_MANAGE_SYSTEM *Permission

var ROLE_SYSTEM_USER *Role
var ROLE_SYSTEM_ADMIN *Role


var BuiltInRoles map[string]*Role



func RoleIdsToString(roles []string) string {
	output := ""
	for _, role := range roles {
		output += role + ", "
	}

	if output == "" {
		return "[<NO ROLES>]"
	}

	return output[:len(output)-1]
}

func init() {

}
