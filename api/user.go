// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.


package api

import (
	"bytes"


	"io"
	"net/http"


	"time"

	l4g "github.com/alecthomas/log4go"

	"github.com/nomadsingles/platform/app"

	"github.com/nomadsingles/platform/model"
	"github.com/nomadsingles/platform/store"
	"github.com/nomadsingles/platform/utils"
)

func InitUser() {
	l4g.Debug(utils.T("api.user.init.debug"))


}

func createUser(c *Context, w http.ResponseWriter, r *http.Request) {
	user := model.UserFromJson(r.Body)

	if user == nil {
		c.SetInvalidParam("createUser", "user")
		return
	}

	hash := r.URL.Query().Get("h")


	var ruser *model.User
	var err *model.AppError
	if len(hash) > 0 {
		ruser, err = app.CreateUserWithHash(user, hash, r.URL.Query().Get("d"))
	} else {
		ruser, err = app.CreateUserFromSignup(user)
	}

	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(ruser.ToJson()))
}

func login(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)

	id := props["id"]
	loginId := props["login_id"]
	password := props["password"]
	mfaToken := props["token"]
	deviceId := props["device_id"]
	ldapOnly := props["ldap_only"] == "true"

	c.LogAudit("attempt - user_id=" + id + " login_id=" + loginId)
	user, err := app.AuthenticateUserForLogin(id, loginId, password, mfaToken, deviceId, ldapOnly)
	if err != nil {
		c.LogAudit("failure - user_id=" + id + " login_id=" + loginId)
		c.Err = err
		return
	}

	c.LogAuditWithUserId(user.Id, "success")

	doLogin(c, w, r, user, deviceId)
	if c.Err != nil {
		return
	}

	user.Sanitize(map[string]bool{})

	w.Write([]byte(user.ToJson()))
}

func LoginByOAuth(c *Context, w http.ResponseWriter, r *http.Request, service string, userData io.Reader) *model.User {
	buf := bytes.Buffer{}
	buf.ReadFrom(userData)

	authData := ""

	if len(authData) == 0 {
		c.Err = model.NewLocAppError("LoginByOAuth", "api.user.login_by_oauth.parse.app_error",
			map[string]interface{}{"Service": service}, "")
		return nil
	}

	var user *model.User
	var err *model.AppError
	if user, err = app.GetUserByAuth(&authData, service); err != nil {
		if err.Id == store.MISSING_AUTH_ACCOUNT_ERROR {
			if user, err = app.CreateOAuthUser(service, bytes.NewReader(buf.Bytes()), ""); err != nil {
				c.Err = err
				return nil
			}
		}
		c.Err = err
		return nil
	}



	doLogin(c, w, r, user, "")
	if c.Err != nil {
		return nil
	}

	return user
}

// User MUST be authenticated completely before calling Login
func doLogin(c *Context, w http.ResponseWriter, r *http.Request, user *model.User, deviceId string) {
	session, err := app.DoLogin(w, r, user, deviceId)
	if err != nil {
		c.Err = err
		return
	}

	c.Session = *session
}

func revokeSession(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)
	id := props["id"]

	if err := app.RevokeSessionById(id); err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(model.MapToJson(props)))
}

func attachDeviceId(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)

	deviceId := props["device_id"]
	if len(deviceId) == 0 {
		c.SetInvalidParam("attachDevice", "deviceId")
		return
	}

	// A special case where we logout of all other sessions with the same device id
	if err := app.RevokeSessionsForDeviceId(c.Session.UserId, deviceId, c.Session.Id); err != nil {
		c.Err = err
		c.Err.StatusCode = http.StatusInternalServerError
		return
	}

	app.ClearSessionCacheForUser(c.Session.UserId)
	c.Session.SetExpireInDays(*utils.Cfg.ServiceSettings.SessionLengthMobileInDays)

	maxAge := *utils.Cfg.ServiceSettings.SessionLengthMobileInDays * 60 * 60 * 24

	secure := false
	if app.GetProtocol(r) == "https" {
		secure = true
	}

	expiresAt := time.Unix(model.GetMillis()/1000+int64(maxAge), 0)
	sessionCookie := &http.Cookie{
		Name:     model.SESSION_COOKIE_TOKEN,
		Value:    c.Session.Token,
		Path:     "/",
		MaxAge:   maxAge,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   secure,
	}

	http.SetCookie(w, sessionCookie)

	if err := app.AttachDeviceId(c.Session.Id, deviceId, c.Session.ExpiresAt); err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(model.MapToJson(props)))
}

func getSessions(c *Context, w http.ResponseWriter, r *http.Request) {


}

func logout(c *Context, w http.ResponseWriter, r *http.Request) {

}

func Logout(c *Context, w http.ResponseWriter, r *http.Request) {

}

func getMe(c *Context, w http.ResponseWriter, r *http.Request) {


}

func getInitialLoad(c *Context, w http.ResponseWriter, r *http.Request) {


}

func getUser(c *Context, w http.ResponseWriter, r *http.Request) {

}





