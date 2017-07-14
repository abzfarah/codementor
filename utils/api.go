// Copyright (c) 2017 Mattermost, Inc. All Rights Reserved.


package utils

import (
	"net/http"
	"strings"
)

type OriginCheckerProc func(*http.Request) bool

func OriginChecker(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	return *Cfg.ServiceSettings.AllowCorsFrom == "*" || strings.Contains(origin, *Cfg.ServiceSettings.AllowCorsFrom)
}

func GetOriginChecker(r *http.Request) OriginCheckerProc {
	if len(*Cfg.ServiceSettings.AllowCorsFrom) > 0 {
		return OriginChecker
	}

	return nil
}
