// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"bytes"

	"fmt"
	"hash/fnv"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"io/ioutil"

	"net/http"
	"strconv"
	"strings"

	l4g "github.com/alecthomas/log4go"

	"github.com/golang/freetype"

	"github.com/mattermost/platform/model"

	"github.com/mattermost/platform/utils"
)

func CreateUserWithHash(user *model.User, hash string, data string) (*model.User, *model.AppError) {
	if err := IsUserSignUpAllowed(); err != nil {
		return nil, err
	}

	props := model.MapFromJson(strings.NewReader(data))

	if hash != model.HashSha256(fmt.Sprintf("%v:%v", data, utils.Cfg.EmailSettings.InviteSalt)) {
		return nil, model.NewLocAppError("CreateUserWithHash", "api.user.create_user.signup_link_invalid.app_error", nil, "")
	}

	if t, err := strconv.ParseInt(props["time"], 10, 64); err != nil || model.GetMillis()-t > 1000*60*60*48 { // 48 hours
		return nil, model.NewLocAppError("CreateUserWithHash", "api.user.create_user.signup_link_expired.app_error", nil, "")
	}



	user.Email = props["email"]
	user.EmailVerified = true

	var ruser *model.User
	var err *model.AppError
	if ruser, err = CreateUser(user); err != nil {
		return nil, err
	}



	return ruser, nil
}



func CreateUserAsAdmin(user *model.User) (*model.User, *model.AppError) {
	ruser, err := CreateUser(user)
	if err != nil {
		return nil, err
	}



	return ruser, nil
}

func CreateUserFromSignup(user *model.User) (*model.User, *model.AppError) {
	if err := IsUserSignUpAllowed(); err != nil {
		return nil, err
	}

	if !IsFirstUserAccount() && !*utils.Cfg.TeamSettings.EnableOpenServer {
		err := model.NewLocAppError("CreateUserFromSignup", "api.user.create_user.no_open_server", nil, "email="+user.Email)
		err.StatusCode = http.StatusForbidden
		return nil, err
	}

	user.EmailVerified = false

	ruser, err := CreateUser(user)
	if err != nil {
		return nil, err
	}



	return ruser, nil
}

func IsUserSignUpAllowed() *model.AppError {
	if !utils.Cfg.EmailSettings.EnableSignUpWithEmail || !utils.Cfg.TeamSettings.EnableUserCreation {
		err := model.NewLocAppError("IsUserSignUpAllowed", "api.user.create_user.signup_email_disabled.app_error", nil, "")
		err.StatusCode = http.StatusNotImplemented
		return err
	}
	return nil
}

func IsFirstUserAccount() bool {
	if SessionCacheLength() == 0 {
		if cr := <-Srv.Store.User().GetTotalUsersCount(); cr.Err != nil {
			l4g.Error(cr.Err)
			return false
		} else {
			count := cr.Data.(int64)
			if count <= 0 {
				return true
			}
		}
	}

	return false
}

func CreateUser(user *model.User) (*model.User, *model.AppError) {
	if !user.IsSSOUser() && !CheckUserDomain(user, utils.Cfg.TeamSettings.RestrictCreationToDomains) {
		return nil, model.NewLocAppError("CreateUser", "api.user.create_user.accepted_domain.app_error", nil, "")
	}

	user.Roles = model.ROLE_SYSTEM_USER.Id

	// Below is a special case where the first user in the entire
	// system is granted the system_admin role
	if result := <-Srv.Store.User().GetTotalUsersCount(); result.Err != nil {
		return nil, result.Err
	} else {
		count := result.Data.(int64)
		if count <= 0 {
			user.Roles = model.ROLE_SYSTEM_ADMIN.Id + " " + model.ROLE_SYSTEM_USER.Id
		}
	}

	user.Locale = *utils.Cfg.LocalizationSettings.DefaultClientLocale

	if ruser, err := createUser(user); err != nil {
		return nil, err
	} else {


		return ruser, nil
	}
}

