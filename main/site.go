// Copyright 2011 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"json"
	"os"
)

// SiteConfig contains global blog and administrator constants
type SiteConfig struct {
	RootURL           string // No slashes before or after the URL
	Title             string
	TitleHTML         string
	SubTitle          string
	Author            string
	AuthorTwitter     string
	DisqusShortname   string
	GoogleAnalyticsID string
	FacebookAdminID   string
}

func ParseSiteConfig(filename string) (*SiteConfig, os.Error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	c := &SiteConfig{}
	err = json.Unmarshal(buf, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
