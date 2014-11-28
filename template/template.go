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

package template

import (
	"html/template"
	"io/ioutil"
)

// Parse finds all extends chain & constructs the final page layout.
func Parse(filename string) (page *template.Template) {
	var err error
	filenames := Extends(filename)

	for _, item := range filenames {
		if bits, err := ioutil.ReadFile(item); err == nil {
			var tmpl *template.Template
			// intialize final page template using the very first ancestor.
			if page == nil {
				page = template.New(item)
			}
			if item == page.Name() {
				tmpl = page
			} else {
				tmpl = page.New(item)
			}
			_, err = tmpl.Parse(string(bits))
		}
	}
	return template.Must(page, err)
}
