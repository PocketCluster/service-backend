package auth

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "regexp"
    "strings"
    "sync"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/jinzhu/gorm"
    "github.com/pkg/errors"
    "golang.org/x/crypto/ripemd160"

    "github.com/stkim1/service-backend/api/auth/model"
    "github.com/stkim1/service-backend/shared/randstr"
)

const (
    regEmailCheck string = "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\\])"
)

func readRequestCSV(filename string) ([]string, error) {
    csvf, err := os.Open(filename)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    defer csvf.Close()

    var (
        csvr = csv.NewReader(csvf)
        entries []string = nil
        headerSkipped = false
    )
    for {
        entry, err := csvr.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, errors.WithStack(err)
        }
        if !headerSkipped {
            headerSkipped = true
            continue
        }
        email := strings.TrimSpace(entry[1])
        match, _ := regexp.MatchString(regEmailCheck, email)
        if match {
            entries = append(entries, email)
        }
    }
    if len(entries) == 0 {
        return nil, errors.New("empty request")
    }
    return entries, nil
}

func updateRequestRecord(requester []string, orm *gorm.DB) error {
    if len(requester) == 0 {
        return errors.New("no requester in the list")
    }
    if orm == nil {
        return errors.New("invalid data repository")
    }

    for _, req := range requester {

        // check if email exist
        var authid model.AuthIdentity
        orm.Where(fmt.Sprintf("%s = ?", model.ColEmail), req).First(&authid)

        if len(authid.Email) != 0 && len(authid.Invitation) == 19 && len(authid.InvHash) == 40 {
            continue
        }

        // generate
        ic := strings.ToUpper(randstr.NewRandomString(20))
        full := fmt.Sprintf("%s-%s-%s-%s-%s", ic[0:4], ic[4:8], ic[8:12], ic[12:16], ic[16:])
        hasher := ripemd160.New()
        hasher.Write([]byte(full))
        iHash := fmt.Sprintf("%x", hasher.Sum(nil))

        // save to orm
        authid.Email = req
        authid.Invitation = full
        authid.InvHash = iHash
        orm.Create(&authid)
    }

    return nil
}

func recordInvitation(orm *gorm.DB, filename string) error {
    if orm == nil {
        return errors.New("invalid data repository")
    }

    var authlist []model.AuthIdentity = nil
    orm.Find(&authlist)
    if len(authlist) == 0 {
        return errors.New("empty invitation list")
    }

    invrec, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return errors.WithStack(err)
    }
    defer invrec.Close()
    invrec.Seek(0, 0)
    invrec.Write([]byte(fmt.Sprintf("Email, Invitation\n")))

    for _, a := range authlist {
        invrec.Write([]byte(fmt.Sprintf("%v, %v\n", a.Email, a.Invitation)))
    }
    return nil
}

func RefreshInvitationList(wg *sync.WaitGroup, isTermC <- chan interface{}, orm *gorm.DB, reqcsv, invcsv string) error {
    if wg == nil {
        return errors.New("invalid workgroup")
    }
    if orm == nil {
        return errors.New("invalid orm")
    }
    if len(reqcsv) == 0 {
        return errors.New("invalid request list")
    }
    if len(invcsv) == 0 {
        return errors.New("invalid request list")
    }

    go func(fwg *sync.WaitGroup, fisTermC <- chan interface{}, form *gorm.DB, freqcsv, finvcsv string) {
        var reftick = time.NewTicker(time.Minute * 5)

        // update invitation list as soon as start
        if reqs, rerr := readRequestCSV(freqcsv); rerr == nil {
            if ierr := updateRequestRecord(reqs, form); ierr == nil {
                if cerr := recordInvitation(orm, finvcsv); cerr != nil {
                    log.Error(cerr.Error())
                }
            } else {
                log.Error(ierr.Error())
            }
        } else {
            log.Error(rerr.Error())
        }

        for {
            select {
                case <- fisTermC: {
                    reftick.Stop()
                    fwg.Done()
                }
                case <- reftick.C: {
                    // update invitation list as soon as start
                    if reqs, rerr := readRequestCSV(freqcsv); rerr == nil {
                        if ierr := updateRequestRecord(reqs, form); ierr == nil {
                            if cerr := recordInvitation(orm, finvcsv); cerr != nil {
                                log.Error(cerr.Error())
                            }
                        } else {
                            log.Error(ierr.Error())
                        }
                    } else {
                        log.Error(rerr.Error())
                    }
                }
            }
        }
    }(wg, isTermC, orm, reqcsv, invcsv)

    return nil
}