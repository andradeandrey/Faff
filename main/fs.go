// Copyright 2011 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/line"
	//"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

// PostFS takes care of reading posts from a file-system GIT directory
type PostFS struct {
	postDir   string
	shouldCfg *ShouldReread
	chro      []*postInfo
	urls      map[string]*postInfo
	lk        sync.Mutex
}

type postInfo struct {
	URLPath  string
	DiskPath string
}

var (
	ErrLineTooLong = os.NewError("line too long")
	ErrParse       = os.NewError("parse")
	ErrNoPost      = os.NewError("no post")
)

// NewPostFS() creates a new PostFS{} object. You must call
// Refresh() before you can use any other methods.
func NewPostFS(postDir string) *PostFS {
	return &PostFS{
		postDir:   postDir,
		shouldCfg: NewShouldReread(path.Join(postDir, "POSTFS")),
		chro:      make([]*postInfo, 0),
		urls:      make(map[string]*postInfo),
	}
}

func (pfs *PostFS) Dir() string { return pfs.postDir }

// GetDiskPath() returns the absolute disk path of the 
// post whose url path is urlpath
func (pfs *PostFS) GetDiskPath(urlpath string) (abs, rel string) {
	pfs.Refresh()
	pfs.lk.Lock()
	defer pfs.lk.Unlock()

	pinfo := pfs.urls[urlpath]
	if pinfo == nil {
		return "", ""
	}
	return path.Join(pfs.postDir, pinfo.DiskPath), pinfo.DiskPath
}

// Latest() returns the URLs of the latest few posts
func (pfs *PostFS) Latest() []string {
	pfs.Refresh()
	pfs.lk.Lock()
	defer pfs.lk.Unlock()

	ll := pfs.chro[:min(20, len(pfs.chro))]
	r := make([]string, len(ll))
	for i, pinfo := range ll {
		r[i] = pinfo.URLPath
	}
	return r
}

// Refresh() reads the PostFS config file and refreshes the current in-memory post DB
func (pfs *PostFS) Refresh() os.Error {

	// Check if POSTFS has changed on disk
	should, err := pfs.shouldCfg.Should()
	if err != nil {
		log.Printf("Warn: POSTFS file seems to be missing (err=%s)\n", err)
		return nil
	}
	if !should {
		return nil
	}
	log.Printf("Refreshing POSTFS\n")

	pfs.lk.Lock()
	defer pfs.lk.Unlock()

	file, err := os.OpenFile(path.Join(pfs.postDir, "POSTFS"), os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	chro, err := parsePostFSConfig(file)
	if err != nil {
		log.Printf("Problem: Error parsing POSTFS\n")
		return err
	}
	pfs.urls = make(map[string]*postInfo)
	for _, pinfo := range chro {
		pfs.urls[pinfo.URLPath] = pinfo
	}
	pfs.chro = chro
	log.Printf("Indexed %d posts in POSTFS\n", len(chro))
	return err
}

var lineRegExp = regexp.MustCompile(`^([^\t ]+)([\t ]+)([^\t ]+)$`)

func parsePostFSConfig(r io.Reader) (chro []*postInfo, err os.Error) {
	chro = []*postInfo{}
	lr := line.NewReader(r, 2000)
	for {
		l, isprefix, err0 := lr.ReadLine()
		if err0 != nil {
			err = err0
			break
		}
		if isprefix {
			err = ErrLineTooLong
			break
		}
		line := strings.TrimSpace(string(l))
		if line == "" || line[0] == '#' {
			continue
		}
		mtch := lineRegExp.FindStringSubmatch(line)
		if mtch == nil {
			return nil, ErrParse
		}
		chro = append(chro, &postInfo{mtch[1], mtch[3]})
	}
	if err != os.EOF {
		return nil, err
	}
	return chro, nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
