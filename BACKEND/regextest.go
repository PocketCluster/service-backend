package main

import (
	"fmt"
	"regexp"
)

func main() {
	// Compile the expression once, usually at init time.
	// Use raw strings to avoid having to quote the backslashes.
	//var validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`)
	var validID = regexp.MustCompile(`^/index.html\?page=(?P<page>[0-9]+)$`)
	//goji.Get(regexp.MustCompile(`^/index.html\?page=(?P<page>\d+)$`), application.Route(controller, "IndexPaged"))
	//goji.Get(regexp.MustCompile(`^/index.html\?page=(?P<page>[0-9]+)$`), application.Route(controller, "IndexPaged"))

	fmt.Println(validID.MatchString("/index.html?page=32"))
	fmt.Println(validID.MatchString("/index.html?pages=32"))
}

