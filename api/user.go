// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.


package api

import (

	"net/http"



	l4g "github.com/alecthomas/log4go"



	"github.com/nomadsingles/platform/utils"
)

func InitUser() {
	l4g.Debug(utils.T("api.user.init.debug"))


}



func login(c *Context, w http.ResponseWriter, r *http.Request) {

}



