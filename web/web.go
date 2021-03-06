// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


package web

import (
	"net/http"
	"strings"


	l4g "github.com/alecthomas/log4go"
	"github.com/nomadsingles/platform/api"
	"github.com/nomadsingles/platform/app"
	"github.com/nomadsingles/platform/model"
	"github.com/nomadsingles/platform/utils"

	"github.com/mssola/user_agent"
)

func InitWeb() {
	l4g.Debug(utils.T("web.init.debug"))

	mainrouter := app.Srv.Router

	if *utils.Cfg.ServiceSettings.WebserverMode != "disabled" {
		staticDir := utils.FindDir(model.CLIENT_DIR)
		l4g.Debug("Using client directory at %v", staticDir)
		if *utils.Cfg.ServiceSettings.WebserverMode == "gzip" {
			mainrouter.PathPrefix("/static/").Handler(staticHandler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))))
		} else {
			mainrouter.PathPrefix("/static/").Handler(staticHandler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))))
		}

		mainrouter.Handle("/{anything:.*}", api.AppHandlerIndependent(root)).Methods("GET")
	}
}

func staticHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		handler.ServeHTTP(w, r)
	})
}

var browsersNotSupported string = "MSIE/8;MSIE/9;MSIE/10;Internet Explorer/8;Internet Explorer/9;Internet Explorer/10;Safari/7;Safari/8"

func CheckBrowserCompatability(c *api.Context, r *http.Request) bool {
	ua := user_agent.New(r.UserAgent())
	bname, bversion := ua.Browser()

	browsers := strings.Split(browsersNotSupported, ";")
	for _, browser := range browsers {
		version := strings.Split(browser, "/")

		if strings.HasPrefix(bname, version[0]) && strings.HasPrefix(bversion, version[1]) {
			return false
		}
	}

	return true

}

func root(c *api.Context, w http.ResponseWriter, r *http.Request) {
	if !CheckBrowserCompatability(c, r) {
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(c.T("web.check_browser_compatibility.app_error")))
		return
	}

	if api.IsApiCall(r) {
		api.Handle404(w, r)
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Connection", "keep-alive")


	http.ServeFile(w, r, utils.FindDir(model.CLIENT_DIR)+"index.html")
}
