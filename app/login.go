// Copyright (c) 2017 Mattermost, Inc. All Rights Reserved.


package app

import (

	"net/http"



	"github.com/nomadsingles/platform/model"

)




func GetProtocol(r *http.Request) string {
	if r.Header.Get(model.HEADER_FORWARDED_PROTO) == "https" || r.TLS != nil {
		return "https"
	} else {
		return "http"
	}
}
