// Copyright 2011 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"sync"
)

// ShouldReread keeps track of whether a file has been updated since last reading it.
type ShouldReread struct {
	fname string
	mtime int64
	lk    sync.Mutex
}

func NewShouldReread(filename string) *ShouldReread {
	return &ShouldReread{fname: filename}
}

func (sh *ShouldReread) Should() (shoud bool, err os.Error) {
	fi, err := os.Stat(sh.fname)
	if err != nil {
		return false, err
	}

	sh.lk.Lock()
	defer sh.lk.Unlock()

	if fi.Mtime_ns > sh.mtime {
		sh.mtime = fi.Mtime_ns
		return true, nil
	}
	return false, nil
}

// XXX: If to be used, this object could benefit from locks
type CachedObject struct {
	fname string
	fetch func(string) (obj interface{}, err os.Error)
	obj   interface{}
	mtime int64
}

func NewCachedObject(filename string, fetch func(string) (interface{}, os.Error)) *CachedObject {
	return &CachedObject{fname: filename, fetch: fetch}
}

func (c *CachedObject) Get() (obj interface{}, err os.Error) {
	if c.obj == nil {
		return c.readFile()
	}
	fi, err := os.Stat(c.fname)
	if err != nil {
		return nil, err
	}
	if fi.Mtime_ns > c.mtime {
		return c.readFile()
	}
	return c.obj, nil
}

func (c *CachedObject) readFile() (obj interface{}, err os.Error) {
	fi, err := os.Stat(c.fname)
	if err != nil {
		return nil, err
	}
	obj, err = c.fetch(c.fname)
	if err != nil {
		return nil, err
	}
	c.obj = obj
	c.mtime = fi.Mtime_ns
	return obj, nil
}
