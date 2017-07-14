// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


package store

import (
	"time"

	l4g "github.com/alecthomas/log4go"
	"github.com/nomadsingles/platform/model"
)

type StoreResult struct {
	Data interface{}
	Err  *model.AppError
}

type StoreChannel chan StoreResult

func Must(sc StoreChannel) interface{} {
	r := <-sc
	if r.Err != nil {
		l4g.Close()
		time.Sleep(time.Second)
		panic(r.Err)
	}

	return r.Data
}

type Store interface {
	User() UserStore
	Session() SessionStore
	OAuth() OAuthStore
	System() SystemStore
	MarkSystemRanUnitTests()
	Close()
	DropAllTables()
	TotalMasterDbConnections() int
	TotalReadDbConnections() int
	TotalSearchDbConnections() int
}


type UserStore interface {
	Save(user *model.User) StoreChannel
	Update(user *model.User, allowRoleUpdate bool) StoreChannel
	UpdateLastPictureUpdate(userId string) StoreChannel
	UpdateUpdateAt(userId string) StoreChannel
	UpdatePassword(userId, newPassword string) StoreChannel
	UpdateAuthData(userId string, service string, authData *string, email string, resetMfa bool) StoreChannel
	UpdateMfaSecret(userId, secret string) StoreChannel
	UpdateMfaActive(userId string, active bool) StoreChannel
	Get(id string) StoreChannel
	GetAll() StoreChannel
	InvalidateProfilesInChannelCacheByUser(userId string)
	InvalidateProfilesInChannelCache(channelId string)

	GetAllProfiles(offset int, limit int) StoreChannel
	GetProfiles(teamId string, offset int, limit int) StoreChannel
	GetProfileByIds(userId []string, allowFromCache bool) StoreChannel
	InvalidatProfileCacheForUser(userId string)
	GetByEmail(email string) StoreChannel
	GetByAuth(authData *string, authService string) StoreChannel
	GetAllUsingAuthService(authService string) StoreChannel
	GetByUsername(username string) StoreChannel
	GetForLogin(loginId string, allowSignInWithUsername, allowSignInWithEmail, ldapEnabled bool) StoreChannel
	VerifyEmail(userId string) StoreChannel
	GetEtagForAllProfiles() StoreChannel
	GetEtagForProfiles(teamId string) StoreChannel
	UpdateFailedPasswordAttempts(userId string, attempts int) StoreChannel
	GetTotalUsersCount() StoreChannel
	GetSystemAdminProfiles() StoreChannel
	PermanentDelete(userId string) StoreChannel
	AnalyticsUniqueUserCount(teamId string) StoreChannel
	AnalyticsActiveCount(time int64) StoreChannel
	GetUnreadCount(userId string) StoreChannel
	GetUnreadCountForChannel(userId string, channelId string) StoreChannel
	GetRecentlyActiveUsersForTeam(teamId string) StoreChannel
	Search(teamId string, term string, options map[string]bool) StoreChannel
	SearchNotInTeam(notInTeamId string, term string, options map[string]bool) StoreChannel
	SearchInChannel(channelId string, term string, options map[string]bool) StoreChannel
	SearchNotInChannel(teamId string, channelId string, term string, options map[string]bool) StoreChannel
	SearchWithoutTeam(term string, options map[string]bool) StoreChannel
	AnalyticsGetInactiveUsersCount() StoreChannel
	AnalyticsGetSystemAdminCount() StoreChannel
	GetProfilesNotInTeam(teamId string, offset int, limit int) StoreChannel
	GetEtagForProfilesNotInTeam(teamId string) StoreChannel
}

type SessionStore interface {
	Save(session *model.Session) StoreChannel
	Get(sessionIdOrToken string) StoreChannel
	GetSessions(userId string) StoreChannel
	GetSessionsWithActiveDeviceIds(userId string) StoreChannel
	Remove(sessionIdOrToken string) StoreChannel
	RemoveAllSessions() StoreChannel
	PermanentDeleteSessionsByUser(teamId string) StoreChannel
	UpdateLastActivityAt(sessionId string, time int64) StoreChannel
	UpdateRoles(userId string, roles string) StoreChannel
	UpdateDeviceId(id string, deviceId string, expiresAt int64) StoreChannel
	AnalyticsSessionCount() StoreChannel
}



type OAuthStore interface {
	SaveApp(app *model.OAuthApp) StoreChannel
	UpdateApp(app *model.OAuthApp) StoreChannel
	GetApp(id string) StoreChannel
	GetAppByUser(userId string) StoreChannel
	GetApps() StoreChannel
	GetAuthorizedApps(userId string) StoreChannel
	DeleteApp(id string) StoreChannel
	SaveAuthData(authData *model.AuthData) StoreChannel
	GetAuthData(code string) StoreChannel
	RemoveAuthData(code string) StoreChannel
	PermanentDeleteAuthDataByUser(userId string) StoreChannel
	SaveAccessData(accessData *model.AccessData) StoreChannel
	UpdateAccessData(accessData *model.AccessData) StoreChannel
	GetAccessData(token string) StoreChannel
	GetAccessDataByUserForApp(userId, clientId string) StoreChannel
	GetAccessDataByRefreshToken(token string) StoreChannel
	GetPreviousAccessData(userId, clientId string) StoreChannel
	RemoveAccessData(token string) StoreChannel
}

type SystemStore interface {
	Save(system *model.System) StoreChannel
	SaveOrUpdate(system *model.System) StoreChannel
	Update(system *model.System) StoreChannel
	Get() StoreChannel
	GetByName(name string) StoreChannel
}


