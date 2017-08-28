#!/usr/bin/env bash

function clean_vendor() {
}

function clean_gopath() {
	rm -rf github.com/ikeikeikeike/go-sitemap-generator && (rmdir github.com/ikeikeikeike > /dev/null 2>&1 || true)
	rm -rf gopkg.in/yaml.v2 && (rmdir gopkg.in > /dev/null 2>&1 || true)
	rm -rf github.com/andybalholm/cascadia && (rmdir github.com/andybalholm > /dev/null 2>&1 || true)
	rm -rf github.com/gorilla/securecookie && (rmdir github.com/gorilla > /dev/null 2>&1 || true)
	rm -rf golang.org/x/crypto && (rmdir golang.org/x > /dev/null 2>&1 || true)
	rm -rf github.com/thoas/stats && (rmdir github.com/thoas > /dev/null 2>&1 || true)
	rm -rf github.com/PuerkitoBio/goquery && (rmdir github.com/PuerkitoBio > /dev/null 2>&1 || true)
	rm -rf github.com/codegangsta/cli && (rmdir github.com/codegangsta > /dev/null 2>&1 || true)
	rm -rf github.com/gorilla/mux && (rmdir github.com/gorilla > /dev/null 2>&1 || true)
	rm -rf github.com/cloudflare/cloudflare-go && (rmdir github.com/cloudflare > /dev/null 2>&1 || true)
	rm -rf github.com/Sirupsen/logrus && (rmdir github.com/Sirupsen > /dev/null 2>&1 || true)
	rm -rf github.com/gorilla/context && (rmdir github.com/gorilla > /dev/null 2>&1 || true)
	rm -rf xi2.org/x/xz && (rmdir xi2.org/x > /dev/null 2>&1 || true)
	rm -rf gopkg.in/check.v1 && (rmdir gopkg.in > /dev/null 2>&1 || true)
	rm -rf github.com/go-utils/uslice && (rmdir github.com/go-utils > /dev/null 2>&1 || true)
	rm -rf github.com/blevesearch/bleve && (rmdir github.com/blevesearch > /dev/null 2>&1 || true)
	rm -rf github.com/dustin/go-humanize && (rmdir github.com/dustin > /dev/null 2>&1 || true)
	rm -rf github.com/cloudflare/cfssl && (rmdir github.com/cloudflare > /dev/null 2>&1 || true)
	rm -rf github.com/jinzhu/inflection && (rmdir github.com/jinzhu > /dev/null 2>&1 || true)
	rm -rf github.com/julienschmidt/httprouter && (rmdir github.com/julienschmidt > /dev/null 2>&1 || true)
	rm -rf gopkg.in/mgo.v2 && (rmdir gopkg.in > /dev/null 2>&1 || true)
	rm -rf github.com/blevesearch/segment && (rmdir github.com/blevesearch > /dev/null 2>&1 || true)
	rm -rf github.com/jinzhu/gorm && (rmdir github.com/jinzhu > /dev/null 2>&1 || true)
	rm -rf github.com/beevik/etree && (rmdir github.com/beevik > /dev/null 2>&1 || true)
	rm -rf github.com/Redundancy/go-sync && (rmdir github.com/Redundancy > /dev/null 2>&1 || true)
	rm -rf golang.org/x/oauth2 && (rmdir golang.org/x > /dev/null 2>&1 || true)
	rm -rf github.com/cbroglie/mustache && (rmdir github.com/cbroglie > /dev/null 2>&1 || true)
	rm -rf github.com/fatih/structs && (rmdir github.com/fatih > /dev/null 2>&1 || true)
	rm -rf github.com/mattn/go-sqlite3 && (rmdir github.com/mattn > /dev/null 2>&1 || true)
	rm -rf github.com/pkg/profile && (rmdir github.com/pkg > /dev/null 2>&1 || true)
	rm -rf github.com/davecgh/go-spew && (rmdir github.com/davecgh > /dev/null 2>&1 || true)
	rm -rf github.com/boltdb/bolt && (rmdir github.com/boltdb > /dev/null 2>&1 || true)
	rm -rf gopkg.in/vmihailenco/msgpack.v2 && (rmdir gopkg.in/vmihailenco > /dev/null 2>&1 || true)
	rm -rf github.com/blevesearch/go-porterstemmer && (rmdir github.com/blevesearch > /dev/null 2>&1 || true)
	rm -rf golang.org/x/sys/unix && (rmdir golang.org/x/sys > /dev/null 2>&1 || true)
	rm -rf github.com/gorilla/sessions && (rmdir github.com/gorilla > /dev/null 2>&1 || true)
	rm -rf github.com/golang/protobuf && (rmdir github.com/golang > /dev/null 2>&1 || true)
	rm -rf github.com/zenazn/goji && (rmdir github.com/zenazn > /dev/null 2>&1 || true)
	rm -rf github.com/google/go-github && (rmdir github.com/google > /dev/null 2>&1 || true)
	rm -rf github.com/google/go-querystring && (rmdir github.com/google > /dev/null 2>&1 || true)
	rm -rf github.com/mailgun/timetools && (rmdir github.com/mailgun > /dev/null 2>&1 || true)
	rm -rf github.com/pkg/errors && (rmdir github.com/pkg > /dev/null 2>&1 || true)
	rm -rf github.com/imdario/mergo && (rmdir github.com/imdario > /dev/null 2>&1 || true)
	rm -rf github.com/steveyen/gtreap && (rmdir github.com/steveyen > /dev/null 2>&1 || true)
}
