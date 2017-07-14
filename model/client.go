// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.


package model

import (

	"fmt"

	"io/ioutil"

	"net/http"
	"net/url"

	"strings"
	"time"

	l4g "github.com/alecthomas/log4go"
)

const (
	HEADER_REQUEST_ID         = "X-Request-ID"
	HEADER_VERSION_ID         = "X-Version-ID"
	HEADER_CLUSTER_ID         = "X-Cluster-ID"
	HEADER_ETAG_SERVER        = "ETag"
	HEADER_ETAG_CLIENT        = "If-None-Match"
	HEADER_FORWARDED          = "X-Forwarded-For"
	HEADER_REAL_IP            = "X-Real-IP"
	HEADER_FORWARDED_PROTO    = "X-Forwarded-Proto"
	HEADER_TOKEN              = "token"
	HEADER_BEARER             = "BEARER"
	HEADER_AUTH               = "Authorization"
	HEADER_REQUESTED_WITH     = "X-Requested-With"
	HEADER_REQUESTED_WITH_XML = "XMLHttpRequest"
	STATUS                    = "status"
	STATUS_OK                 = "OK"
	STATUS_FAIL               = "FAIL"
	STATUS_REMOVE             = "REMOVE"

	CLIENT_DIR = "webapp/dist"

	API_URL_SUFFIX_V1 = "/api/v1"
	API_URL_SUFFIX_V3 = "/api/v3"
	API_URL_SUFFIX_V4 = "/api/v4"
	API_URL_SUFFIX    = API_URL_SUFFIX_V4
)

type Result struct {
	RequestId string
	Etag      string
	Data      interface{}
}

type ResponseMetadata struct {
	StatusCode int
	Error      *AppError
	RequestId  string
	Etag       string
}

type Client struct {
	Url           string       // The location of the server like "http://localhost:8065"
	ApiUrl        string       // The api location of the server like "http://localhost:8065/api/v3"
	HttpClient    *http.Client // The http client
	AuthToken     string
	AuthType      string
	TeamId        string
	RequestId     string
	Etag          string
	ServerVersion string
}

// NewClient constructs a new client with convienence methods for talking to
// the server.
func NewClient(url string) *Client {
	return &Client{url, url + API_URL_SUFFIX_V3, &http.Client{}, "", "", "", "", "", ""}
}

func closeBody(r *http.Response) {
	if r.Body != nil {
		ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
}

func (c *Client) SetOAuthToken(token string) {
	c.AuthToken = token
	c.AuthType = HEADER_TOKEN
}

func (c *Client) ClearOAuthToken() {
	c.AuthToken = ""
	c.AuthType = HEADER_BEARER
}

func (c *Client) SetTeamId(teamId string) {
	c.TeamId = teamId
}

func (c *Client) GetTeamId() string {
	if len(c.TeamId) == 0 {
		println(`You are trying to use a route that requires a team_id, 
        	but you have not called SetTeamId() in client.go`)
	}

	return c.TeamId
}

func (c *Client) ClearTeamId() {
	c.TeamId = ""
}

func (c *Client) GetTeamRoute() string {
	return fmt.Sprintf("/teams/%v", c.GetTeamId())
}

func (c *Client) GetChannelRoute(channelId string) string {
	return fmt.Sprintf("/teams/%v/channels/%v", c.GetTeamId(), channelId)
}

func (c *Client) GetUserRequiredRoute(userId string) string {
	return fmt.Sprintf("/users/%v", userId)
}

func (c *Client) GetChannelNameRoute(channelName string) string {
	return fmt.Sprintf("/teams/%v/channels/name/%v", c.GetTeamId(), channelName)
}

func (c *Client) GetEmojiRoute() string {
	return "/emoji"
}

func (c *Client) GetGeneralRoute() string {
	return "/general"
}

func (c *Client) GetFileRoute(fileId string) string {
	return fmt.Sprintf("/files/%v", fileId)
}

func (c *Client) DoPost(url, data, contentType string) (*http.Response, *AppError) {
	rq, _ := http.NewRequest("POST", c.Url+url, strings.NewReader(data))
	rq.Header.Set("Content-Type", contentType)
	rq.Close = true

	if rp, err := c.HttpClient.Do(rq); err != nil {
		return nil, NewLocAppError(url, "model.client.connecting.app_error", nil, err.Error())
	} else if rp.StatusCode >= 300 {
		defer closeBody(rp)
		return nil, AppErrorFromJson(rp.Body)
	} else {
		return rp, nil
	}
}

func (c *Client) DoApiPost(url string, data string) (*http.Response, *AppError) {
	rq, _ := http.NewRequest("POST", c.ApiUrl+url, strings.NewReader(data))
	rq.Close = true

	if len(c.AuthToken) > 0 {
		rq.Header.Set(HEADER_AUTH, c.AuthType+" "+c.AuthToken)
	}

	if rp, err := c.HttpClient.Do(rq); err != nil {
		return nil, NewLocAppError(url, "model.client.connecting.app_error", nil, err.Error())
	} else if rp.StatusCode >= 300 {
		defer closeBody(rp)
		return nil, AppErrorFromJson(rp.Body)
	} else {
		return rp, nil
	}
}

func (c *Client) DoApiGet(url string, data string, etag string) (*http.Response, *AppError) {
	rq, _ := http.NewRequest("GET", c.ApiUrl+url, strings.NewReader(data))
	rq.Close = true

	if len(etag) > 0 {
		rq.Header.Set(HEADER_ETAG_CLIENT, etag)
	}

	if len(c.AuthToken) > 0 {
		rq.Header.Set(HEADER_AUTH, c.AuthType+" "+c.AuthToken)
	}

	if rp, err := c.HttpClient.Do(rq); err != nil {
		return nil, NewLocAppError(url, "model.client.connecting.app_error", nil, err.Error())
	} else if rp.StatusCode == 304 {
		return rp, nil
	} else if rp.StatusCode >= 300 {
		defer closeBody(rp)
		return rp, AppErrorFromJson(rp.Body)
	} else {
		return rp, nil
	}
}

func getCookie(name string, resp *http.Response) *http.Cookie {
	for _, cookie := range resp.Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}

	return nil
}

