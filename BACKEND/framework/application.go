package framework

import (
    "crypto/sha256"
    "html/template"
    "io"
    "net/http"
    "reflect"
    "sync"
    "sync/atomic"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/gorilla/sessions"
    "github.com/zenazn/goji/web"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"

    "github.com/stkim1/BACKEND/model"
    "github.com/stkim1/BACKEND/control"
    "github.com/stkim1/BACKEND/config"
    "github.com/stkim1/BACKEND/storage"
    "github.com/stkim1/BACKEND/storage/boltbk"
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
    Controller          *control.Controller
    Config              *config.Config
    Template            *template.Template
    Store               *sessions.CookieStore
    MetaDB              *gorm.DB
    SuppleDB            storage.Nosql
    CsrfProtection      *csrfProtection

    // waiter
    UpdateWait          sync.WaitGroup
    IsMetaUpdating      atomic.Value
    QuitMetaUpdate      chan bool
    IsSuppUpdating      atomic.Value
    QuitSuppUpdate      chan bool
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

    // (SQLITE) metadata
    metadb, err := gorm.Open(a.Config.Database.DatabaseType, a.Config.Database.DatabasePath)
    if err != nil {
        log.Fatal(trace.Wrap(err))
    }
    // Migrate the schema
    metadb.AutoMigrate(&model.Repository{}, &model.Author{}, &model.RepoContributor{});
    a.MetaDB = metadb;
    a.QuitMetaUpdate = make(chan bool)
    a.IsMetaUpdating.Store(false)

    // (BOLTDB) supplementary
    suppledb, err := boltbk.New(a.Config.Supplement.DatabasePath)
    if err != nil {
        log.Fatal(trace.Wrap(err))
    }
    a.SuppleDB = suppledb
    a.QuitSuppUpdate = make(chan bool)
    a.IsSuppUpdating.Store(false)

    a.CsrfProtection = &csrfProtection{
        Key:       a.Config.CSRF.Key,
        Cookie:    a.Config.CSRF.Cookie,
        Header:    a.Config.CSRF.Header,
        Secure:    a.Config.Cookie.Secure,
    }
}

func (a *Application) Close() {
    log.Info("Wait for graceful completion...")
    a.QuitSuppUpdate <- true
    a.QuitMetaUpdate <- true
    a.UpdateWait.Wait()
    close(a.QuitSuppUpdate)
    close(a.QuitMetaUpdate)

    if a.MetaDB != nil {
        a.MetaDB.Close()
    }

    if a.SuppleDB != nil {
        a.SuppleDB.Close()
    }

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
