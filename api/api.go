// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/platform/app"

	"github.com/mattermost/platform/model"


	_ "github.com/nicksnyder/go-i18n/i18n"
)

type Routes struct {
	Root    *mux.Router // ''
	ApiRoot *mux.Router // 'api/v3'

	Users    *mux.Router // 'api/v3/users'
	NeedUser *mux.Router // 'api/v3/users/{user_id:[A-Za-z0-9]+}'

	Teams    *mux.Router // 'api/v3/teams'
	NeedTeam *mux.Router // 'api/v3/teams/{team_id:[A-Za-z0-9]+}'

	Channels        *mux.Router // 'api/v3/teams/{team_id:[A-Za-z0-9]+}/channels'
	NeedChannel     *mux.Router // 'api/v3/teams/{team_id:[A-Za-z0-9]+}/channels/{channel_id:[A-Za-z0-9]+}'
	NeedChannelName *mux.Router // 'api/v3/teams/{team_id:[A-Za-z0-9]+}/channels/name/{channel_name:[A-Za-z0-9_-]+}'

	Posts    *mux.Router // 'api/v3/teams/{team_id:[A-Za-z0-9]+}/channels/{channel_id:[A-Za-z0-9]+}/posts'
	NeedPost *mux.Router // 'api/v3/teams/{team_id:[A-Za-z0-9]+}/channels/{channel_id:[A-Za-z0-9]+}/posts/{post_id:[A-Za-z0-9]+}'

	Commands *mux.Router // 'api/v3/teams/{team_id:[A-Za-z0-9]+}/commands'
	Hooks    *mux.Router // 'api/v3/teams/{team_id:[A-Za-z0-9]+}/hooks'

	TeamFiles *mux.Router // 'api/v3/teams/{team_id:[A-Za-z0-9]+}/files'
	Files     *mux.Router // 'api/v3/files'
	NeedFile  *mux.Router // 'api/v3/files/{file_id:[A-Za-z0-9]+}'

	OAuth *mux.Router // 'api/v3/oauth'

	Admin *mux.Router // 'api/v3/admin'

	General *mux.Router // 'api/v3/general'

	Preferences *mux.Router // 'api/v3/preferences'

	License *mux.Router // 'api/v3/license'

	Public *mux.Router // 'api/v3/public'

	Emoji *mux.Router // 'api/v3/emoji'

	Webrtc *mux.Router // 'api/v3/webrtc'
}

var BaseRoutes *Routes

func InitRouter() {
	app.Srv.Router = mux.NewRouter()
	app.Srv.Router.NotFoundHandler = http.HandlerFunc(Handle404)
}

func InitApi() {
	BaseRoutes = &Routes{}
	BaseRoutes.Root = app.Srv.Router
	BaseRoutes.ApiRoot = app.Srv.Router.PathPrefix(model.API_URL_SUFFIX_V3).Subrouter()


	// 404 on any api route before web.go has a chance to serve it
	app.Srv.Router.Handle("/api/{anything:.*}", http.HandlerFunc(Handle404))


}

func HandleEtag(etag string, routeName string, w http.ResponseWriter, r *http.Request) bool {


	return false
}

func ReturnStatusOK(w http.ResponseWriter) {
	m := make(map[string]string)
	m[model.STATUS] = model.STATUS_OK
	w.Write([]byte(model.MapToJson(m)))
}
