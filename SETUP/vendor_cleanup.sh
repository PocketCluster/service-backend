#!/usr/bin/env bash

function clean_vendor() {
}

function clean_gopath() {
	rm -rf github.com/jinzhu/gorm && (rmdir github.com/jinzhu > /dev/null 2>&1 || true)
}
