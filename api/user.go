// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.


package api

import (
	"bytes"

	"fmt"
	"io"
	"net/http"
	"strconv"

	"time"

	l4g "github.com/alecthomas/log4go"
	"github.com/gorilla/mux"
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

	params := mux.Vars(r)
	id := params["user_id"]

	if !app.SessionHasPermissionToUser(c.Session, id) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	if sessions, err := app.GetSessions(id); err != nil {
		c.Err = err
		return
	} else {
		for _, session := range sessions {
			session.Sanitize()
		}

		w.Write([]byte(model.SessionsToJson(sessions)))
	}
}

func logout(c *Context, w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	data["user_id"] = c.Session.UserId

	Logout(c, w, r)
	if c.Err == nil {
		w.Write([]byte(model.MapToJson(data)))
	}
}

func Logout(c *Context, w http.ResponseWriter, r *http.Request) {
	c.LogAudit("")
	c.RemoveSessionCookie(w, r)
	if c.Session.Id != "" {
		if err := app.RevokeSessionById(c.Session.Id); err != nil {
			c.Err = err
			return
		}
	}
}

func getMe(c *Context, w http.ResponseWriter, r *http.Request) {

	if user, err := app.GetUser(c.Session.UserId); err != nil {
		c.Err = err
		c.RemoveSessionCookie(w, r)
		l4g.Error(utils.T("api.user.get_me.getting.error"), c.Session.UserId)
		return
	} else if HandleEtag(user.Etag(utils.Cfg.PrivacySettings.ShowFullName, utils.Cfg.PrivacySettings.ShowEmailAddress), "Get Me", w, r) {
		return
	} else {
		user.Sanitize(map[string]bool{})
		w.Header().Set(model.HEADER_ETAG_SERVER, user.Etag(utils.Cfg.PrivacySettings.ShowFullName, utils.Cfg.PrivacySettings.ShowEmailAddress))
		w.Write([]byte(user.ToJson()))
		return
	}
}

func getInitialLoad(c *Context, w http.ResponseWriter, r *http.Request) {

	il := model.InitialLoad{}

	if len(c.Session.UserId) != 0 {
		var err *model.AppError

		il.User, err = app.GetUser(c.Session.UserId)
		if err != nil {
			c.Err = err
			return
		}
		il.User.Sanitize(map[string]bool{})


	}

	if app.SessionCacheLength() == 0 {
		// Below is a special case when intializating a new server
		// Lets check to make sure the server is really empty

		il.NoAccounts = app.IsFirstUserAccount()
	}

	il.ClientCfg = utils.ClientCfg
	if app.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		il.LicenseCfg = utils.ClientLicense
	} else {
		il.LicenseCfg = utils.GetSanitizedClientLicense()
	}

	w.Write([]byte(il.ToJson()))
}

func getUser(c *Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["user_id"]

	var user *model.User
	var err *model.AppError

	if user, err = app.GetUser(id); err != nil {
		c.Err = err
		return
	}

	etag := user.Etag(utils.Cfg.PrivacySettings.ShowFullName, utils.Cfg.PrivacySettings.ShowEmailAddress)

	if HandleEtag(etag, "Get User", w, r) {
		return
	} else {
		app.SanitizeProfile(user, c.IsSystemAdmin())
		w.Header().Set(model.HEADER_ETAG_SERVER, etag)
		w.Write([]byte(user.ToJson()))
		return
	}
}

func getByUsername(c *Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	var user *model.User
	var err *model.AppError

	if user, err = app.GetUserByUsername(username); err != nil {
		c.Err = err
		return
	} else if HandleEtag(user.Etag(utils.Cfg.PrivacySettings.ShowFullName, utils.Cfg.PrivacySettings.ShowEmailAddress), "Get By Username", w, r) {
		return
	} else {
		sanitizeProfile(c, user)

		w.Header().Set(model.HEADER_ETAG_SERVER, user.Etag(utils.Cfg.PrivacySettings.ShowFullName, utils.Cfg.PrivacySettings.ShowEmailAddress))
		w.Write([]byte(user.ToJson()))
		return
	}
}

func getByEmail(c *Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	email := params["email"]

	if user, err := app.GetUserByEmail(email); err != nil {
		c.Err = err
		return
	} else if HandleEtag(user.Etag(utils.Cfg.PrivacySettings.ShowFullName, utils.Cfg.PrivacySettings.ShowEmailAddress), "Get By Email", w, r) {
		return
	} else {
		sanitizeProfile(c, user)

		w.Header().Set(model.HEADER_ETAG_SERVER, user.Etag(utils.Cfg.PrivacySettings.ShowFullName, utils.Cfg.PrivacySettings.ShowEmailAddress))
		w.Write([]byte(user.ToJson()))
		return
	}
}

