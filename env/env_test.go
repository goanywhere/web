/**
 * ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 * ----------------------------------------------------------------------
 *  Copyright © 2014 GoAnywhere Ltd. All Rights Reserved.
 * ----------------------------------------------------------------------*/

package env

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/goanywhere/web/crypto"
	. "github.com/smartystreets/goconvey/convey"
)

type Spec struct {
	App     string
	Debug   bool
	Total   int
	Version float32
	Tag     string `web:"multiple_words_tag"`
}

func setup() {
	os.Clearenv()
	Set("app", "example")
	Set("debug", "true")
	Set("total", "100")
	Set("version", "32.1")
	Set("multiple_words_tag", "ALT")
}

func TestFindKeyValue(t *testing.T) {
	Convey("[private] Find key value pair from string", t, func() {
		k, v := findKeyValue(" test: value")
		So(k, ShouldEqual, "test")
		So(v, ShouldEqual, "value")

		k, v = findKeyValue(" test: value")
		So(k, ShouldEqual, "test")
		So(v, ShouldEqual, "value")

		k, v = findKeyValue("\ttest:\tvalue\t\n")
		So(k, ShouldEqual, "test")
		So(v, ShouldEqual, "value")
	})
}

func TestLoad(t *testing.T) {
	var pool = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*(-_+)")
	filename := "/tmp/.env"
	// plain value without quote
	if env, err := os.Create(filename); err == nil {
		defer env.Close()
		defer os.Remove(filename)
		secret := crypto.RandomString(64, pool)
		buffer := bufio.NewWriter(env)
		buffer.WriteString(fmt.Sprintf("secret=%s\n", secret))
		buffer.WriteString("app=myapp\n")
		buffer.Flush()

		Convey("Load key/value from dotenv file", t, func() {
			Set("root", "/tmp")
			Load()
			So(Get("secret"), ShouldEqual, secret)
		})
	}
	// value with `` quote
	if env, err := os.Create(filename); err == nil {
		defer env.Close()
		defer os.Remove(filename)
		secret := crypto.RandomString(64, pool)
		buffer := bufio.NewWriter(env)
		buffer.WriteString(fmt.Sprintf("secret='%s'\n", secret))
		buffer.WriteString("app=myapp\n")
		buffer.Flush()

		Convey("Load key/value from dotenv file with \"'\"", t, func() {
			Set("root", "/tmp")
			Load()
			So(Get("secret"), ShouldEqual, secret)
		})
	}
	// value with `"` quote
	if env, err := os.Create(filename); err == nil {
		defer env.Close()
		defer os.Remove(filename)
		secret := crypto.RandomString(64, pool)
		buffer := bufio.NewWriter(env)
		buffer.WriteString(fmt.Sprintf("secret=\"%s\"\n", secret))
		buffer.WriteString("app=myapp\n")
		buffer.Flush()

		Convey("Load key/value from dotenv file with '\"'", t, func() {
			Set("root", "/tmp")
			Load()
			So(Get("secret"), ShouldEqual, secret)
		})
	}
}

func TestLoadSpec(t *testing.T) {
	Convey("Load key/value from `os.Envrion in pre-defined struct`", t, func() {
		setup()
		var spec Spec
		err := LoadSpec(&spec)
		So(err, ShouldBeNil)
		So(spec.App, ShouldEqual, "example")
		So(spec.Debug, ShouldBeTrue)
		So(spec.Total, ShouldEqual, 100)
		So(spec.Version, ShouldEqual, 32.1)
		So(spec.Tag, ShouldEqual, "ALT")

		Set("app", "myapplication")
		LoadSpec(&spec)
		So(spec.App, ShouldEqual, "myapplication")
	})
}

func TestGetString(t *testing.T) {
	Convey("GetString from os.Environ", t, func() {
		secret := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*(-_=+)"
		Set("secret", secret)
		So(Get("secret"), ShouldEqual, secret)
		So(Get("SeCrEt"), ShouldEqual, secret)
	})
}

func TestGetBool(t *testing.T) {
	Convey("GetBool from os.Environ", t, func() {
		value, err := GetBool("NotFound")
		So(value, ShouldBeFalse)
		So(err, ShouldNotBeNil)

		Set("enabled", "true")
		value, err = GetBool("enabled")
		So(value, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}

func TestGetInt(t *testing.T) {
	Convey("GetInt from os.Environ", t, func() {
		value, err := GetInt("NotFound")
		So(value, ShouldEqual, 0)
		So(err, ShouldNotBeNil)

		Set("int", "123")
		value, err = GetInt("int")
		So(value, ShouldEqual, 123)
		So(err, ShouldBeNil)
	})
}

func TestGetFloat(t *testing.T) {
	Convey("GetFloat from os.Environ", t, func() {
		value, err := GetFloat("NotFound")
		So(value, ShouldEqual, 0.0)
		So(err, ShouldNotBeNil)

		Set("float", "32.1")
		value, err = GetFloat("float")
		So(value, ShouldEqual, 32.1)
		So(err, ShouldBeNil)
	})
}

func TestAccess(t *testing.T) {
	Convey("Get/Set access to os.Environ", t, func() {
		Set("shell", "/bin/zsh")
		So(Get("shell"), ShouldEqual, "/bin/zsh")

		Set("AnyThiNg", "content")
		So(Get("anything"), ShouldEqual, "content")

		So(Get("NotFound"), ShouldEqual, "")
	})
}

func TestValues(t *testing.T) {
	Convey("Getting values from os.Environ", t, func() {
		os.Clearenv()
		values := Values()
		So(len(values), ShouldBeZeroValue)

		Set("app", "me")
		values = Values()
		So(len(values), ShouldEqual, 1)
		So(values[Prefix+"_APP"], ShouldEqual, "me")
	})
}
