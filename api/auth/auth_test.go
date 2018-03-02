package auth

import (
    "net/http"
    "net/http/httptest"
    "os"
    "path/filepath"
    "runtime"
    "testing"

    _ "github.com/mattn/go-sqlite3"
    "github.com/jinzhu/gorm"
    "github.com/julienschmidt/httprouter"

    "github.com/stkim1/api"
)

func openRouteWithAuth() (*httprouter.Router, *gorm.DB, error) {
    var (
        _, testfile, _, _ = runtime.Caller(0)
        dbfile = filepath.Join(filepath.Dir(testfile), "auth.sql")
    )
    orm, err := gorm.Open("sqlite3", dbfile)
    if err != nil {
        return nil, nil, err
    }
    authsrvc, err := NewAuthGateway(orm)
    if err != nil {
        return nil, nil, err
    }
    router := httprouter.New()
    router.POST(api.URLAuthCheck,  authsrvc.IsUserAuthValid)
    return router, orm, nil
}

func closeTestOrm(orm *gorm.DB) error {
    if orm == nil {
        return nil
    }
    orm.Close()
    var (
        _, testfile, _, _ = runtime.Caller(0)
        dbfile = filepath.Join(filepath.Dir(testfile), "auth.sql")
    )
    return os.Remove(dbfile)
}

func TestEmptyValue(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    defer closeTestOrm(orm)
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }

    r, _ := http.NewRequest("POST", api.URLAuthCheck, nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, r)
    if w.Code != http.StatusOK {
        t.Errorf("failed to serve " + api.URLAuthCheck)
        t.FailNow()
    }
}
