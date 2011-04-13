# Copyright 2011 Petar Maymounkov. All rights reserved.
# Use of this source code is governed by a GPL-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

all:	install

install:
	cd main && make install

clean:
	cd main && make clean

nuke:
	cd main && make nuke
