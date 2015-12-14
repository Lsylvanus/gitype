// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"runtime"
	"testing"

	"github.com/issue9/assert"
)

func TestCheckConfigDir(t *testing.T) {
	a := assert.New(t)

	a.Error(checkConfigDir("", "uploadDir"))
	a.Error(checkConfigDir("/abc", "uploadDir"))
	a.NotError(checkConfigDir("/abc/", "uploadDir"))

	if runtime.GOOS == "windows" {
		a.NotError(checkConfigDir("/abc\\", "uploadDir"))
	}
}

func TestCheckConfigURL(t *testing.T) {
	a := assert.New(t)

	a.Error(checkConfigURL("", "uploadURL"))
	a.Error(checkConfigURL("/abc/", "uploadURL"))
	a.NotError(checkConfigURL("/abc", "uploadURL"))
}

func TestLoadConfig(t *testing.T) {
	a := assert.New(t)

	cfg, err := loadConfig("./testdata/app.json")
	a.NotError(err).NotNil(cfg)

	a.Equal(cfg.FrontAPIPrefix, "/api").
		Equal(cfg.DBDriver, "sqlite3")
}
