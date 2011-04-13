// Copyright 2011 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"os"
	"github.com/petar/GoHTTP/http"
)

type IndexPage struct {
	Page  *Page
	Index []*Post
}

func (pm *PostMan) RenderIndex() (body io.ReadCloser, err os.Error) {
	indexpage := &IndexPage{
		Page:  &pm.page,
		Index: pm.Latest(),
	}
	var w bytes.Buffer

	pm.lk.Lock()
	defer pm.lk.Unlock()

	templ, err := pm.indexTempl.Get()
	if err != nil {
		return nil, err
	}
	if err = templ.Execute(&w, indexpage); err != nil {
		return nil, err
	}
	return http.NopCloser{&w}, nil
}
