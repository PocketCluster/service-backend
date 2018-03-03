package auth

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"

    "github.com/pkg/errors"
    "github.com/jinzhu/gorm"
    "golang.org/x/crypto/ripemd160"

    "github.com/stkim1/api/auth/model"
    "github.com/stkim1/sharedpkg/randstr"
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
        entries = append(entries, strings.TrimSpace(entry[1]))
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

        var authid model.AuthIdentity
        orm.Where(fmt.Sprintf("%s = ?", model.ColEmail), req).First(&authid)

        if len(authid.Email) != 0 {
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