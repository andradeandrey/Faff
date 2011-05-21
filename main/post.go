// Copyright 2011 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"encoding/line"
	//"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/textproto"
	"os"
	"path"
	"strings"
	"sync"
	"time"
	"github.com/petar/GoHTTP/http"
	"github.com/petar/GoHTTP/template"
)

// PostMan takes care of parsing posts from their respective files and rendering (and maybe caching)
// them
type PostMan struct {
	gitCmd     string
	templDir   string  // template directory
	fs         *PostFS // manager of URL-to-disk mappings
	page       Page
	postTempl  *template.CachedTemplate
	indexTempl *template.CachedTemplate
	lk         sync.Mutex
}

type Post struct {
	ID          string
	URL         string
	URLEscaped  string
	Title       string
	Body        string
	Exerpt      string
	Tags        []string
	Created     string
	Updated     string
	CreatedTime *time.Time
	UpdatedTime *time.Time
	Hidden      bool	// Hidden posts are accessible through their URL, but not shown in the index page
}

type PostPage struct {
	Page *Page
	Post *Post
}

const PostTimeLayout = "Jan _2, 2006"

// make and init postman
func NewPostMan(gitCmd, templDir, postDir, disqusDevMode string, config *SiteConfig) *PostMan {
	return &PostMan{
		gitCmd:     gitCmd,
		templDir:   templDir,
		fs:         NewPostFS(postDir),
		page:       Page{
			RootURL:           config.RootURL,
			Title:             config.Title,
			TitleHTML:         config.TitleHTML,
			SubTitle:          config.SubTitle,
			Author:            config.Author,
			AuthorTwitter:     config.AuthorTwitter,
			DisqusDevMode:     GetDisqusDevMode(disqusDevMode),
			DisqusShortname:   config.DisqusShortname,
			GoogleAnalyticsID: config.GoogleAnalyticsID,
			FacebookAdminID:   config.FacebookAdminID,
		},
		postTempl:  template.NewCachedTemplate(path.Join(templDir, "post.templ"), nil),
		indexTempl: template.NewCachedTemplate(path.Join(templDir, "index.templ"), nil),
	}
}

func (pm *PostMan) RefreshFS() os.Error { return pm.fs.Refresh() }

// Latest() returns an array of the latest posts together with their rendered bodies
func (pm *PostMan) Latest() []*Post {
	purls := pm.fs.Latest()
	r := []*Post{}
	for _, purl := range purls {
		post, err := pm.Parse(purl)
		if err != nil {
			log.Printf("Problem: Parsing for permalink '%s' (%s)\n", purl, err)
			continue
		}
		if !post.Hidden {
			r = append(r, post)
		}
	}
	return r
}

func (pm *PostMan) ParseAndRender(urlPath string) (io.ReadCloser, os.Error) {
	post, err := pm.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	return pm.Render(post)
}

// Parse() returns the Post{} object of the post for urlPath
func (pm *PostMan) Parse(urlPath string) (*Post, os.Error) {

	// If urlPath has a trailing '/', remove it
	if len(urlPath) > 0 && urlPath[len(urlPath)-1] == '/' {
		urlPath = urlPath[:len(urlPath)-1]
	}

	// Obtain the post filename
	postfileAbs, postfileRel := pm.fs.GetDiskPath(urlPath)
	if postfileAbs == "" {
		return nil, ErrNoPost
	}

	// Read GIT info
	ct, ut, err := GITGetCreateUpdateTime(pm.gitCmd, pm.fs.Dir(), postfileRel)
	if err != nil {
		log.Printf("Problem: Reading GIT times from '%s' (%s)\n", postfileRel, err)
		return nil, err
	}

	// Read post contents
	rawbuf, err := ioutil.ReadFile(postfileAbs)
	if err != nil {
		return nil, err
	}
	bufr := bufio.NewReader(bytes.NewBuffer(rawbuf))

	// Read header
	mimer := textproto.NewReader(bufr)
	h, err := mimer.ReadMIMEHeader()
	if err != nil {
		log.Printf("Problem: Reading MIME header of '%s': (%s)\n", postfileRel, err)
		return nil, err
	}

	// Parse tags
	tags := []string{}
	for _, l := range h["Tags"] {
		ll := strings.Split(l, ",", -1)
		for _, t := range ll {
			tt := strings.TrimSpace(t)
			if tt == "" {
				continue
			}
			tags = append(tags, tt)
		}
	}

	// Parse options
	var hidden bool
	for _, l := range h["Options"] {
		ll := strings.Split(l, ",", -1)
		for _, t := range ll {
			switch strings.ToLower(strings.TrimSpace(t)) {
			case "hidden":
				hidden = true
			}
		}
	}

	// Read title and body
	liner := line.NewReader(bufr, 2000)
	var title string
	for {
		title0, _, err := liner.ReadLine()
		if err != nil {
			log.Printf("Problem: Reading title of '%s': (%s)\n", postfileRel, err)
			return nil, err
		}
		title1 := make([]byte, len(title0))
		copy(title1, title0)
		title = strings.TrimSpace(string(title1))
		if title != "" {
			break
		}
	}

	// Read body of post and exerpt
	exerpt := ""
	body := []byte{}
	for {
		line, _, err := liner.ReadLine()
		if err != nil {
			if err == os.EOF {
				break
			}
			return nil, err
		}
		l := strings.TrimSpace(string(line))
		if l == "<!--more-->" {
			exerpt = string(body)
		}
		body = append(body, line...)
		body = append(body, '\n')
	}

	// Prepare post object
	post := &Post{
		ID:          h.Get("Id"),
		Title:       title,
		URL:         urlPath,
		URLEscaped:  http.URLEscape(urlPath),
		Body:        string(body),
		Exerpt:      exerpt,
		Tags:        tags,
		Created:     ct.Format(PostTimeLayout),
		CreatedTime: ct,
		UpdatedTime: ut,
		Hidden:      hidden,
	}
	if ut != nil {
		post.Updated = ut.Format(PostTimeLayout)
	}

	return post, nil
}

// Render() renders the post page and returns the corresponding io.ReadCloser
func (pm *PostMan) Render(post *Post) (body io.ReadCloser, err os.Error) {
	pm.lk.Lock()
	defer pm.lk.Unlock()

	var w bytes.Buffer
	templ, err := pm.postTempl.Get()
	if err != nil {
		log.Printf("Problem: Fetching post template (%s)\n", err)
		return nil, err
	}
	postpage := &PostPage{
		Page: &pm.page,
		Post: post,
	}
	if err = templ.Execute(&w, postpage); err != nil {
		log.Printf("Problem: Applying post template (%s)\n", err)
		return nil, err
	}
	return http.NopCloser{&w}, nil
}
