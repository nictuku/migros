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
	"os"
	"flag"
	"log"
)

var username = flag.String("u", "nictuku", "username")
var password = flag.String("p", "", "password")

func requireLogin() {
	if *username == "" || *password == "" {
		fmt.Fprintf(os.Stderr, "\nERROR: Username or password missing.\n\n")
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	requireLogin()
	client := &minecraftClient{}
	err := client.Login(*username, *password)
	if err != nil {
		log.Println(err)
	}
}