// Must is a convenience function used for testing.
func (c *Client) Must(result *Result, err *AppError) *Result {
	if err != nil {
		l4g.Close()
		time.Sleep(time.Second)
		panic(err)
	}

	return result
}

// MustGeneric is a convenience function used for testing.
func (c *Client) MustGeneric(result interface{}, err *AppError) interface{} {
	if err != nil {
		l4g.Close()
		time.Sleep(time.Second)
		panic(err)
	}

	return result
}

// CheckStatusOK is a convenience function for checking the return of Web Service
// call that return the a map of status=OK.
func (c *Client) CheckStatusOK(r *http.Response) bool {
	m := MapFromJson(r.Body)
	defer closeBody(r)

	if m != nil && m[STATUS] == STATUS_OK {
		return true
	}

	return false
}

func (c *Client) fillInExtraProperties(r *http.Response) {
	c.RequestId = r.Header.Get(HEADER_REQUEST_ID)
	c.Etag = r.Header.Get(HEADER_ETAG_SERVER)
	c.ServerVersion = r.Header.Get(HEADER_VERSION_ID)
}

func (c *Client) clearExtraProperties() {
	c.RequestId = ""
	c.Etag = ""
	c.ServerVersion = ""
}

// General Routes Section

// GetClientProperties returns properties needed by the client to show/hide
// certian features.  It returns a map of strings.
func (c *Client) GetClientProperties() (map[string]string, *AppError) {
	c.clearExtraProperties()
	if r, err := c.DoApiGet(c.GetGeneralRoute()+"/client_props", "", ""); err != nil {
		return nil, err
	} else {
		defer closeBody(r)
		c.fillInExtraProperties(r)
		return MapFromJson(r.Body), nil
	}
}

// LogClient is a convenience Web Service call so clients can log messages into
// the server-side logs.  For example we typically log javascript error messages
// into the server-side.  It returns true if the logging was successful.
func (c *Client) LogClient(message string) (bool, *AppError) {
	c.clearExtraProperties()
	m := make(map[string]string)
	m["level"] = "ERROR"
	m["message"] = message

	if r, err := c.DoApiPost(c.GetGeneralRoute()+"/log_client", MapToJson(m)); err != nil {
		return false, err
	} else {
		defer closeBody(r)
		c.fillInExtraProperties(r)
		return c.CheckStatusOK(r), nil
	}
}

// GetPing returns a map of strings with server time, server version, and node Id.
// Systems that want to check on health status of the server should check the
// url /api/v3/ping for a 200 status response.
func (c *Client) GetPing() (map[string]string, *AppError) {
	c.clearExtraProperties()
	if r, err := c.DoApiGet(c.GetGeneralRoute()+"/ping", "", ""); err != nil {
		return nil, err
	} else {
		defer closeBody(r)
		c.fillInExtraProperties(r)
		return MapFromJson(r.Body), nil
	}
}








// User Routes Section

