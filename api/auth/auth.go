package auth

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/julienschmidt/httprouter"
    "github.com/jinzhu/gorm"

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
    const (
        kInviation  string = "invitation"
        kDevice     string = "device"
    )
    var (
        authid = model.AuthIdentity{}
    )

    if err := cforigin.IsOriginAllowedCountry(r); err != nil {
        log.Debugf(errors.WithStack(err).Error())
        abnormal.ResponseJsonError(w, errmsg.ErrMsgJsonUnallowedCountry, http.StatusForbidden)
        return
    }

    // c37ace13-e333-4f51-bb38-eb5728d14a38 -> 14b43cd77d391e05b1f24f5235aa596f63cf1bf5 | ^[a-z0-9]{40}$
    iHash := strings.TrimSpace(r.FormValue(kInviation))
    if !model.IsValidHash(iHash) {
        abnormal.ResponseJsonError(w, errmsg.ErrMsgJsonInvalidInvitation, http.StatusForbidden)
        return
    }

    // G8815052XYL -> 6c458e9b8821e4b5b6053bf91dc46723ad0e42d3 | ^[a-z0-9]{40}$
    dHash := strings.TrimSpace(r.FormValue(kDevice))
    if !model.IsValidHash(dHash) {
        abnormal.ResponseJsonError(w, errmsg.ErrMsgJsonUnsubmittedDevice, http.StatusForbidden)
        return
    }

    // find invitation
    a.orm.Where(fmt.Sprintf("%s = ?", model.ColInvHash), iHash).First(&authid)
    if len(authid.InvHash) == 0 {
        abnormal.ResponseJsonError(w, errmsg.ErrMsgJsonInvalidInvitation, http.StatusForbidden)
        return
    }

    // check device hash
    if len(authid.DevHash) != 0 {

        // 1. when device hash is not equal to submitted hash
        if authid.DevHash != dHash {
            abnormal.ResponseJsonError(w, errmsg.ErrMsgJsonUnsubmittedDevice, http.StatusForbidden)
            return
        }

        // 2. when device hash is found and equal to submitted, then pass

    } else {
        // 3. when no device hash found, this is the first use
        a.orm.Model(&authid).Update(model.ColDevHash, dHash)
    }

    w.Header().Set("Server", "PocketCluster API Service")
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Cache-Control", "max-age=3600") // 1 hr
    w.Header().Set("Expires", time.Now().UTC().Add(time.Hour).Format("Mon, 2 Jan 2006 15:04:05 MST"))
    json.NewEncoder(w).Encode(map[string]string{"auth":"pass", "error":""})
    w.WriteHeader(http.StatusOK)
}