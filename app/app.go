// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


package app

import (
	"io/ioutil"
	"net/http"
)

func CloseBody(r *http.Response) {
	if r.Body != nil {
		ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
}
