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
	// Home page
	goji.Get("/", application.Route(controller, "Index"))
	goji.Get("/index.html", application.Route(controller, "Index"))

	goji.Get("/sitemap.xml", application.Route(controller, "Sitemap"))

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

	// Category Index
	goji.Get(regexp.MustCompile(`^/category/(?P<cat>[a-z]+).html$`), application.Route(controller, "Category"))
	// Respotory
	goji.Get(regexp.MustCompile(`^/(?P<repo>[a-z0-9-]+).html$`), application.Route(controller, "Repository"))

	graceful.PostHook(func() {
		application.Close()
	})
	goji.Serve()
}