package main

import (
    "os"
    "path"
    "regexp"
    "strings"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/PuerkitoBio/goquery"
    "github.com/blevesearch/bleve"

    "github.com/stkim1/backend/config"
    "github.com/stkim1/backend/model"
    pocketsearch "github.com/stkim1/backend/search"

)

func readmeTokenizing(filename string) (string, error) {
    var (
        cndnsLead          = regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
        cndnsInside        = regexp.MustCompile(`[\s\p{Zs}]{2,}`)
        textContent string = ""
    )

    f, err := os.Open(filename)
    if err != nil {
        return textContent, err
    }
    defer f.Close()

    readme, err := goquery.NewDocumentFromReader(f)
    if err != nil {
        return textContent, err
    }
    readme.Find("pre").Each(func(_ int, s *goquery.Selection) {
        s.Empty().Remove()
    })
    textContent = cndnsLead.ReplaceAllString(readme.Text(), "")
    textContent = strings.Replace(textContent,"\t"," ", -1)
    textContent = cndnsInside.ReplaceAllString(textContent, " ")
    return textContent, nil
}

func main() {
    var (
        rsIndx bleve.Index
        err error
    )
    // config
    cfgPath, ok := os.LookupEnv(config.EnvConfigFilePath)
    if !ok {
        cfgPath = "config.yaml"
    }
    cfg, err := config.NewConfig(cfgPath)
    if err != nil {
        log.Fatal(err)
    }

    rsIndx, err = bleve.Open(cfg.Search.IndexStoragePath)
    if err != nil {
        m, err := pocketsearch.BuildIndexMapping()
        if err != nil {
            log.Fatal(err)
        }
        rsIndx, err = bleve.New(cfg.Search.IndexStoragePath, m)
        if err != nil {
            log.Fatal(err)
        }
    }
    defer rsIndx.Close()

    // database
    metaDB, err := gorm.Open(cfg.Database.DatabaseType, cfg.Database.DatabasePath)
    if err != nil {
        log.Fatal(errors.WithStack(err))
    }
    defer metaDB.Close()

    var repos []model.Repository
    metaDB.Find(&repos)
    for i, meta := range repos {
        readme, err := readmeTokenizing(path.Join(cfg.General.ReadmePath, meta.Slug + ".html"))
        if err != nil {
            log.Error(err)
        } else {
            sr := pocketsearch.NewSerachRepo(&(repos[i]), &readme)
            sr.Index(rsIndx)
        }
    }

}
