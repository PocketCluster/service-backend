package main

import (
    "os"
    "regexp"
    "runtime"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/gorilla/context"
    "github.com/zenazn/goji"
    "github.com/zenazn/goji/graceful"

    "github.com/stkim1/backend/framework"
    "github.com/stkim1/backend/control"
    "github.com/stkim1/backend/config"
    "github.com/stkim1/backend/util"
)

func main() {
    var (
        app *framework.Application
        ctrl *control.Controller
    )

    cfgPath, ok := os.LookupEnv(config.EnvConfigFilePath)
    // develop mode
    if ok {
        log.SetLevel(log.DebugLevel)

    // release mode
    } else {
        util.InitSysLogger()
        cfgPath = "config.yaml"
    }

    cfg, err := config.NewConfig(cfgPath)
    if err != nil {
        log.Panic(errors.WithMessage(err, "Cannot load config"))
        return
    }
    runtime.GOMAXPROCS(cfg.General.MaxConcurrency)
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

    // TODO : Handle Notfound
    //goji.NotFound()

    // TODO : what submix for?
    //v1Mux :=goji.SubMux()

    // dashboard
    goji.Get("/pocketcluster/dashboard/:mode",                                    app.AddRoute(ctrl.DashboardFront))
    goji.Post("/pocketcluster/dashboard/repository/:mode",                        app.AddRoute(ctrl.DashboardRepository))

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
    goji.Get("/sitemap.xml",                                                      app.AddRoute(ctrl.Sitemap))

    // Home page
    // FIXME: all three regexp fail. WTF? (https://github.com/zenazn/goji/issues/75) & (https://github.com/zenazn/goji/blob/master/web/regexp_pattern.go#L56)
    //goji.Get(regexp.MustCompile(`^/index.html\?page=(?P<page>\d+)$`), application.Route(controller, "IndexPaged"))
    //goji.Get(regexp.MustCompile(`^/index.html\?page=(?P<page>[0-9]+)$`), application.Route(controller, "IndexPaged"))
    //goji.Get(regexp.MustCompile(`^/index.html[?]page=(?P<page>[0-9]+)$`), application.Route(controller, "IndexPaged"))
    goji.Get(regexp.MustCompile(`^/index(?P<page>[0-9]+).html$`),                    app.AddRoute(ctrl.IndexPaged))
    goji.Get("/index.html",                                                       app.AddRoute(ctrl.Index))
    goji.Get("/",                                                                 app.AddRoute(ctrl.Index))

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