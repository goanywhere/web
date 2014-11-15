/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       create.go
 *  @date       2014-11-02
 *  @author     Jim Zhan <jim.zhan@me.com>
 *
 *  Copyright © 2014 Jim Zhan.
 *  ------------------------------------------------------------
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *  ------------------------------------------------------------
 */
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/goanywhere/web"
	"github.com/goanywhere/web/crypto"
)

var pool = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*(-_+)")

func createProject(path string) error {
	// create .env with secret.
	if env, err := os.Create(filepath.Join(path, ".env")); err == nil {
		defer env.Close()
		buffer := bufio.NewWriter(env)
		buffer.WriteString(fmt.Sprintf("secret=%s\n", crypto.RandomString(64, pool)))
		buffer.Flush()
	}
	return nil
}

func Create(context *cli.Context) {
	args := context.Args()
	if len(args) != 1 {
		web.Info("Valid Project Name Missing")
	} else {
		if cwd, err := os.Getwd(); err != nil {
			panic(err)
		} else {
			path := filepath.Join(cwd, args[0])
			if err := os.Mkdir(path, os.ModePerm); err != nil {
				panic(err)
			}
			createProject(path)
		}
	}
}
