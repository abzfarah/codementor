// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.


package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nomadsingles/platform/app"

	"github.com/nomadsingles/platform/model"


	_ "github.com/nicksnyder/go-i18n/i18n"
)

type Routes struct {
	Root    *mux.Router // ''
	ApiRoot *mux.Router // 'api/v3'

	Users    *mux.Router // 'api/v3/users'

	OAuth *mux.Router // 'api/v3/oauth'


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