func createUser(user *model.User) (*model.User, *model.AppError) {
	user.MakeNonNil()

	if err := utils.IsPasswordValid(user.Password); user.AuthService == "" && err != nil {
		return nil, err
	}

	if result := <-Srv.Store.User().Save(user); result.Err != nil {
		l4g.Error(utils.T("api.user.create_user.save.error"), result.Err)
		return nil, result.Err
	} else {
		ruser := result.Data.(*model.User)

		if user.EmailVerified {
			if err := VerifyUserEmail(ruser.Id); err != nil {
				l4g.Error(utils.T("api.user.create_user.verified.error"), err)
			}
		}



		ruser.Sanitize(map[string]bool{})

		return ruser, nil
	}
}

func CreateOAuthUser(service string, userData io.Reader, teamId string) (*model.User, *model.AppError) {
	if !utils.Cfg.TeamSettings.EnableUserCreation {
		return nil, model.NewAppError("CreateOAuthUser", "api.user.create_user.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	var user *model.User


	if user == nil {
		return nil, model.NewLocAppError("CreateOAuthUser", "api.user.create_oauth_user.create.app_error", map[string]interface{}{"Service": service}, "")
	}

	suchan := Srv.Store.User().GetByAuth(user.AuthData, service)
	euchan := Srv.Store.User().GetByEmail(user.Email)

	found := true
	count := 0
	for found {
		if found = IsUsernameTaken(user.Username); found {
			user.Username = user.Username + strconv.Itoa(count)
			count += 1
		}
	}

	if result := <-suchan; result.Err == nil {
		return nil, model.NewLocAppError("CreateOAuthUser", "api.user.create_oauth_user.already_used.app_error", map[string]interface{}{"Service": service}, "email="+user.Email)
	}

	if result := <-euchan; result.Err == nil {
		authService := result.Data.(*model.User).AuthService
		if authService == "" {
			return nil, model.NewLocAppError("CreateOAuthUser", "api.user.create_oauth_user.already_attached.app_error",
				map[string]interface{}{"Service": service, "Auth": model.USER_AUTH_SERVICE_EMAIL}, "email="+user.Email)
		} else {
			return nil, model.NewLocAppError("CreateOAuthUser", "api.user.create_oauth_user.already_attached.app_error",
				map[string]interface{}{"Service": service, "Auth": authService}, "email="+user.Email)
		}
	}

	user.EmailVerified = true

	ruser, err := CreateUser(user)
	if err != nil {
		return nil, err
	}


	return ruser, nil
}

// Check that a user's email domain matches a list of space-delimited domains as a string.
func CheckUserDomain(user *model.User, domains string) bool {
	if len(domains) == 0 {
		return true
	}

	domainArray := strings.Fields(strings.TrimSpace(strings.ToLower(strings.Replace(strings.Replace(domains, "@", " ", -1), ",", " ", -1))))

	matched := false
	for _, d := range domainArray {
		if strings.HasSuffix(strings.ToLower(user.Email), "@"+d) {
			matched = true
			break
		}
	}

	return matched
}

// Check if the username is already used by another user. Return false if the username is invalid.
func IsUsernameTaken(name string) bool {

	if !model.IsValidUsername(name) {
		return false
	}

	if result := <-Srv.Store.User().GetByUsername(name); result.Err != nil {
		return false
	}

	return true
}

func GetUser(userId string) (*model.User, *model.AppError) {
	if result := <-Srv.Store.User().Get(userId); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.(*model.User), nil
	}
}

func GetUserByUsername(username string) (*model.User, *model.AppError) {
	if result := <-Srv.Store.User().GetByUsername(username); result.Err != nil && result.Err.Id == "store.sql_user.get_by_username.app_error" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else {
		return result.Data.(*model.User), nil
	}
}

func GetUserByEmail(email string) (*model.User, *model.AppError) {

	if result := <-Srv.Store.User().GetByEmail(email); result.Err != nil && result.Err.Id == "store.sql_user.missing_account.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.User), nil
	}
}

func GetUserByAuth(authData *string, authService string) (*model.User, *model.AppError) {
	if result := <-Srv.Store.User().GetByAuth(authData, authService); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.(*model.User), nil
	}
}



func GetUsers(offset int, limit int) ([]*model.User, *model.AppError) {
	if result := <-Srv.Store.User().GetAllProfiles(offset, limit); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.([]*model.User), nil
	}
}

