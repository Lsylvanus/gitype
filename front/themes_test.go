// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package front

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLoadThemeFile(t *testing.T) {
	a := assert.New(t)

	theme, err := loadThemeFile("./testdata/theme1/theme.json")
	a.NotError(err)
	a.Equal(theme.Name, "default").Equal(theme.Author.Name, "caixw")
}