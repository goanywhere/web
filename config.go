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
package rex

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/goanywhere/rex/crypto"
	"github.com/goanywhere/rex/template"
	"github.com/goanywhere/x/env"
)

var (
	once   sync.Once
	loader *template.Loader

	settings  *config
	signature *crypto.Signature
)

type config struct {
	Root   string
	Debug  bool
	Secret string

	Host string
	Port int

	Templates string

	X_Frame_Options        string
	X_Content_Type_Options string
	X_XSS_Protection       string
	X_UA_Compatible        string
}

// Settings returns a singleton settings access point.
func configure() *config {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to retrieve project root: %v", err)
	}

	once.Do(func() {
		settings = new(config)
		settings.Debug = true
		settings.Host = "localhost"
		settings.Port = 5000

		settings.Templates = "templates"

		settings.X_Frame_Options = "deny"
		settings.X_Content_Type_Options = "nosniff"
		settings.X_XSS_Protection = "1; mode=block"
		settings.X_UA_Compatible = "IE=Edge, chrome=1"

		settings.Root, _ = filepath.Abs(cwd)
		env.Load(filepath.Join(settings.Root, ".env"))
		env.Dump(settings)

		if settings.Secret == "" {
			log.Fatal("Secret key missing")
		}

		loader = template.NewLoader(settings.Templates)
		// creates a signature for accessing securecookie.
		signature = crypto.NewSignature(settings.Secret)
	})
	return settings
}