// Copyright 2011 Yves Junqueira
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"http"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/garyburd/twister/web"
)

const (
	MINECRAFT_LOGIN_URL   = "https://www.minecraft.net/login.jsp"
	MINECRAFT_NET_TIMEOUT = 10
)

type minecraftClient struct{}

// Login to the minecraft.net website using HTTPS.
func (c *minecraftClient) Login(username, password string) (err os.Error) {
	param := make(web.ParamMap)
	param.Set("username", username)
	param.Set("password", password)
	param.Set("use_secure", "true")
	url := MINECRAFT_LOGIN_URL + "?" + param.FormEncodedString()

	var resp *http.Response
	done := make(chan bool, 1)
	go func() {
		resp, err = http.PostForm(url, param.StringMap())
		done <- true
	}()

	timeout := time.After(MINECRAFT_NET_TIMEOUT * 1e9) //
	select {
	case <-done:
		break
	case <-timeout:
		return os.NewError("Login attempt timed out - " + MINECRAFT_LOGIN_URL)
	}
	if resp == nil {
		return os.NewError("Login server responded with a null body.")
	}

	switch resp.StatusCode {
	// A successful minecraft.net login is followed by a 302.
	case 302:
		if r, ok := resp.Header["Location"]; ok && len(r) > 0 {
			if r[0] == "https://www.minecraft.net/" {
				log.Println("Login successful.")
				return nil
			} else {
				return os.NewError("Login redirecting to unknown page: " + r[0])
			}
		}
	case 200:
		return os.NewError("Login failed.")
	}
	return os.NewError("Login return code: " + string(resp.StatusCode))
}

func readHttpResponse(resp *http.Response, httpErr os.Error) (p []byte, err os.Error) {
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	err = httpErr
	if err != nil {
		log.Println(err.String())
		return nil, err
	}
	p, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		err = os.NewError(fmt.Sprintf("Server Error code: %d; msg: %v", resp.StatusCode, string(p)))
		return nil, err
	}
	return p, nil
}
