// Copyright 2011 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

type Page struct {
	RootURL           string	// Root URL of blog without 'http://' part, e.g. 'popalg.org'
	Title             string	// Blog title
	TitleHTML         string	// HTML to prettify the title heading on top of pages
	SubTitle          string	// Blog subtitle
	Author            string	// Blog author name
	AuthorTwitter     string	// Twitter username of blog author
	DisqusDevMode     int		// 1 = Run DISQUS in dev mode, 0 = run in production mode
	DisqusShortname   string	// Administrator's DISQUS username (aka shortname)
	GoogleAnalyticsID string	// Google Analytics ID
	FacebookAdminID   string	// Facebook Admin ID 
}

func GetDisqusDevMode(mode string) int {
	if mode != "prod" {
		return 1
	}
	return 0
}
