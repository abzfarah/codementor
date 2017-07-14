// Copyright (c) 2017 Mattermost, Inc. All Rights Reserved.


package app

import (
	"fmt"
	"net/http"
	"time"


	"github.com/nomadsingles/platform/model"
	"github.com/nomadsingles/platform/utils"
	"github.com/mssola/user_agent"
)



func DoLogin(w http.ResponseWriter, r *http.Request, user *model.User, deviceId string) (*model.Session, *model.AppError) {
	session := &model.Session{UserId: user.Id, Roles: user.GetRawRoles(), DeviceId: deviceId, IsOAuth: false}

	maxAge := *utils.Cfg.ServiceSettings.SessionLengthWebInDays * 60 * 60 * 24

	if len(deviceId) > 0 {
		session.SetExpireInDays(*utils.Cfg.ServiceSettings.SessionLengthMobileInDays)


	} else {
		session.SetExpireInDays(*utils.Cfg.ServiceSettings.SessionLengthWebInDays)
	}

	ua := user_agent.New(r.UserAgent())

	plat := ua.Platform()
	if plat == "" {
		plat = "unknown"
	}

	os := ua.OS()
	if os == "" {
		os = "unknown"
	}

	bname, bversion := ua.Browser()
	if bname == "" {
		bname = "unknown"
	}

	if bversion == "" {
		bversion = "0.0"
	}

	session.AddProp(model.SESSION_PROP_PLATFORM, plat)
	session.AddProp(model.SESSION_PROP_OS, os)
	session.AddProp(model.SESSION_PROP_BROWSER, fmt.Sprintf("%v/%v", bname, bversion))




	w.Header().Set(model.HEADER_TOKEN, session.Token)

	secure := false
	if GetProtocol(r) == "https" {
		secure = true
	}

	expiresAt := time.Unix(model.GetMillis()/1000+int64(maxAge), 0)
	sessionCookie := &http.Cookie{
		Name:     model.SESSION_COOKIE_TOKEN,
		Value:    session.Token,
		Path:     "/",
		MaxAge:   maxAge,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   secure,
	}

	http.SetCookie(w, sessionCookie)

	return session, nil
}

func GetProtocol(r *http.Request) string {
	if r.Header.Get(model.HEADER_FORWARDED_PROTO) == "https" || r.TLS != nil {
		return "https"
	} else {
		return "http"
	}
}
