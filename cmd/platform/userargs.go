// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.

package main

import (
	"github.com/nomadsingles/platform/app"
	"github.com/nomadsingles/platform/model"
)

func getUsersFromUserArgs(userArgs []string) []*model.User {
	users := make([]*model.User, 0, len(userArgs))
	for _, userArg := range userArgs {
		user := getUserFromUserArg(userArg)
		users = append(users, user)
	}
	return users
}

func getUserFromUserArg(userArg string) *model.User {
	var user *model.User
	if result := <-app.Srv.Store.User().GetByEmail(userArg); result.Err == nil {
		user = result.Data.(*model.User)
	}

	if user == nil {
		if result := <-app.Srv.Store.User().GetByUsername(userArg); result.Err == nil {
			user = result.Data.(*model.User)
		}
	}

	if user == nil {
		if result := <-app.Srv.Store.User().Get(userArg); result.Err == nil {
			user = result.Data.(*model.User)
		}
	}

	return user
}