func getProfiles(c *Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	offset, err := strconv.Atoi(params["offset"])
	if err != nil {
		c.SetInvalidParam("getProfiles", "offset")
		return
	}

	limit, err := strconv.Atoi(params["limit"])
	if err != nil {
		c.SetInvalidParam("getProfiles", "limit")
		return
	}

	etag := app.GetUsersEtag() + params["offset"] + "." + params["limit"]
	if HandleEtag(etag, "Get Profiles", w, r) {
		return
	}

	if profiles, err := app.GetUsersMap(offset, limit, c.IsSystemAdmin()); err != nil {
		c.Err = err
		return
	} else {
		w.Header().Set(model.HEADER_ETAG_SERVER, etag)
		w.Write([]byte(model.UserMapToJson(profiles)))
	}
}

func getProfilesInTeam(c *Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	teamId := params["team_id"]



	offset, err := strconv.Atoi(params["offset"])
	if err != nil {
		c.SetInvalidParam("getProfilesInTeam", "offset")
		return
	}

	limit, err := strconv.Atoi(params["limit"])
	if err != nil {
		c.SetInvalidParam("getProfilesInTeam", "limit")
		return
	}

	etag := app.GetUsersInTeamEtag(teamId)
	if HandleEtag(etag, "Get Profiles In Team", w, r) {
		return
	}

	if profiles, err := app.GetUsersInTeamMap(teamId, offset, limit, c.IsSystemAdmin()); err != nil {
		c.Err = err
		return
	} else {
		w.Header().Set(model.HEADER_ETAG_SERVER, etag)
		w.Write([]byte(model.UserMapToJson(profiles)))
	}
}

func getProfilesInChannel(c *Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	channelId := params["channel_id"]


	if !app.SessionHasPermissionToChannel(c.Session, channelId, model.PERMISSION_READ_CHANNEL) {
		c.SetPermissionError(model.PERMISSION_READ_CHANNEL)
		return
	}


}

func getProfilesNotInChannel(c *Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	channelId := params["channel_id"]



	if !app.SessionHasPermissionToChannel(c.Session, channelId, model.PERMISSION_READ_CHANNEL) {
		c.SetPermissionError(model.PERMISSION_READ_CHANNEL)
		return
	}



}

func getAudits(c *Context, w http.ResponseWriter, r *http.Request) {


}

func getProfileImage(c *Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["user_id"]
	readFailed := false

	var etag string

	if users, err := app.GetUsersByIds([]string{id}, false); err != nil {
		c.Err = err
		return
	} else {
		if len(users) == 0 {
			c.Err = model.NewLocAppError("getProfileImage", "store.sql_user.get_profiles.app_error", nil, "")
			return
		}

		user := users[0]
		etag = strconv.FormatInt(user.LastPictureUpdate, 10)
		if HandleEtag(etag, "Profile Image", w, r) {
			return
		}

		var img []byte

		if err != nil {
			c.Err = err
			return
		}

		if readFailed {
			w.Header().Set("Cache-Control", "max-age=300, public") // 5 mins
		} else {
			w.Header().Set("Cache-Control", "max-age=86400, public") // 24 hrs
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set(model.HEADER_ETAG_SERVER, etag)
		w.Write(img)
	}
}



func updateUser(c *Context, w http.ResponseWriter, r *http.Request) {
	user := model.UserFromJson(r.Body)

	if user == nil {
		c.SetInvalidParam("updateUser", "user")
		return
	}

	if !app.SessionHasPermissionToUser(c.Session, user.Id) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	if ruser, err := app.UpdateUserAsUser(user, c.IsSystemAdmin()); err != nil {
		c.Err = err
		return
	} else {
		c.LogAudit("")
		w.Write([]byte(ruser.ToJson()))
	}
}

func updatePassword(c *Context, w http.ResponseWriter, r *http.Request) {
	c.LogAudit("attempted")

	props := model.MapFromJson(r.Body)
	userId := props["user_id"]
	if len(userId) != 26 {
		c.SetInvalidParam("updatePassword", "user_id")
		return
	}

	currentPassword := props["current_password"]
	if len(currentPassword) <= 0 {
		c.SetInvalidParam("updatePassword", "current_password")
		return
	}

	newPassword := props["new_password"]

	if userId != c.Session.UserId {
		c.Err = model.NewLocAppError("updatePassword", "api.user.update_password.context.app_error", nil, "")
		c.Err.StatusCode = http.StatusForbidden
		return
	}

	if err := app.UpdatePasswordAsUser(userId, currentPassword, newPassword); err != nil {
		c.LogAudit("failed")
		c.Err = err
		return
	} else {
		c.LogAudit("completed")

		data := make(map[string]string)
		data["user_id"] = c.Session.UserId
		w.Write([]byte(model.MapToJson(data)))
	}
}

func updateRoles(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)
	params := mux.Vars(r)

	userId := params["user_id"]
	if len(userId) != 26 {
		c.SetInvalidParam("updateMemberRoles", "user_id")
		return
	}

	newRoles := props["new_roles"]
	if !(model.IsValidUserRoles(newRoles)) {
		c.SetInvalidParam("updateMemberRoles", "new_roles")
		return
	}

	if !app.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_ROLES) {
		c.SetPermissionError(model.PERMISSION_MANAGE_ROLES)
		return
	}

	if _, err := app.UpdateUserRoles(userId, newRoles); err != nil {
		return
	} else {
		c.LogAuditWithUserId(userId, "roles="+newRoles)
	}

	rdata := map[string]string{}
	rdata["status"] = "ok"
	w.Write([]byte(model.MapToJson(rdata)))
}

