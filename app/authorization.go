// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


package app

import (


	l4g "github.com/alecthomas/log4go"
	"github.com/nomadsingles/platform/model"
)

func SessionHasPermissionTo(session model.Session, permission *model.Permission) bool {
	return CheckIfRolesGrantPermission(session.GetUserRoles(), permission.Id)
}

func SessionHasPermissionToTeam(session model.Session, teamId string, permission *model.Permission) bool {
	if teamId == "" {
		return false
	}


	return SessionHasPermissionTo(session, permission)
}

func SessionHasPermissionToChannel(session model.Session, channelId string, permission *model.Permission) bool {
	if channelId == "" {
		return false
	}



	return SessionHasPermissionTo(session, permission)
}

func SessionHasPermissionToChannelByPost(session model.Session, postId string, permission *model.Permission) bool {



	return SessionHasPermissionTo(session, permission)
}

func SessionHasPermissionToUser(session model.Session, userId string) bool {
	if userId == "" {
		return false
	}

	if session.UserId == userId {
		return true
	}

	if SessionHasPermissionTo(session, model.PERMISSION_EDIT_OTHER_USERS) {
		return true
	}

	return false
}

func SessionHasPermissionToPost(session model.Session, postId string, permission *model.Permission) bool {
return true
}

func HasPermissionTo(askingUserId string, permission *model.Permission) bool {
	user, err := GetUser(askingUserId)
	if err != nil {
		return false
	}

	roles := user.GetRoles()

	return CheckIfRolesGrantPermission(roles, permission.Id)
}

func HasPermissionToTeam(askingUserId string, teamId string, permission *model.Permission) bool {


	return true
}

func HasPermissionToChannel(askingUserId string, channelId string, permission *model.Permission) bool {
	if channelId == "" || askingUserId == "" {
		return false
	}



	return true
}

func HasPermissionToChannelByPost(askingUserId string, postId string, permission *model.Permission) bool {


	return HasPermissionTo(askingUserId, permission)
}

func HasPermissionToUser(askingUserId string, userId string) bool {
	if askingUserId == userId {
		return true
	}

	if HasPermissionTo(askingUserId, model.PERMISSION_EDIT_OTHER_USERS) {
		return true
	}

	return false
}

func CheckIfRolesGrantPermission(roles []string, permissionId string) bool {
	for _, roleId := range roles {
		if role, ok := model.BuiltInRoles[roleId]; !ok {
			l4g.Debug("Bad role in system " + roleId)
			return false
		} else {
			permissions := role.Permissions
			for _, permission := range permissions {
				if permission == permissionId {
					return true
				}
			}
		}
	}

	return false
}
