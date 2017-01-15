package framework

import (
    "io/ioutil"
    "sync"
    "sync/atomic"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/jinzhu/gorm"

    "github.com/stkim1/BACKEND/storage"
    "github.com/stkim1/BACKEND/config"
    "github.com/stkim1/BACKEND/update"
)

func (a *Application) ScheduleMetaUpdate() {
    go func (quit <- chan bool, metaWaiter *sync.WaitGroup, metaDB *gorm.DB, cfg *config.Config, isUpdating *atomic.Value) {
        var (
            updateTicker *time.Ticker   = time.NewTicker(time.Minute)
            lastRec time.Time
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
        log.Info("Last Meta update" + lastRec.String())

        for {
            select {
            case launch := <- updateTicker.C: {
                    if !isUpdating.Load().(bool) && (time.Minute * time.Duration(cfg.Update.MetaUpdateInterval)) < launch.Sub(lastRec) {
                        ioutil.WriteFile(cfg.Update.MetaUpdateRecord, []byte(launch.Format(time.RFC3339)), 0600)
                        lastRec = launch

                        go update.UpdateAllRepoMeta(metaDB, cfg, metaWaiter, isUpdating)
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
    }(a.QuitMetaUpdate, &a.UpdateWait, a.MetaDB, a.Config, &a.IsMetaUpdating)
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
        log.Info("Last Supp update" + lastRec.String())

        for {
            select {
            case launch := <- updateTicker.C: {
                    if !isUpdating.Load().(bool) && (time.Minute * time.Duration(cfg.Update.SuppUpdateInterval)) < launch.Sub(lastRec) {
                        ioutil.WriteFile(cfg.Update.SuppUpdateRecord, []byte(launch.Format(time.RFC3339)), 0600)
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
