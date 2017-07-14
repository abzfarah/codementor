// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.


package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"


	l4g "github.com/alecthomas/log4go"

	goi18n "github.com/nicksnyder/go-i18n/i18n"

	"github.com/nomadsingles/platform/app"

	"github.com/nomadsingles/platform/model"
	"github.com/nomadsingles/platform/utils"
)

type Context struct {
	Session       model.Session
	RequestId     string
	IpAddress     string
	Path          string
	Err           *model.AppError
	siteURLHeader string
	T             goi18n.TranslateFunc
	Locale        string
}

func ApiAppHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &handler{h, false, false, true}
}

func AppHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &handler{h, false, false, false}
}

func AppHandlerIndependent(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &handler{h, false, false, false}
}

func ApiUserRequired(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &handler{h, true, false, true}
}

func UserRequired(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &handler{h, true, false, false}
}

func AppHandlerTrustRequester(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &handler{h, false, false, false}
}

func ApiAdminSystemRequired(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &handler{h, true, true, true}
}

func ApiAdminSystemRequiredTrustRequester(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &handler{h, true, true, true}
}

func ApiAppHandlerTrustRequester(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &handler{h, false, false, true}
}

func ApiUserRequiredTrustRequester(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &handler{h, true, false, true}
}

func ApiAppHandlerTrustRequesterIndependent(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &handler{h, false, false, true}
}

type handler struct {
	handleFunc         func(*Context, http.ResponseWriter, *http.Request)
	requireUser        bool
	requireSystemAdmin bool
	isApi              bool
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	l4g.Debug("%v", r.URL.Path)



	c := &Context{}
	c.T, c.Locale = utils.GetTranslationsAndLocale(w, r)
	c.RequestId = model.NewId()
	c.IpAddress = utils.GetIpAddress(r)

	token := ""


	// Attempt to parse token out of the header
	authHeader := r.Header.Get(model.HEADER_AUTH)
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:6]) == model.HEADER_BEARER {
		// Default session token
		token = authHeader[7:]

	} else if len(authHeader) > 5 && strings.ToLower(authHeader[0:5]) == model.HEADER_TOKEN {
		// OAuth token
		token = authHeader[6:]
	}

	// Attempt to parse the token from the cookie
	if len(token) == 0 {
		if cookie, err := r.Cookie(model.SESSION_COOKIE_TOKEN); err == nil {
			token = cookie.Value

			if (h.requireSystemAdmin || h.requireUser) {
				if r.Header.Get(model.HEADER_REQUESTED_WITH) != model.HEADER_REQUESTED_WITH_XML {
					c.Err = model.NewLocAppError("ServeHTTP", "api.context.session_expired.app_error", nil, "token="+token+" Appears to be a CSRF attempt")
					token = ""
				}
			}
		}
	}



	c.SetSiteURLHeader(app.GetProtocol(r) + "://" + r.Host)

	w.Header().Set(model.HEADER_REQUEST_ID, c.RequestId)
	w.Header().Set(model.HEADER_VERSION_ID, fmt.Sprintf("%v.%v.%v.%v", model.CurrentVersion, model.BuildNumber, utils.ClientCfgHash))


	// Instruct the browser not to display us in an iframe unless is the same origin for anti-clickjacking
	if !h.isApi {
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'self'")
	} else {
		// All api response bodies will be JSON formatted by default
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "GET" {
			w.Header().Set("Expires", "0")
		}
	}


	if c.Err == nil && (h.requireUser || h.requireSystemAdmin) {
		//check if teamId exist
		c.CheckTeamId()
	}

	if c.Err == nil {
		h.handleFunc(c, w, r)
	}

	// Handle errors that have occoured
	if c.Err != nil {
		c.Err.Translate(c.T)
		c.Err.RequestId = c.RequestId
		c.LogError(c.Err)
		c.Err.Where = r.URL.Path

		// Block out detailed error when not in developer mode
		if !*utils.Cfg.ServiceSettings.EnableDeveloper {
			c.Err.DetailedError = ""
		}

		if h.isApi {
			w.WriteHeader(c.Err.StatusCode)
			w.Write([]byte(c.Err.ToJson()))


		}

	}


	}




func (c *Context) LogError(err *model.AppError) {

	// filter out endless reconnects
	if c.Path == "/api/v3/users/websocket" && err.StatusCode == 401 || err.Id == "web.check_browser_compatibility.app_error" {
		c.LogDebug(err)
	} else {
		l4g.Error(utils.T("api.context.log.error"), c.Path, err.Where, err.StatusCode,
			c.RequestId, c.Session.UserId, c.IpAddress, err.SystemMessage(utils.T), err.DetailedError)
	}
}

func (c *Context) LogDebug(err *model.AppError) {
	l4g.Debug(utils.T("api.context.log.error"), c.Path, err.Where, err.StatusCode,
		c.RequestId, c.Session.UserId, c.IpAddress, err.SystemMessage(utils.T), err.DetailedError)
}


func (c *Context) SetInvalidParam(where string, name string) {
	c.Err = NewInvalidParamError(where, name)
}

func NewInvalidParamError(where string, name string) *model.AppError {
	err := model.NewLocAppError(where, "api.context.invalid_param.app_error", map[string]interface{}{"Name": name}, "")
	err.StatusCode = http.StatusBadRequest
	return err
}

func (c *Context) SetUnknownError(where string, details string) {
	c.Err = model.NewLocAppError(where, "api.context.unknown.app_error", nil, details)
}




func (c *Context) SetSiteURLHeader(url string) {
	c.siteURLHeader = strings.TrimRight(url, "/")
}




func (c *Context) GetSiteURLHeader() string {
	return c.siteURLHeader
}



func IsApiCall(r *http.Request) bool {
	return strings.Index(r.URL.Path, "/api/") == 0
}

func RenderWebError(err *model.AppError, w http.ResponseWriter, r *http.Request) {
	T, _ := utils.GetTranslationsAndLocale(w, r)

	title := T("api.templates.error.title", map[string]interface{}{"SiteName": utils.ClientCfg["SiteName"]})
	message := err.Message
	details := err.DetailedError
	link := "/"
	linkMessage := T("api.templates.error.link")

	status := http.StatusTemporaryRedirect
	if err.StatusCode != http.StatusInternalServerError {
		status = err.StatusCode
	}

	http.Redirect(
		w,
		r,
		"/error?title="+url.QueryEscape(title)+
			"&message="+url.QueryEscape(message)+
			"&details="+url.QueryEscape(details)+
			"&link="+url.QueryEscape(link)+
			"&linkmessage="+url.QueryEscape(linkMessage),
		status)
}

func Handle404(w http.ResponseWriter, r *http.Request) {
	err := model.NewLocAppError("Handle404", "api.context.404.app_error", nil, "")
	err.Translate(utils.T)
	err.StatusCode = http.StatusNotFound

	l4g.Debug("%v: code=404 ip=%v", r.URL.Path, utils.GetIpAddress(r))

	if IsApiCall(r) {
		w.WriteHeader(err.StatusCode)
		err.DetailedError = "There doesn't appear to be an api call for the url='" + r.URL.Path + "'.  Typo? are you missing a team_id or user_id as part of the url?"
		w.Write([]byte(err.ToJson()))
	} else {
		RenderWebError(err, w, r)
	}
}

func (c *Context) CheckTeamId() {

}
