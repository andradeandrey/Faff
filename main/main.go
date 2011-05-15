// Copyright 2011 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"exec"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"github.com/petar/GoHTTP/http"
	"github.com/petar/GoHTTP/server"
	"github.com/petar/GoHTTP/server/subs"
	"github.com/petar/GoGauge/pprof"
)

func envOrDefault(cmd, dfl string) string {
	p, err := exec.LookPath(cmd)
	if err != nil {
		return dfl
	}
	return p
}

var (
	flagStaticDir = flag.String("static", "", "Path to framework static directory")
	flagSiteDir   = flag.String("site", "", "Path to site-specific static directory")
	flagDir       = flag.String("dir", "", "Path to content GIT directory")
	flagBind      = flag.String("bind", "0.0.0.0:80", "Address to bind web server to")
	flagGIT       = flag.String("git", envOrDefault("git", "/usr/local/git/bin/git"), "GIT command name")
	flagMode      = flag.String("mode", "dev", "Choose between 'dev' and 'prod' mode")
	flagConfig    = flag.String("config", "", "Name of site config file")
)

func main() {
	pprof.InstallPprof()
	fmt.Fprintf(os.Stderr, "Faff — 2011 — by Petar Maymounkov, petar@5ttt.org\n")
	flag.Parse()
	MonitorMemProfile()

	config, err := ParseSiteConfig(*flagConfig)
	if err != nil {
		log.Printf("Problem reading config file: %s\n", err)
		os.Exit(1)
	}

	postman = NewPostMan(*flagGIT, *flagStaticDir, *flagDir, *flagMode, config)
	err = postman.RefreshFS()
	if err != nil {
		log.Printf("Problem starting postman: %s\n", err)
		os.Exit(1)
	}

	srv, err := server.NewServerEasy(*flagBind)
	if err != nil {
		log.Printf("Problem binding WWW server: %s\n", err)
		os.Exit(1)
	}
	srv.AddSub("/s/", subs.NewStaticSub(*flagStaticDir))	
	srv.AddSub("/t/", subs.NewStaticSub(*flagSiteDir))	
	fmt.Printf("Faff server running ...\n")
	for {
		q, err := srv.Read()
		if err != nil {
			log.Printf("Problem responding to %s: %s\n", q.Req.RawURL, err)
			srv.Shutdown()
			os.Exit(1)
		}
		var resp *http.Response
		url := q.Req.URL
		if url == nil {
			resp = http.NewResponse404String(q.Req, "could not parse URL")
		} else {
			resp, err = respond(q.Req, url)
			if err != nil {
				resp = http.NewResponse404String(q.Req, err.String())
			}
		}
		q.ContinueAndWrite(resp)
	}
}

var postman *PostMan

func respond(req *http.Request, url *http.URL) (resp *http.Response, err os.Error) {
	var body io.ReadCloser
	if url.Path == "/" {
		body, err = postman.RenderIndex()
	} else {
		body, err = postman.ParseAndRender(url.Path)
	}
	if err != nil && err != ErrNoPost {
		log.Printf("Render: %s\n", err)
	}
	if err != nil {
		return nil, err
	}
	return http.NewResponseWithBody(req, body), nil
}