func GetUsersMap(offset int, limit int, asAdmin bool) (map[string]*model.User, *model.AppError) {
	users, err := GetUsers(offset, limit)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*model.User, len(users))

	for _, user := range users {
		SanitizeProfile(user, asAdmin)
		userMap[user.Id] = user
	}

	return userMap, nil
}

func GetUsersPage(page int, perPage int, asAdmin bool) ([]*model.User, *model.AppError) {
	users, err := GetUsers(page*perPage, perPage)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		SanitizeProfile(user, asAdmin)
	}

	return users, nil
}

func GetUsersEtag() string {
	return (<-Srv.Store.User().GetEtagForAllProfiles()).Data.(string)
}

func GetUsersInTeam(teamId string, offset int, limit int) ([]*model.User, *model.AppError) {
	if result := <-Srv.Store.User().GetProfiles(teamId, offset, limit); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.([]*model.User), nil
	}
}

func GetUsersNotInTeam(teamId string, offset int, limit int) ([]*model.User, *model.AppError) {
	if result := <-Srv.Store.User().GetProfilesNotInTeam(teamId, offset, limit); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.([]*model.User), nil
	}
}

func GetUsersInTeamMap(teamId string, offset int, limit int, asAdmin bool) (map[string]*model.User, *model.AppError) {
	users, err := GetUsersInTeam(teamId, offset, limit)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*model.User, len(users))

	for _, user := range users {
		SanitizeProfile(user, asAdmin)
		userMap[user.Id] = user
	}

	return userMap, nil
}

func GetUsersInTeamPage(teamId string, page int, perPage int, asAdmin bool) ([]*model.User, *model.AppError) {
	users, err := GetUsersInTeam(teamId, page*perPage, perPage)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		SanitizeProfile(user, asAdmin)
	}

	return users, nil
}

func GetUsersNotInTeamPage(teamId string, page int, perPage int, asAdmin bool) ([]*model.User, *model.AppError) {
	users, err := GetUsersNotInTeam(teamId, page*perPage, perPage)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		SanitizeProfile(user, asAdmin)
	}

	return users, nil
}

func GetUsersInTeamEtag(teamId string) string {
	return (<-Srv.Store.User().GetEtagForProfiles(teamId)).Data.(string)
}

func GetUsersNotInTeamEtag(teamId string) string {
	return (<-Srv.Store.User().GetEtagForProfilesNotInTeam(teamId)).Data.(string)
}




func GetUsersByIds(userIds []string, asAdmin bool) ([]*model.User, *model.AppError) {
	if result := <-Srv.Store.User().GetProfileByIds(userIds, true); result.Err != nil {
		return nil, result.Err
	} else {
		users := result.Data.([]*model.User)

		for _, u := range users {
			SanitizeProfile(u, asAdmin)
		}

		return users, nil
	}
}


