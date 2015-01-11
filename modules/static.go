/* ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 *
 * (C) Copyright 2015 GoAnywhere (http://goanywhere.io).
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
package modules

import (
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/goanywhere/rex"
	"github.com/goanywhere/rex/template"
)

type static struct {
	Dir  string
	URL  string
	once sync.Once
}

func (self *static) init(options Options) {
	self.Dir = options.Get("Dir", "build").(string)
	self.URL = options.Get("URL", "/static/").(string)

	self.once.Do(func() {
		template.Functions["static"] = func(path string) string {
			return strings.Join([]string{
				strings.TrimRight(self.URL, "/"),
				strings.TrimLeft(path, "/")},
				"/")
		}
	})
}

func (self *static) serve(w http.ResponseWriter, r *http.Request) {
	var dir = http.Dir(filepath.Join(rex.Settings.Root, self.Dir))
	var path = r.URL.Path[len(self.URL):]

	var file, err = dir.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return
	}
	// try serving index.html
	if stat.IsDir() {
		// redirect if missing trailing slash
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusFound)
			return
		}

		path = filepath.Join(path, "index.html")
		file, err = dir.Open(path)
		if err != nil {
			return
		}
		defer file.Close()

		stat, err = file.Stat()
		if err != nil || stat.IsDir() {
			return
		}
	}

	http.ServeContent(w, r, path, stat.ModTime(), file)
}

func Static(options Options) func(http.Handler) http.Handler {
	s := new(static)
	s.init(options)

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" && strings.HasPrefix(r.URL.Path, s.URL) {
				s.serve(w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)
	}
}