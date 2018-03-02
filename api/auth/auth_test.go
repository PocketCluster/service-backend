package auth

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "net/url"
    "os"
    "path/filepath"
    "runtime"
    "strings"
    "testing"

    _ "github.com/mattn/go-sqlite3"
    "github.com/jinzhu/gorm"
    "github.com/julienschmidt/httprouter"

    "github.com/stkim1/api"
    "github.com/stkim1/api/auth/model"
    "github.com/stkim1/sharedpkg/errmsg"
)

const (
    valid_inviation = "14b43cd77d391e05b1f24f5237aa596f63cf1bf5"
    valid_device    = "6c458e9b8821e4b5b6053bf91dc46723ad0e42d3"
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

func TestUncoveredCountry(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    defer closeTestOrm(orm)
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }

    r, _ := http.NewRequest("POST", api.URLAuthCheck, nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, r)
    if w.Code != http.StatusForbidden {
        t.Errorf("[%v] invalid response code : %v", api.URLAuthCheck, w.Code)
        t.FailNow()
    }
    if w.Body.String() != fmt.Sprintln(errmsg.ErrMsgJsonUnallowedCountry) {
        t.Errorf("[%v] invalid response body : %v", api.URLAuthCheck, w.Body.String())
        t.FailNow()
    }
}

func TestEmptyPostValue(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    defer closeTestOrm(orm)
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }

    r, _ := http.NewRequest("POST", api.URLAuthCheck, nil)
    r.Header.Set("cf-ipcountry", "US")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, r)
    if w.Code != http.StatusForbidden {
        t.Errorf("[%v] invalid response code : %v", api.URLAuthCheck, w.Code)
        t.FailNow()
    }
    if w.Body.String() != fmt.Sprintln(errmsg.ErrMsgJsonInvalidInvitation) {
        t.Errorf("[%v] invalid response body : %v", api.URLAuthCheck, w.Body.String())
        t.FailNow()
    }
}

func TestInvalidInviatation(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    defer closeTestOrm(orm)
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }

    v := url.Values{}
    v.Set("invitation", "14b43cd77d391e05b1f24f523 aa596f63cf1bf5")

    r, err := http.NewRequest("POST", api.URLAuthCheck, strings.NewReader(v.Encode()))
    if err != nil {
        t.Errorf("request construction error : %v", err.Error())
        t.FailNow()
    }
    r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Set("cf-ipcountry", "US")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, r)
    if w.Code != http.StatusForbidden {
        t.Errorf("[%v] invalid response code : %v", api.URLAuthCheck, w.Code)
        t.FailNow()
    }
    if w.Body.String() != fmt.Sprintln(errmsg.ErrMsgJsonInvalidInvitation) {
        t.Errorf("[%v] invalid response body : %v", api.URLAuthCheck, w.Body.String())
        t.FailNow()
    }
}

func TestNoDeviceHash(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    defer closeTestOrm(orm)
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }

    v := url.Values{}
    v.Set("invitation", valid_inviation)

    r, err := http.NewRequest("POST", api.URLAuthCheck, strings.NewReader(v.Encode()))
    if err != nil {
        t.Errorf("request construction error : %v", err.Error())
        t.FailNow()
    }
    r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Set("cf-ipcountry", "US")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, r)
    if w.Code != http.StatusForbidden {
        t.Errorf("[%v] invalid response code : %v", api.URLAuthCheck, w.Code)
        t.FailNow()
    }
    if w.Body.String() != fmt.Sprintln(errmsg.ErrMsgJsonUnsubmittedDevice) {
        t.Errorf("[%v] invalid response body : %v", api.URLAuthCheck, w.Body.String())
        t.FailNow()
    }
}

func TestInvitationNotFound(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    defer closeTestOrm(orm)
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }

    v := url.Values{}
    v.Set("invitation", valid_inviation)
    v.Set("device", valid_device)

    r, err := http.NewRequest("POST", api.URLAuthCheck, strings.NewReader(v.Encode()))
    if err != nil {
        t.Errorf("request construction error : %v", err.Error())
        t.FailNow()
    }
    r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Set("cf-ipcountry", "US")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, r)
    if w.Code != http.StatusForbidden {
        t.Errorf("[%v] invalid response code : %v", api.URLAuthCheck, w.Code)
        t.FailNow()
    }
    if w.Body.String() != fmt.Sprintln(errmsg.ErrMsgJsonInvalidInvitation) {
        t.Errorf("[%v] invalid response body : %v", api.URLAuthCheck, w.Body.String())
        t.FailNow()
    }
}

func TestInvitationWithoutDevicePair(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    defer closeTestOrm(orm)
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }

    var (
        a, b = model.AuthIdentity{}, model.AuthIdentity{}
    )
    a.Invitation = valid_inviation
    orm.Create(&a)

    v := url.Values{}
    v.Set("invitation", valid_inviation)
    v.Set("device", valid_device)

    r, err := http.NewRequest("POST", api.URLAuthCheck, strings.NewReader(v.Encode()))
    if err != nil {
        t.Errorf("request construction error : %v", err.Error())
        t.FailNow()
    }
    r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Set("cf-ipcountry", "US")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, r)
    if w.Code != http.StatusOK {
        t.Errorf("[%v] invalid response. code : %v | message %v", api.URLAuthCheck, w.Code, w.Body.String())
        t.FailNow()
    }
    if w.Body.String() != fmt.Sprintln(`{"auth":"pass"}`) {
        t.Errorf("[%v] invalid response body : %v", api.URLAuthCheck, w.Body.String())
        t.FailNow()
    }

    orm.Where("invitation = ?", valid_inviation).First(&b)
    if len(b.Invitation) == 0 {
        t.Errorf("unable to find inviation with valid code")
        t.FailNow()
    }
    if b.Invitation != valid_inviation {
        t.Errorf("unable to find the corresponding invitation with code")
        t.FailNow()
    }
    if b.Device != valid_device {
        t.Errorf("incorrect device hash for invitation")
        t.FailNow()
    }
}