// CreateUser creates a user in the system based on the provided user struct.
func (c *Client) CreateUser(user *User, hash string) (*Result, *AppError) {
	if r, err := c.DoApiPost("/users/create", user.ToJson()); err != nil {
		return nil, err
	} else {
		defer closeBody(r)
		return &Result{r.Header.Get(HEADER_REQUEST_ID),
			r.Header.Get(HEADER_ETAG_SERVER), UserFromJson(r.Body)}, nil
	}
}


func (c *Client) CreateUserFromSignup(user *User, data string, hash string) (*Result, *AppError) {
	if r, err := c.DoApiPost("/users/create?d="+url.QueryEscape(data)+"&h="+hash, user.ToJson()); err != nil {
		return nil, err
	} else {
		defer closeBody(r)
		return &Result{r.Header.Get(HEADER_REQUEST_ID),
			r.Header.Get(HEADER_ETAG_SERVER), UserFromJson(r.Body)}, nil
	}
}

// GetUser returns a user based on a provided user id string. Must be authenticated.
func (c *Client) GetUser(id string, etag string) (*Result, *AppError) {
	if r, err := c.DoApiGet("/users/"+id+"/get", "", etag); err != nil {
		return nil, err
	} else {
		defer closeBody(r)
		return &Result{r.Header.Get(HEADER_REQUEST_ID),
			r.Header.Get(HEADER_ETAG_SERVER), UserFromJson(r.Body)}, nil
	}
}

// getByUsername returns a user based on a provided username string. Must be authenticated.
func (c *Client) GetByUsername(username string, etag string) (*Result, *AppError) {
	if r, err := c.DoApiGet(fmt.Sprintf("/users/name/%v", username), "", etag); err != nil {
		return nil, err
	} else {
		defer closeBody(r)
		return &Result{r.Header.Get(HEADER_REQUEST_ID),
			r.Header.Get(HEADER_ETAG_SERVER), UserFromJson(r.Body)}, nil
	}
}

// LoginById authenticates a user by user id and password.
func (c *Client) LoginById(id string, password string) (*Result, *AppError) {
	m := make(map[string]string)
	m["id"] = id
	m["password"] = password
	return c.login(m)
}

// Login authenticates a user by login id, which can be username, email or some sort
// of SSO identifier based on configuration, and a password.
func (c *Client) Login(loginId string, password string) (*Result, *AppError) {
	m := make(map[string]string)
	m["login_id"] = loginId
	m["password"] = password
	return c.login(m)
}




func (c *Client) login(m map[string]string) (*Result, *AppError) {
	if r, err := c.DoApiPost("/users/login", MapToJson(m)); err != nil {
		return nil, err
	} else {
		c.AuthToken = r.Header.Get(HEADER_TOKEN)
		c.AuthType = HEADER_BEARER
		sessionToken := getCookie(SESSION_COOKIE_TOKEN, r)

		if c.AuthToken != sessionToken.Value {
			NewLocAppError("/users/login", "model.client.login.app_error", nil, "")
		}

		defer closeBody(r)
		return &Result{r.Header.Get(HEADER_REQUEST_ID),
			r.Header.Get(HEADER_ETAG_SERVER), UserFromJson(r.Body)}, nil
	}
}

// Logout terminates the current user's session.
func (c *Client) Logout() (*Result, *AppError) {
	if r, err := c.DoApiPost("/users/logout", ""); err != nil {
		return nil, err
	} else {
		c.AuthToken = ""
		c.AuthType = HEADER_BEARER
		c.TeamId = ""

		defer closeBody(r)
		return &Result{r.Header.Get(HEADER_REQUEST_ID),
			r.Header.Get(HEADER_ETAG_SERVER), MapFromJson(r.Body)}, nil
	}
}





func (c *Client) RevokeSession(sessionAltId string) (*Result, *AppError) {
	m := make(map[string]string)
	m["id"] = sessionAltId

	if r, err := c.DoApiPost("/users/revoke_session", MapToJson(m)); err != nil {
		return nil, err
	} else {
		defer closeBody(r)
		return &Result{r.Header.Get(HEADER_REQUEST_ID),
			r.Header.Get(HEADER_ETAG_SERVER), MapFromJson(r.Body)}, nil
	}
}

func (c *Client) GetSessions(id string) (*Result, *AppError) {
	if r, err := c.DoApiGet("/users/"+id+"/sessions", "", ""); err != nil {
		return nil, err
	} else {
		defer closeBody(r)
		return &Result{r.Header.Get(HEADER_REQUEST_ID),
			r.Header.Get(HEADER_ETAG_SERVER), SessionsFromJson(r.Body)}, nil
	}
}