func updateActive(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)

	userId := props["user_id"]
	if len(userId) != 26 {
		c.SetInvalidParam("updateActive", "user_id")
		return
	}

	active := props["active"] == "true"

	// true when you're trying to de-activate yourself
	isSelfDeactive := !active && userId == c.Session.UserId

	if !isSelfDeactive && !app.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		c.Err = model.NewLocAppError("updateActive", "api.user.update_active.permissions.app_error", nil, "userId="+userId)
		c.Err.StatusCode = http.StatusForbidden
		return
	}

	if ruser, err := app.UpdateActiveNoLdap(userId, active); err != nil {
		c.Err = err
	} else {
		c.LogAuditWithUserId(ruser.Id, fmt.Sprintf("active=%v", active))
		w.Write([]byte(ruser.ToJson()))
	}
}

func sendPasswordReset(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)

	email := props["email"]
	if len(email) == 0 {
		c.SetInvalidParam("sendPasswordReset", "email")
		return
	}

	if sent, err := app.SendPasswordReset(email, utils.GetSiteURL()); err != nil {
		c.Err = err
		return
	} else if sent {
		c.LogAudit("sent=" + email)
	}

	w.Write([]byte(model.MapToJson(props)))
}

func resetPassword(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)

	code := props["code"]


	newPassword := props["new_password"]

	c.LogAudit("attempt - code=" + code)

	if err := app.ResetPasswordFromCode(code, newPassword); err != nil {
		c.LogAudit("fail - code=" + code)
		c.Err = err
		return
	}

	c.LogAudit("success - code=" + code)

	rdata := map[string]string{}
	rdata["status"] = "ok"
	w.Write([]byte(model.MapToJson(rdata)))
}

func updateUserNotify(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)

	userId := props["user_id"]
	if len(userId) != 26 {
		c.SetInvalidParam("updateUserNotify", "user_id")
		return
	}

	if !app.SessionHasPermissionToUser(c.Session, userId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	delete(props, "user_id")

	email := props["email"]
	if len(email) == 0 {
		c.SetInvalidParam("updateUserNotify", "email")
		return
	}

	desktop_sound := props["desktop_sound"]
	if len(desktop_sound) == 0 {
		c.SetInvalidParam("updateUserNotify", "desktop_sound")
		return
	}

	desktop := props["desktop"]
	if len(desktop) == 0 {
		c.SetInvalidParam("updateUserNotify", "desktop")
		return
	}

	comments := props["comments"]
	if len(comments) == 0 {
		c.SetInvalidParam("updateUserNotify", "comments")
		return
	}

	ruser, err := app.UpdateUserNotifyProps(userId, props)
	if err != nil {
		c.Err = err
		return
	}

	c.LogAuditWithUserId(ruser.Id, "")

	options := utils.Cfg.GetSanitizeOptions()
	options["passwordupdate"] = false
	ruser.Sanitize(options)
	w.Write([]byte(ruser.ToJson()))
}


func verifyEmail(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)

	userId := props["uid"]
	if len(userId) != 26 {
		c.SetInvalidParam("verifyEmail", "uid")
		return
	}

	hashedId := props["hid"]
	if len(hashedId) == 0 {
		c.SetInvalidParam("verifyEmail", "hid")
		return
	}

	if hashedId == model.HashSha256(userId+utils.Cfg.EmailSettings.InviteSalt) {
		if c.Err = app.VerifyUserEmail(userId); c.Err != nil {
			return
		} else {
			c.LogAudit("Email Verified")
			return
		}
	}

	c.Err = model.NewLocAppError("verifyEmail", "api.user.verify_email.bad_link.app_error", nil, "")
	c.Err.StatusCode = http.StatusBadRequest
}




func sanitizeProfile(c *Context, user *model.User) *model.User {
	options := utils.Cfg.GetSanitizeOptions()

	if app.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		options["email"] = true
		options["fullname"] = true
		options["authservice"] = true
	}

	user.SanitizeProfile(options)

	return user
}


func getProfilesByIds(c *Context, w http.ResponseWriter, r *http.Request) {
	userIds := model.ArrayFromJson(r.Body)

	if len(userIds) == 0 {
		c.SetInvalidParam("getProfilesByIds", "user_ids")
		return
	}

	if profiles, err := app.GetUsersByIds(userIds, c.IsSystemAdmin()); err != nil {
		c.Err = err
		return
	} else {
		profileMap := map[string]*model.User{}
		for _, p := range profiles {
			profileMap[p.Id] = p
		}
		w.Write([]byte(model.UserMapToJson(profileMap)))
	}
}
