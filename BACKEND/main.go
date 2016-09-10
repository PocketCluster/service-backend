package main

import (
	"flag"
	"regexp"

	"github.com/golang/glog"
	"github.com/gorilla/context"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	"github.com/stkim1/BACKEND/framework"
	"github.com/stkim1/BACKEND/control"
)

func main() {

	filename := flag.String("config", "config.toml", "Path to configuration file")
	flag.Parse()
	defer glog.Flush()

	var application = &framework.Application{}
	application.Init(filename)

	// Apply middleware
	//goji.Use(application.ApplySessions)
	goji.Use(application.ApplyDbMap)
	//goji.Use(application.ApplyAuth)
	goji.Use(application.ApplyIsXhr)
	//goji.Use(application.ApplyCsrfProtection)
	goji.Use(context.ClearHandler)

	// Setup Routers
	controller := &control.Controller{}

	// admin
	goji.Get("/pocketcluster/dashboard/overview", application.Route(controller, "DashboardOverview"))
	goji.Get("/pocketcluster/dashboard/repository", application.Route(controller, "DashboardRepository"))
/*
	// Sign In routes
	goji.Get("/signin", application.Route(controller, "SignIn"))
	goji.Post("/signin", application.Route(controller, "SignInPost"))

	// Sign Up routes
	goji.Get("/signup", application.Route(controller, "SignUp"))
	goji.Post("/signup", application.Route(controller, "SignUpPost"))

	// KTHXBYE
	goji.Get("/logout", application.Route(controller, "Logout"))
*/

	// sitemap
	goji.Get("/sitemap.xml", application.Route(controller, "Sitemap"))

	// Home page
	// FIXME: all three regexp fail. WTF? (https://github.com/zenazn/goji/issues/75) & (https://github.com/zenazn/goji/blob/master/web/regexp_pattern.go#L56)
	//goji.Get(regexp.MustCompile(`^/index.html\?page=(?P<page>\d+)$`), application.Route(controller, "IndexPaged"))
	//goji.Get(regexp.MustCompile(`^/index.html\?page=(?P<page>[0-9]+)$`), application.Route(controller, "IndexPaged"))
	//goji.Get(regexp.MustCompile(`^/index.html[?]page=(?P<page>[0-9]+)$`), application.Route(controller, "IndexPaged"))
	goji.Get(regexp.MustCompile(`^/index.html/(?P<page>[0-9]+)$`), application.Route(controller, "IndexPaged"))
	goji.Get("/", application.Route(controller, "Index"))
	goji.Get("/index.html", application.Route(controller, "Index"))

	// Category Index
	goji.Get(regexp.MustCompile(`^/category/(?P<cat>[a-z]+).html/(?P<page>[0-9]+)$`), application.Route(controller, "CategoryPaged"))
	goji.Get(regexp.MustCompile(`^/category/(?P<cat>[a-z]+).html$`), application.Route(controller, "Category"))

	// Respotory
	goji.Get(regexp.MustCompile(`^/(?P<repo>[a-z0-9-]+).html$`), application.Route(controller, "Repository"))

	graceful.PostHook(func() {
		application.Close()
	})
	goji.Serve()
}