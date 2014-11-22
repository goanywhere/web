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

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goanywhere/fs"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateProject(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	workspace := filepath.Join(strings.Split(gopath, ";")[0], "src")

	Convey("[CMD] Create Project", t, func() {
		path := "github.com/goanywhere/CreateProjectTest"
		project, err := createProject(path)
		So(project, ShouldEqual, filepath.Join(workspace, path))
		So(fs.Exists(project), ShouldBeTrue)
		So(err, ShouldBeNil)
		os.RemoveAll(project)
	})

	Convey("[CMD] Create Project with Invalid Paths", t, func() {
		path := "github.com/goanywhere/web"
		project, err := createProject(path)
		So(project, ShouldEqual, "")
		So(err, ShouldEqual, os.ErrExist)
	})

}