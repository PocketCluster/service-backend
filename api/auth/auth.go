package auth

import (
    "net/http"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/julienschmidt/httprouter"
    "github.com/jinzhu/gorm"
    //"golang.org/x/crypto/ripemd160"

    "github.com/stkim1/api/abnormal"
    "github.com/stkim1/api/auth/model"
    "github.com/stkim1/sharedpkg/errmsg"
    "github.com/stkim1/sharedpkg/cforigin"
)

type AuthGateway interface {
    IsUserAuthValid(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}

type authGateway struct {
    orm     *gorm.DB
}

func NewAuthGateway(orm *gorm.DB) (AuthGateway, error) {
    var (
        authid = &model.AuthIdentity{}
    )
    if !orm.HasTable(authid) {
        orm.CreateTable(authid)
    } else {
        orm.AutoMigrate(authid)
    }
    return &authGateway{
        orm:    orm,
    }, nil
}

func (a *authGateway) IsUserAuthValid(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

    if err := cforigin.IsOriginAllowedCountry(r); err != nil {
        log.Debugf(errors.WithStack(err).Error())
        abnormal.ResponseJsonError(w, errmsg.ErrMsgJsonUnallowedCountry, http.StatusForbidden)
        return
    }

    iHash := r.FormValue("invitation")
    dHash := r.FormValue("device")
    log.Debugf("invitation %v | device %v", iHash, dHash)

    w.Header().Set("Server", "PocketCluster API Service")
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Connection", "keep-alive")
    //w.Header().Set("Etag", d.Name()) // file name
    w.Header().Set("Cache-Control", "max-age=3600") // 1 hr
    w.Header().Set("Expires", time.Now().UTC().Add(time.Hour).Format("Mon, 2 Jan 2006 15:04:05 MST"))
}