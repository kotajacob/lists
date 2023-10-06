# lists
# See LICENSE for copyright and license details.
.POSIX:

PREFIX ?= /usr
GO ?= go
GOFLAGS ?= -buildvcs=false
RM ?= rm -f

all: lists

lists:
	$(GO) build $(GOFLAGS) .

install: all
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f lists $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/lists

uninstall:
	$(RM) $(DESTDIR)$(PREFIX)/bin/lists

clean:
	$(RM) lists

run:
	go run .

watch:
	fd -e go -e tmpl | entr -rs "go run ."

.PHONY: all lists install uninstall clean run watch
