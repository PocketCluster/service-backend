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

    "github.com/stkim1/service-backend/api"
    "github.com/stkim1/service-backend/api/auth/model"
    "github.com/stkim1/service-backend/shared/errmsg"
)

const (
    valid_inviation = "14b43cd77d391e05b1f24f5237aa596f63cf1bf5"
    valid_device    = "6c458e9b8821e4b5b6053bf91dc46723ad0e42d3"
)

func testOrmFilename() string {
    var (
        _, testfile, _, _ = runtime.Caller(0)
    )
    return filepath.Join(filepath.Dir(testfile), "auth.sql")
}

func openTestOrm() (*gorm.DB, error) {
    orm, err := gorm.Open("sqlite3", testOrmFilename())
    if err != nil {
        return nil, err
    }
    orm.CreateTable(&model.AuthIdentity{})
    return orm, nil
}

func closeTestOrm(orm *gorm.DB) error {
    if orm == nil {
        return nil
    }
    orm.Close()
    return os.Remove(testOrmFilename())
}

func openRouteWithAuth() (*httprouter.Router, *gorm.DB, error) {
    orm, err := openTestOrm()
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

func Test_UncoveredCountry(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }
    defer closeTestOrm(orm)

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

func Test_Empty_Post_Value(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }
    defer closeTestOrm(orm)

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

func Test_Invalid_Inviatation(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }
    defer closeTestOrm(orm)

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

func Test_No_DeviceHash(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }
    defer closeTestOrm(orm)

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

func Test_Invitation_With_Invalid_DevicePair(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }
    defer closeTestOrm(orm)

    var (
        a = model.AuthIdentity{
            InvHash: valid_inviation,
            DevHash: valid_device,
        }
    )
    orm.Create(&a)

    v := url.Values{}
    v.Set("invitation", valid_inviation)
    v.Set("device", "1c458e9b8821e4b5b6053bf91dc46723ad0e42d3")

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

func Test_Invitation_NotFound(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }
    defer closeTestOrm(orm)

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

func Test_Invitation_Without_DevicePair(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }
    defer closeTestOrm(orm)

    var (
        a = model.AuthIdentity{
            InvHash: valid_inviation,
        }
        b = model.AuthIdentity{}
    )
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
    if w.Body.String() != fmt.Sprintln(`{"auth":"pass","error":""}`) {
        t.Errorf("[%v] invalid response body : %v", api.URLAuthCheck, w.Body.String())
        t.FailNow()
    }

    orm.Where(fmt.Sprintf("%s = ?", model.ColInvHash), valid_inviation).First(&b)
    if len(b.InvHash) == 0 {
        t.Errorf("unable to find inviation with valid code")
        t.FailNow()
    }
    if b.InvHash != valid_inviation {
        t.Errorf("unable to find the corresponding invitation with code")
        t.FailNow()
    }
    if b.DevHash != valid_device {
        t.Errorf("incorrect device hash for invitation")
        t.FailNow()
    }
}

func Test_Invitation_With_DevicePair(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }
    defer closeTestOrm(orm)

    var (
        a = model.AuthIdentity{
            InvHash: valid_inviation,
            DevHash: valid_device,
        }
        b = model.AuthIdentity{}
    )
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
    if w.Body.String() != fmt.Sprintln(`{"auth":"pass","error":""}`) {
        t.Errorf("[%v] invalid response body : %v", api.URLAuthCheck, w.Body.String())
        t.FailNow()
    }

    orm.Where(fmt.Sprintf("%s = ?", model.ColInvHash), valid_inviation).First(&b)
    if len(b.InvHash) == 0 {
        t.Errorf("unable to find inviation with valid code")
        t.FailNow()
    }
    if b.InvHash != valid_inviation {
        t.Errorf("unable to find the corresponding invitation with code")
        t.FailNow()
    }
    if b.DevHash != valid_device {
        t.Errorf("incorrect device hash for invitation")
        t.FailNow()
    }
}


func Test_Invitation_Check_With_Pool(t *testing.T) {
    router, orm, err := openRouteWithAuth()
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }
    defer closeTestOrm(orm)

    var (
        a = model.AuthIdentity{
            InvHash: valid_inviation,
            DevHash: valid_device,
        }
        a1 = model.AuthIdentity{
            InvHash: "24b43cd77d391e05b1f24f5237aa596f63cf1bf5",
            DevHash: valid_device,
        }
        b = model.AuthIdentity{}
    )
    orm.Create(&a).Create(&a1)

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
    if w.Body.String() != fmt.Sprintln(`{"auth":"pass","error":""}`) {
        t.Errorf("[%v] invalid response body : %v", api.URLAuthCheck, w.Body.String())
        t.FailNow()
    }

    orm.Where(fmt.Sprintf("%s = ?", model.ColInvHash), valid_inviation).First(&b)
    if len(b.InvHash) == 0 {
        t.Errorf("unable to find inviation with valid code")
        t.FailNow()
    }
    if b.InvHash != valid_inviation {
        t.Errorf("unable to find the corresponding invitation with code")
        t.FailNow()
    }
    if b.DevHash != valid_device {
        t.Errorf("incorrect device hash for invitation")
        t.FailNow()
    }
}