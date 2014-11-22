/* ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 *
 * (C) Copyright 2014 GoAnywhere (http://goanywhere.io).
 * ----------------------------------------------------------------------
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
 * ----------------------------------------------------------------------*/

package web

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/goanywhere/env"
	"github.com/gorilla/mux"
)

// Shortcut to create map.
type H map[string]interface{}

// New creates an application instance & setup its default settings..
func New() *Mux {
	app := &Mux{mux.NewRouter(), nil}
	return app
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if cwd, err := os.Getwd(); err == nil {
		if root, err := filepath.Abs(cwd); err == nil {
			env.Set("root", root)
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
	// Application Defaults
	env.Set("debug", "true")
	env.Set("host", "localhost")
	env.Set("port", "5000")
	env.Set("templates", "templates")
}