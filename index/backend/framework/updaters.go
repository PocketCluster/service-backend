package framework

import (
    "io/ioutil"
    "sync"
    "sync/atomic"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/jinzhu/gorm"
    "github.com/blevesearch/bleve"

    "github.com/stkim1/service-backend/index/backend/storage"
    "github.com/stkim1/service-backend/index/backend/config"
    "github.com/stkim1/service-backend/index/backend/update"
    "github.com/stkim1/service-backend/index/backend/model"
)

func (a *Application) ScheduleMetaUpdate() {
    go func (quit <- chan bool, metaWaiter *sync.WaitGroup, metaDB *gorm.DB, searchIndex bleve.Index, cfg *config.Config, totalRepoCount, isUpdating *atomic.Value) {
        var (
            updateTicker *time.Ticker   = time.NewTicker(time.Minute)
            lastRec time.Time
            repoCount int64 = 0
            oldTS []byte
            err error
        )

        metaWaiter.Add(1)
        oldTS, err = ioutil.ReadFile(cfg.Update.MetaUpdateRecord)
        if err == nil {
            lastRec, err = time.Parse(time.RFC3339, string(oldTS))
        }
        if err != nil {
            lastRec = time.Now().Add((time.Minute * time.Duration(30) - time.Minute * time.Duration(cfg.Update.MetaUpdateInterval)))
        }
        log.Info("Last Meta Updated : " + lastRec.String())

        for {
            select {
            case launch := <- updateTicker.C: {
                    if !isUpdating.Load().(bool) && (time.Minute * time.Duration(cfg.Update.MetaUpdateInterval)) < launch.Sub(lastRec) {
                        // update repo count
                        metaDB.Model(&model.Repository{}).Count(&repoCount)
                        totalRepoCount.Store(repoCount)

                        // update (2017/10/10 : we'll give group-read permission for backup)
                        ioutil.WriteFile(cfg.Update.MetaUpdateRecord, []byte(launch.Format(time.RFC3339)), 0640)
                        lastRec = launch
                        go update.UpdateAllRepoMeta(metaDB, searchIndex, cfg, metaWaiter, isUpdating)
                    }
                }

            case <- quit: {
                    log.Info("[ScheduleMetaUpdate] time to quit...")
                    metaWaiter.Done()
                    updateTicker.Stop()
                    return
                }
            }
        }
    }(a.QuitMetaUpdate, &a.UpdateWait, a.MetaDB, a.SearchIndex, a.Config, &(a.Controller.TotalRepoCount), &a.IsMetaUpdating)
}

func (a *Application) ScheduleSuppUpdate() {
    go func (quit <- chan bool, suppWaiter *sync.WaitGroup, metaDB *gorm.DB, suppDB storage.Nosql, cfg *config.Config, isUpdating *atomic.Value) {
        var (
            updateTicker *time.Ticker   = time.NewTicker(time.Minute)
            lastRec time.Time
            oldTS []byte
            err error
        )

        suppWaiter.Add(1)
        oldTS, err = ioutil.ReadFile(cfg.Update.SuppUpdateRecord)
        if err == nil {
            lastRec, err = time.Parse(time.RFC3339, string(oldTS))
        }
        if err != nil {
            lastRec = time.Now().Add((time.Minute * time.Duration(60) - time.Minute * time.Duration(cfg.Update.SuppUpdateInterval)))
        }
        log.Info("Last Supp Updated : " + lastRec.String())

        for {
            select {
            case launch := <- updateTicker.C: {
                    if !isUpdating.Load().(bool) && (time.Minute * time.Duration(cfg.Update.SuppUpdateInterval)) < launch.Sub(lastRec) {
                        // (2017/10/10 : we'll give group-read permission for backup)
                        ioutil.WriteFile(cfg.Update.SuppUpdateRecord, []byte(launch.Format(time.RFC3339)), 0640)
                        lastRec = launch

                        go update.UpdateAllRepoSupplement(metaDB, suppDB, cfg, suppWaiter, isUpdating)
                    }
                }

            case <- quit: {
                    log.Info("[ScheduleSuppUpdate] time to quit...")
                    suppWaiter.Done()
                    updateTicker.Stop()
                    return
                }
            }
        }
    }(a.QuitSuppUpdate, &a.UpdateWait, a.MetaDB, a.SuppleDB, a.Config, &a.IsSuppUpdating)
}
