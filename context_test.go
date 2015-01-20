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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// ---------------------------------------------------------------------------
//  Enhancements for native http.ResponseWriter
// ---------------------------------------------------------------------------
func TestContextStatus(t *testing.T) {
	Convey("Response Status Code", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w, r)
			ctx.String("200 Response")
		}))
		defer server.Close()

		if response, err := http.Get(server.URL); err == nil {
			So(response.StatusCode, ShouldEqual, http.StatusOK)
		}

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w, r)
			ctx.WriteHeader(http.StatusNotFound)
			ctx.String("404 Response")
		}))
		defer server.Close()
		if response, err := http.Get(server.URL); err == nil {
			So(response.StatusCode, ShouldEqual, http.StatusNotFound)
		}
	})
}

func TestContextSize(t *testing.T) {
	Convey("Response Size", t, func() {
		value := "Hello 中文測試"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w, r)
			ctx.String(value)
		}))
		defer server.Close()

		if response, err := http.Get(server.URL); err == nil {
			body, _ := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			So(len(body), ShouldEqual, len([]byte(value)))
		}

	})
}

func TestContextWritten(t *testing.T) {
	Convey("Response's Written Flag", t, func() {
		var flag bool
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w, r)
			flag = ctx.Written()
			ctx.String("Hello World")
		}))
		defer server.Close()

		if _, err := http.Get(server.URL); err == nil {
			So(flag, ShouldBeFalse)
		}
	})
}

// ---------------------------------------------------------------------------
//  HTTP Request Context Data
// ---------------------------------------------------------------------------
/*
func TestContextId(t *testing.T) {
	Convey("Unique Context Id", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w, r)
			ctx.String(ctx.Id())
		}))
		defer server.Close()

		var a, b string
		if response, err := http.Get(server.URL); err == nil {
			body, _ := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			a = string(body)
			So(len(a), ShouldEqual, 40)
			So(strings.HasSuffix(a, "1"), ShouldBeTrue)
		}

		if response, err := http.Get(server.URL); err == nil {
			body, _ := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			b = string(body)
			So(len(b), ShouldEqual, 40)
			So(strings.HasSuffix(b, "2"), ShouldBeTrue)
		}

		So(a[0:len(a)-1], ShouldEqual, b[0:len(b)-1])
	})
}
*/

func TestContextGet(t *testing.T) {
	Convey("Context Data Get", t, func() {
		var value = "example"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w, r)
			ctx.data["name"] = value
			ctx.String(ctx.Get("name").(string))
		}))
		defer server.Close()

		if response, err := http.Get(server.URL); err == nil {
			body, _ := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			So(string(body), ShouldEqual, "example")
		}
	})
}

func TestContextSet(t *testing.T) {
	Convey("Context Data Set", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w, r)
			ctx.Set("name", "example")
			ctx.String(ctx.Get("name").(string))
		}))
		defer server.Close()

		if response, err := http.Get(server.URL); err == nil {
			body, _ := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			So(string(body), ShouldEqual, "example")
		}
	})
}

func TestContextClear(t *testing.T) {
	Convey("Context Data Clear", t, func() {
		var a, b int
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w, r)
			ctx.Set("name", "example")
			a = len(ctx.data)
			ctx.Clear()
			b = len(ctx.data)
			ctx.String("DONE")
		}))
		defer server.Close()

		if response, err := http.Get(server.URL); err == nil {
			body, _ := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			So(a, ShouldEqual, 1)
			So(b, ShouldEqual, 0)
			So(response.StatusCode, ShouldEqual, http.StatusOK)
			So(string(body), ShouldEqual, "DONE")
		}
	})
}

func TestContextDelete(t *testing.T) {
	Convey("Context Data Delete", t, func() {
		var a, b int
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w, r)
			ctx.Set("name", "example")
			a = len(ctx.data)
			ctx.Delete("name")
			b = len(ctx.data)
			ctx.String("DONE")
		}))
		defer server.Close()

		if response, err := http.Get(server.URL); err == nil {
			body, _ := ioutil.ReadAll(response.Body)
			defer response.Body.Close()

			So(a, ShouldEqual, 1)
			So(b, ShouldEqual, 0)
			So(response.StatusCode, ShouldEqual, http.StatusOK)
			So(string(body), ShouldEqual, "DONE")
		}
	})
}

// ---------------------------------------------------------------------------
//  HTTP Cookies
// ---------------------------------------------------------------------------
func TestCookie(t *testing.T) {
	Convey("context#Cookie", t, func() {
		cookie := &http.Cookie{Name: "number", Value: "123", Path: "/"}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w, r)
			ctx.String(ctx.Cookie(cookie.Name))
		}))
		defer server.Close()

		client := new(http.Client)
		request, _ := http.NewRequest("GET", server.URL, nil)
		request.AddCookie(cookie)

		response, _ := client.Do(request)
		defer response.Body.Close()

		body, _ := ioutil.ReadAll(response.Body)
		So(string(body), ShouldEqual, cookie.Value)
	})
}

func TestSetCookie(t *testing.T) {
	Convey("context#SetCookie", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w, r)
			ctx.SetCookie(&http.Cookie{Name: "number", Value: "123", Path: "/"})
			ctx.String("Hello Cookie")
			return
		}))
		defer server.Close()

		if response, err := http.Get(server.URL); err == nil {
			cookie := response.Cookies()[0]
			So(cookie.Name, ShouldEqual, "number")
			So(cookie.Value, ShouldEqual, "123")
		}
	})
}

// ---------------------------------------------------------------------------
//  HTTP Response Rendering
// ---------------------------------------------------------------------------
func TestContextHTML(t *testing.T) {
	Convey("Rendering HTML", t, func() {

	})
}

func TestContextJSON(t *testing.T) {
	Convey("Rendering JSON", t, func() {

	})
}

func TestContextXML(t *testing.T) {
	Convey("Rendering XML", t, func() {

	})
}

func TestContextString(t *testing.T) {
	Convey("Rendering String", t, func() {

	})
}