func CreateProfileImage(username string, userId string) ([]byte, *model.AppError) {
	colors := []color.NRGBA{
		{197, 8, 126, 255},
		{227, 207, 18, 255},
		{28, 181, 105, 255},
		{35, 188, 224, 255},
		{116, 49, 196, 255},
		{197, 8, 126, 255},
		{197, 19, 19, 255},
		{250, 134, 6, 255},
		{227, 207, 18, 255},
		{123, 201, 71, 255},
		{28, 181, 105, 255},
		{35, 188, 224, 255},
		{116, 49, 196, 255},
		{197, 8, 126, 255},
		{197, 19, 19, 255},
		{250, 134, 6, 255},
		{227, 207, 18, 255},
		{123, 201, 71, 255},
		{28, 181, 105, 255},
		{35, 188, 224, 255},
		{116, 49, 196, 255},
		{197, 8, 126, 255},
		{197, 19, 19, 255},
		{250, 134, 6, 255},
		{227, 207, 18, 255},
		{123, 201, 71, 255},
	}

	h := fnv.New32a()
	h.Write([]byte(userId))
	seed := h.Sum32()

	initial := string(strings.ToUpper(username)[0])

	fontBytes, err := ioutil.ReadFile(utils.FindDir("fonts") + utils.Cfg.FileSettings.InitialFont)
	if err != nil {
		return nil, model.NewLocAppError("CreateProfileImage", "api.user.create_profile_image.default_font.app_error", nil, err.Error())
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, model.NewLocAppError("CreateProfileImage", "api.user.create_profile_image.default_font.app_error", nil, err.Error())
	}

	width := int(utils.Cfg.FileSettings.ProfileWidth)
	height := int(utils.Cfg.FileSettings.ProfileHeight)
	color := colors[int64(seed)%int64(len(colors))]
	dstImg := image.NewRGBA(image.Rect(0, 0, width, height))
	srcImg := image.White
	draw.Draw(dstImg, dstImg.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
	size := float64((width + height) / 4)

	c := freetype.NewContext()
	c.SetFont(font)
	c.SetFontSize(size)
	c.SetClip(dstImg.Bounds())
	c.SetDst(dstImg)
	c.SetSrc(srcImg)

	pt := freetype.Pt(width/6, height*2/3)
	_, err = c.DrawString(initial, pt)
	if err != nil {
		return nil, model.NewLocAppError("CreateProfileImage", "api.user.create_profile_image.initial.app_error", nil, err.Error())
	}

	buf := new(bytes.Buffer)

	if imgErr := png.Encode(buf, dstImg); imgErr != nil {
		return nil, model.NewLocAppError("CreateProfileImage", "api.user.create_profile_image.encode.app_error", nil, imgErr.Error())
	} else {
		return buf.Bytes(), nil
	}
}


func UpdatePasswordAsUser(userId, currentPassword, newPassword string) *model.AppError {
	var user *model.User
	var err *model.AppError

	if user, err = GetUser(userId); err != nil {
		return err
	}

	if user == nil {
		err = model.NewAppError("updatePassword", "api.user.update_password.valid_account.app_error", nil, "", http.StatusBadRequest)
		return err
	}

	if user.AuthData != nil && *user.AuthData != "" {
		err = model.NewAppError("updatePassword", "api.user.update_password.oauth.app_error", nil, "auth_service="+user.AuthService, http.StatusBadRequest)
		return err
	}

	if err := doubleCheckPassword(user, currentPassword); err != nil {
		if err.Id == "api.user.check_user_password.invalid.app_error" {
			err = model.NewAppError("updatePassword", "api.user.update_password.incorrect.app_error", nil, "", http.StatusBadRequest)
		}
		return err
	}

	T := utils.GetUserTranslations(user.Locale)

	if err := UpdatePasswordSendEmail(user, newPassword, T("api.user.update_password.menu")); err != nil {
		return err
	}

	return nil
}

func UpdateActiveNoLdap(userId string, active bool) (*model.User, *model.AppError) {
	var user *model.User
	var err *model.AppError
	if user, err = GetUser(userId); err != nil {
		return nil, err
	}



	return UpdateActive(user, active)
}

func UpdateActive(user *model.User, active bool) (*model.User, *model.AppError) {
	if active {
		user.DeleteAt = 0
	} else {
		user.DeleteAt = model.GetMillis()
	}

	if result := <-Srv.Store.User().Update(user, true); result.Err != nil {
		return nil, result.Err
	} else {
		if user.DeleteAt > 0 {
			if err := RevokeAllSessions(user.Id); err != nil {
				return nil, err
			}
		}

		if extra := <-Srv.Store.Channel().ExtraUpdateByUser(user.Id, model.GetMillis()); extra.Err != nil {
			return nil, extra.Err
		}

		ruser := result.Data.([2]*model.User)[0]
		options := utils.Cfg.GetSanitizeOptions()
		options["passwordupdate"] = false
		ruser.Sanitize(options)



		return ruser, nil
	}
}

func SanitizeProfile(user *model.User, asAdmin bool) {
	options := utils.Cfg.GetSanitizeOptions()
	if asAdmin {
		options["email"] = true
		options["fullname"] = true
		options["authservice"] = true
	}
	user.SanitizeProfile(options)
}

func UpdateUserAsUser(user *model.User, asAdmin bool) (*model.User, *model.AppError) {
	updatedUser, err := UpdateUser(user, true)
	if err != nil {
		return nil, err
	}

	sendUpdatedUserEvent(*updatedUser, asAdmin)

	return updatedUser, nil
}

func PatchUser(userId string, patch *model.UserPatch, asAdmin bool) (*model.User, *model.AppError) {
	user, err := GetUser(userId)
	if err != nil {
		return nil, err
	}

	user.Patch(patch)

	updatedUser, err := UpdateUser(user, true)
	if err != nil {
		return nil, err
	}

	sendUpdatedUserEvent(*updatedUser, asAdmin)

	return updatedUser, nil
}

func sendUpdatedUserEvent(user model.User, asAdmin bool) {

}

func UpdateUser(user *model.User, sendNotifications bool) (*model.User, *model.AppError) {
	if result := <-Srv.Store.User().Update(user, false); result.Err != nil {
		return nil, result.Err
	} else {
		rusers := result.Data.([2]*model.User)

		if sendNotifications {
			if rusers[0].Email != rusers[1].Email {

			}

			if rusers[0].Username != rusers[1].Username {

			}
		}


		return rusers[0], nil
	}
}

func UpdateUserNotifyProps(userId string, props map[string]string) (*model.User, *model.AppError) {
	var user *model.User
	var err *model.AppError
	if user, err = GetUser(userId); err != nil {
		return nil, err
	}

	user.NotifyProps = props

	var ruser *model.User
	if ruser, err = UpdateUser(user, true); err != nil {
		return nil, err
	}

	return ruser, nil
}



func UpdatePasswordByUserIdSendEmail(userId, newPassword, method string) *model.AppError {
	var user *model.User
	var err *model.AppError
	if user, err = GetUser(userId); err != nil {
		return err
	}

	return UpdatePasswordSendEmail(user, newPassword, method)
}

func UpdatePassword(user *model.User, newPassword string) *model.AppError {
	if err := utils.IsPasswordValid(newPassword); err != nil {
		return err
	}

	hashedPassword := model.HashPassword(newPassword)

	if result := <-Srv.Store.User().UpdatePassword(user.Id, hashedPassword); result.Err != nil {
		return model.NewLocAppError("UpdatePassword", "api.user.update_password.failed.app_error", nil, result.Err.Error())
	}

	return nil
}

func UpdatePasswordSendEmail(user *model.User, newPassword, method string) *model.AppError {
	if err := UpdatePassword(user, newPassword); err != nil {
		return err
	}


	return nil
}

func ResetPasswordFromCode(code, newPassword string) *model.AppError {


	return nil
}

func SendPasswordReset(email string, siteURL string) (bool, *model.AppError) {

	return true, nil
}


func DeletePasswordRecoveryForUser(userId string) *model.AppError {
	if result := <-Srv.Store.PasswordRecovery().Delete(userId); result.Err != nil {
		return result.Err
	}

	return nil
}

func UpdateUserRoles(userId string, newRoles string) (*model.User, *model.AppError) {
	var user *model.User
	var err *model.AppError
	if user, err = GetUser(userId); err != nil {
		err.StatusCode = http.StatusBadRequest
		return nil, err
	}

	user.Roles = newRoles
	uchan := Srv.Store.User().Update(user, true)
	schan := Srv.Store.Session().UpdateRoles(user.Id, newRoles)

	var ruser *model.User
	if result := <-uchan; result.Err != nil {
		return nil, result.Err
	} else {
		ruser = result.Data.([2]*model.User)[0]
	}

	if result := <-schan; result.Err != nil {
		// soft error since the user roles were still updated
		l4g.Error(result.Err)
	}

	ClearSessionCacheForUser(user.Id)

	return ruser, nil
}

func PermanentDeleteUser(user *model.User) *model.AppError {
	l4g.Warn(utils.T("api.user.permanent_delete_user.attempting.warn"), user.Email, user.Id)
	if user.IsInRole(model.ROLE_SYSTEM_ADMIN.Id) {
		l4g.Warn(utils.T("api.user.permanent_delete_user.system_admin.warn"), user.Email)
	}

	if _, err := UpdateActive(user, false); err != nil {
		return err
	}

	if result := <-Srv.Store.Session().PermanentDeleteSessionsByUser(user.Id); result.Err != nil {
		return result.Err
	}

	if result := <-Srv.Store.OAuth().PermanentDeleteAuthDataByUser(user.Id); result.Err != nil {
		return result.Err
	}

	if result := <-Srv.Store.Webhook().PermanentDeleteIncomingByUser(user.Id); result.Err != nil {
		return result.Err
	}

	if result := <-Srv.Store.Webhook().PermanentDeleteOutgoingByUser(user.Id); result.Err != nil {
		return result.Err
	}

	if result := <-Srv.Store.Command().PermanentDeleteByUser(user.Id); result.Err != nil {
		return result.Err
	}

	if result := <-Srv.Store.Preference().PermanentDeleteByUser(user.Id); result.Err != nil {
		return result.Err
	}

	if result := <-Srv.Store.Channel().PermanentDeleteMembersByUser(user.Id); result.Err != nil {
		return result.Err
	}

	if result := <-Srv.Store.Post().PermanentDeleteByUser(user.Id); result.Err != nil {
		return result.Err
	}

	if result := <-Srv.Store.User().PermanentDelete(user.Id); result.Err != nil {
		return result.Err
	}

	if result := <-Srv.Store.Audit().PermanentDeleteByUser(user.Id); result.Err != nil {
		return result.Err
	}

	if result := <-Srv.Store.Team().RemoveAllMembersByUser(user.Id); result.Err != nil {
		return result.Err
	}

	if result := <-Srv.Store.PasswordRecovery().Delete(user.Id); result.Err != nil {
		return result.Err
	}

	l4g.Warn(utils.T("api.user.permanent_delete_user.deleted.warn"), user.Email, user.Id)

	return nil
}

func PermanentDeleteAllUsers() *model.AppError {
	if result := <-Srv.Store.User().GetAll(); result.Err != nil {
		return result.Err
	} else {
		users := result.Data.([]*model.User)
		for _, user := range users {
			PermanentDeleteUser(user)
		}
	}

	return nil
}

func VerifyUserEmail(userId string) *model.AppError {
	if err := (<-Srv.Store.User().VerifyEmail(userId)).Err; err != nil {
		return err
	}

	return nil
}



func SearchUsersInChannel(channelId string, term string, searchOptions map[string]bool, asAdmin bool) ([]*model.User, *model.AppError) {
	if result := <-Srv.Store.User().SearchInChannel(channelId, term, searchOptions); result.Err != nil {
		return nil, result.Err
	} else {
		users := result.Data.([]*model.User)

		for _, user := range users {
			SanitizeProfile(user, asAdmin)
		}

		return users, nil
	}
}

func SearchUsersNotInChannel(teamId string, channelId string, term string, searchOptions map[string]bool, asAdmin bool) ([]*model.User, *model.AppError) {
	if result := <-Srv.Store.User().SearchNotInChannel(teamId, channelId, term, searchOptions); result.Err != nil {
		return nil, result.Err
	} else {
		users := result.Data.([]*model.User)

		for _, user := range users {
			SanitizeProfile(user, asAdmin)
		}

		return users, nil
	}
}

func SearchUsersInTeam(teamId string, term string, searchOptions map[string]bool, asAdmin bool) ([]*model.User, *model.AppError) {
	if result := <-Srv.Store.User().Search(teamId, term, searchOptions); result.Err != nil {
		return nil, result.Err
	} else {
		users := result.Data.([]*model.User)

		for _, user := range users {
			SanitizeProfile(user, asAdmin)
		}

		return users, nil
	}
}

func SearchUsersNotInTeam(notInTeamId string, term string, searchOptions map[string]bool, asAdmin bool) ([]*model.User, *model.AppError) {
	if result := <-Srv.Store.User().SearchNotInTeam(notInTeamId, term, searchOptions); result.Err != nil {
		return nil, result.Err
	} else {
		users := result.Data.([]*model.User)

		for _, user := range users {
			SanitizeProfile(user, asAdmin)
		}

		return users, nil
	}
}

func SearchUsersWithoutTeam(term string, searchOptions map[string]bool, asAdmin bool) ([]*model.User, *model.AppError) {
	if result := <-Srv.Store.User().SearchWithoutTeam(term, searchOptions); result.Err != nil {
		return nil, result.Err
	} else {
		users := result.Data.([]*model.User)

		for _, user := range users {
			SanitizeProfile(user, asAdmin)
		}

		return users, nil
	}
}


