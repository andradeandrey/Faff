## Overview

Faff is a blogging platform implemented in Go, which retrieves posts
from a local GIT repository.

## Maturity

Faff has been the server behind [Population Algorithms](http://popalg.org) for
over a month now with no downtime. Some features include:

* Keepalive
* Pipelining
* Integrated DISQUS threads
* Integrated Tweet and Facebook like buttons
* Tagging

There are many features that will come as the need arises:

* Search by keywords or tags
* Pages
* etc.

## Installation

You need a working installation of [Go](http://golang.org). Pull this
repository

	git clone git://github.com/petar/Faff.git Faff-git

Install the [GoHTTP]() packages

	goinstall github.com/petar/GoHTTP/util
	goinstall github.com/petar/GoHTTP/http
	goinstall github.com/petar/GoHTTP/server
	goinstall github.com/petar/GoHTTP/server/subs
	goinstall github.com/petar/GoHTTP/template

Build the Faff sever

	cd Faff-git
	make && make install

Create a config file for your server. You can see an example in

	example/example.config

Run the Faff server

	faff -bind=:80 -static=Faff-git/static -config=your-site-config -dir=your-git-posts-directory

Note that `your-git-posts-directory` is a GIT working directory on your file system that contains
your posts. The next section explains how to set this up.

## The GIT posts directory

Coming soon ...

## About

Faff was written by [Petar Maymounkov](http://pdos.csail.mit.edu/~petar/). 

Follow me on [Twitter @maymounkov](http://www.twitter.com/maymounkov)!
