// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"net/http"



	"github.com/nomadsingles/platform/model"
	"github.com/nomadsingles/platform/utils"
)

func CheckPasswordAndAllCriteria(user *model.User, password string, mfaToken string) *model.AppError {
	if err := CheckUserAdditionalAuthenticationCriteria(user, mfaToken); err != nil {
		return err
	}

	if err := checkUserPassword(user, password); err != nil {
		return err
	}

	return nil
}

// This to be used for places we check the users password when they are already logged in
func doubleCheckPassword(user *model.User, password string) *model.AppError {
	if err := checkUserLoginAttempts(user); err != nil {
		return err
	}

	if err := checkUserPassword(user, password); err != nil {
		return err
	}

	return nil
}

func checkUserPassword(user *model.User, password string) *model.AppError {
	if !model.ComparePassword(user.Password, password) {
		if result := <-Srv.Store.User().UpdateFailedPasswordAttempts(user.Id, user.FailedAttempts+1); result.Err != nil {
			return result.Err
		}

		return model.NewLocAppError("checkUserPassword", "api.user.check_user_password.invalid.app_error", nil, "user_id="+user.Id)
	} else {
		if result := <-Srv.Store.User().UpdateFailedPasswordAttempts(user.Id, 0); result.Err != nil {
			return result.Err
		}

		return nil
	}
}



func CheckUserAdditionalAuthenticationCriteria(user *model.User, mfaToken string) *model.AppError {

	if err := checkEmailVerified(user); err != nil {
		return err
	}

	if err := checkUserNotDisabled(user); err != nil {
		return err
	}

	if err := checkUserLoginAttempts(user); err != nil {
		return err
	}

	return nil
}



func checkUserLoginAttempts(user *model.User) *model.AppError {
	if user.FailedAttempts >= utils.Cfg.ServiceSettings.MaximumLoginAttempts {
		return model.NewAppError("checkUserLoginAttempts", "api.user.check_user_login_attempts.too_many.app_error", nil, "user_id="+user.Id, http.StatusForbidden)
	}

	return nil
}

func checkEmailVerified(user *model.User) *model.AppError {
	if !user.EmailVerified && utils.Cfg.EmailSettings.RequireEmailVerification {
		return model.NewLocAppError("Login", "api.user.login.not_verified.app_error", nil, "user_id="+user.Id)
	}
	return nil
}

func checkUserNotDisabled(user *model.User) *model.AppError {
	if user.DeleteAt > 0 {
		return model.NewLocAppError("Login", "api.user.login.inactive.app_error", nil, "user_id="+user.Id)
	}
	return nil
}

func authenticateUser(user *model.User, password, mfaToken string) (*model.User, *model.AppError) {


 if user.AuthService != "" {
		authService := user.AuthService

		err := model.NewLocAppError("login", "api.user.login.use_auth_service.app_error", map[string]interface{}{"AuthService": authService}, "")
		err.StatusCode = http.StatusBadRequest
		return user, err
	} else {
		if err := CheckPasswordAndAllCriteria(user, password, mfaToken); err != nil {
			err.StatusCode = http.StatusUnauthorized
			return user, err
		} else {
			return user, nil
		}
	}
}
