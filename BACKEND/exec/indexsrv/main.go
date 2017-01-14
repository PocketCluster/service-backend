package main

import (
    "os"
    "regexp"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/gorilla/context"
    "github.com/zenazn/goji"
    "github.com/zenazn/goji/graceful"

    "github.com/stkim1/BACKEND/framework"
    "github.com/stkim1/BACKEND/control"
    "github.com/stkim1/BACKEND/config"
)

func main() {
    var (
        app *framework.Application
        ctrl *control.Controller
    )
    cfgPath, ok := os.LookupEnv(config.EnvConfigFilePath)
    if !ok {
        cfgPath = "config.yaml"
    }
    cfg, err := config.NewConfig(cfgPath)
    if err != nil {
        log.Panic(trace.Wrap(err, "Cannot load config"))
        return
    }

    // Setup Controller
    ctrl = control.NewController(cfg)
    // setup Application
    app = framework.NewApplication(cfg, ctrl)

    // Apply middleware
    //goji.Use(app.ApplySessions)
    //goji.Use(app.ApplyCsrfProtection)
    //goji.Use(app.ApplyAuth)
    goji.Use(app.ApplyDbMap)
    goji.Use(app.ApplyIsXhr)
    goji.Use(context.ClearHandler)

    // dashboard
    goji.Get("/pocketcluster/dashboard/:mode",                                       app.AddRoute(ctrl.DashboardFront))
    goji.Post("/pocketcluster/dashboard/repository/:mode",                           app.AddRoute(ctrl.DashboardRepository))

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
    goji.Get("/sitemap.xml",                                                         app.AddRoute(ctrl.Sitemap))

    // Home page
    // FIXME: all three regexp fail. WTF? (https://github.com/zenazn/goji/issues/75) & (https://github.com/zenazn/goji/blob/master/web/regexp_pattern.go#L56)
    //goji.Get(regexp.MustCompile(`^/index.html\?page=(?P<page>\d+)$`), application.Route(controller, "IndexPaged"))
    //goji.Get(regexp.MustCompile(`^/index.html\?page=(?P<page>[0-9]+)$`), application.Route(controller, "IndexPaged"))
    //goji.Get(regexp.MustCompile(`^/index.html[?]page=(?P<page>[0-9]+)$`), application.Route(controller, "IndexPaged"))
    goji.Get(regexp.MustCompile(`^/index(?P<page>[0-9]+).html$`),                    app.AddRoute(ctrl.IndexPaged))
    goji.Get("/index.html",                                                          app.AddRoute(ctrl.Index))
    goji.Get("/",                                                                    app.AddRoute(ctrl.Index))

    // Category Index
    goji.Get(regexp.MustCompile(`^/category/(?P<cat>[a-z]+)(?P<page>[0-9]+).html$`), app.AddRoute(ctrl.CategoryPaged))
    goji.Get(regexp.MustCompile(`^/category/(?P<cat>[a-z]+).html$`),                 app.AddRoute(ctrl.Category))

    // Respotory
    goji.Get(regexp.MustCompile(`^/(?P<repo>[a-z0-9-]+).html$`),                     app.AddRoute(ctrl.Repository))

    graceful.PostHook(func() {
        app.Close()
    })

    // just before going into serve, initiate updater
    app.ScheduleMetaUpdate()
    app.ScheduleSuppUpdate()
    goji.Serve()
}