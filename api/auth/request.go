package auth

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "regexp"
    "strings"

    "github.com/pkg/errors"
    "github.com/jinzhu/gorm"
    "golang.org/x/crypto/ripemd160"

    "github.com/stkim1/api/auth/model"
    "github.com/stkim1/sharedpkg/randstr"
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
        ic := randstr.NewCapRandString(16)
        full := fmt.Sprintf("%s-%s-%s-%s", ic[0:4], ic[4:8], ic[8:12], ic[12:])
        hasher := ripemd160.New()
        hasher.Write([]byte(full))
        iHash := fmt.Sprintf("%x\n", hasher.Sum(nil))

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

    invrec, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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

