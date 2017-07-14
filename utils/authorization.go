// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


package utils

import (
	"github.com/nomadsingles/platform/model"
)

func SetDefaultRolesBasedOnConfig() {
	// Reset the roles to default to make this logic easier
	model.InitalizeRoles()


}
