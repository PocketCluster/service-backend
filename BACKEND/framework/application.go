package framework

import (
    "html/template"
    "io"
    "net/http"
    "reflect"
    "crypto/sha256"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/gorilla/sessions"
    "github.com/zenazn/goji/web"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"

    "github.com/stkim1/BACKEND/model"
    "github.com/stkim1/BACKEND/control"
    "github.com/stkim1/BACKEND/config"
)

func NewApplication(config *config.Config, control *control.Controller) *Application {
    app := &Application{
        Config:         config,
        Controller:     control,
    }
    app.init()
    return app
}

type csrfProtection struct {
    Key    string
    Cookie string
    Header string
    Secure bool
}

// Application-wide resource management
type Application struct {
    Controller     *control.Controller
    Config         *config.Config
    Template       *template.Template
    Store          *sessions.CookieStore
    GORM           *gorm.DB
    CsrfProtection *csrfProtection
}

func (a *Application) init() {
    hash := sha256.New()
    io.WriteString(hash, a.Config.Cookie.MacSecret)
    a.Store = sessions.NewCookieStore(hash.Sum(nil))
    a.Store.Options = &sessions.Options{
        Path:     "/",
        HttpOnly: true,
        Secure:   a.Config.Cookie.Secure,
    }

    db, err := gorm.Open(a.Config.Database.DatabaseType, a.Config.Database.DatabasePath)
    if err != nil {
        log.Error(trace.Wrap(err,"Failed to open database"))
    }
    // Migrate the schema
    db.AutoMigrate(&model.Author{}, &model.Repository{}, &model.RepoCommit{}, &model.RepoVersion{}, &model.RepoLanguage{}, &model.RepoContributor{});

    // set relation
    // db.Model(&model.Repository{}).Related(&model.RepoVersion{})
    // db.Model(&model.Repository{}).Related(&model.RepoCommit{})
    // db.Model(&model.Repository{}).Related(&model.RepoLanguage{})
    a.GORM = db;

    a.CsrfProtection = &csrfProtection{
        Key:       a.Config.CSRF.Key,
        Cookie:    a.Config.CSRF.Cookie,
        Header:    a.Config.CSRF.Header,
        Secure:    a.Config.Cookie.Secure,
    }
}

func (a *Application) Close() {
    log.Info("!!!Application terminating!!!")
}

func (a *Application) Route(controller interface{}, route string) interface{} {
    fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
        c.Env["Content-Type"] = "text/html"

        methodValue := reflect.ValueOf(controller).MethodByName(route)
        methodInterface := methodValue.Interface()
        method := methodInterface.(func(c web.C, r *http.Request) (string, int))

        body, code := method(c, r)

        if session, exists := c.Env["Session"]; exists {
            err := session.(*sessions.Session).Save(r, w)
            if err != nil {
                log.Error(trace.Wrap(err, "Can't save session"))
            }
        }

        switch code {
        case http.StatusOK:
            if _, exists := c.Env["Content-Type"]; exists {
                w.Header().Set("Content-Type", c.Env["Content-Type"].(string))
            }
            io.WriteString(w, body)
        case http.StatusNotFound:
            http.Error(w, http.StatusText(404), 404)
        case http.StatusBadRequest:
            // FIXME : replace "error" with err.Error()
            http.Error(w, "error", http.StatusBadRequest)
        case http.StatusSeeOther, http.StatusFound:
            http.Redirect(w, r, body, code)
        }
    }
    return fn
}

func (a *Application) AddRoute(method func(c web.C, r *http.Request) (string, int)) interface{} {
    fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
        c.Env["Content-Type"] = "text/html"

        body, code := method(c, r)

        if session, exists := c.Env["Session"]; exists {
            err := session.(*sessions.Session).Save(r, w)
            if err != nil {
                log.Error(trace.Wrap(err, "Can't save session"))
            }
        }

        switch code {
        case http.StatusOK:
            if _, exists := c.Env["Content-Type"]; exists {
                w.Header().Set("Content-Type", c.Env["Content-Type"].(string))
            }
            io.WriteString(w, body)
        case http.StatusNotFound:
            http.Error(w, http.StatusText(404), 404)
        case http.StatusBadRequest:
            // FIXME : replace "error" with err.Error()
            http.Error(w, "error", http.StatusBadRequest)
        case http.StatusSeeOther, http.StatusFound:
            http.Redirect(w, r, body, code)
        }
    }
    return fn
}