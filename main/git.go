// Copyright 2011 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/line"
	"exec"
	//"fmt"
	"path"
	"os"
	"time"
)

const GITTimeLayout = "Mon, 2 Jan 2006 15:04:05 -0700"

func GITGetCreateUpdateTime(gitcmd, repo, filename string) (create, update *time.Time, err os.Error) {
	aa, err := GITGetAddTimes(gitcmd, repo, filename)
	if err != nil {
		return nil, nil, err
	}
	if len(aa) == 0 {
		return nil, nil, os.ENOENT
	}
	create = aa[len(aa)-1]
	if len(aa) > 1 {
		update = aa[0]
	}
	return
}

func GITGetAddTimes(gitcmd, repo, filename string) ([]*time.Time, os.Error) {
	ss, err := Run(gitcmd, "",
		[]string{" ", "--git-dir=" + path.Join(repo, ".git"),
			"log", "--format=%aD", "--", delPrefixSlash(filename)})
	if err != nil {
		return nil, err
	}
	tt := make([]*time.Time, len(ss))
	for i, s := range ss {
		t, err := time.Parse(GITTimeLayout, s)
		if err != nil {
			return nil, err
		}
		tt[i] = t
	}
	return tt, nil
}

func Run(prog, dir string, argv []string) ([]string, os.Error) {
	cmd, err := exec.Run(prog, argv, nil, dir, exec.Pipe, exec.Pipe, exec.Pipe)
	if err != nil {
		return []string{""}, err
	}
	liner := line.NewReader(cmd.Stdout, 2000)
	r := []string{}
	var l []byte
	for {
		l, _, err = liner.ReadLine()
		if err != nil {
			break
		}
		r = append(r, string(l))
	}
	if err != nil && err != os.EOF {
		return nil, err
	}
	err = cmd.Close()
	return r, err
}

func delPrefixSlash(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == '/' {
		return s[1:]
	}
	return s
}
