package framework

import (
    "sync"
    "time"
    "sync/atomic"

    log "github.com/Sirupsen/logrus"
    "github.com/jinzhu/gorm"

    "github.com/stkim1/BACKEND/storage"
    "github.com/stkim1/BACKEND/config"
)

func (a *Application) ScheduleMetaUpdate() {
    go func (quit <- chan bool, metaWaiter *sync.WaitGroup, metaDB *gorm.DB, cfg *config.Config, lastUpdated *time.Time) {
        var (
            isWorking atomic.Value
            updateTicker = time.NewTicker(time.Minute * 60)
        )

        //time.Duration(cfg.MetaUpdateInterval)
        metaWaiter.Add(1)
        isWorking.Store(false)

        for {
            select {
            case launch := <- updateTicker.C:
                if !isWorking.Load().(bool) {
                    log.Infof("[ScheduleMetaUpdate] %v should begin", launch)
                    //go updateSupplement(wg, &isWorking)
                }

            case <- quit: {
                log.Info("[ScheduleMetaUpdate] time to quit...")
                metaWaiter.Done()
                updateTicker.Stop()
                return
            }
            }
        }
    }(a.QuitMetaUpdate, &a.UpdateWait, a.MetaDB, a.Config)
}

func (a *Application) ScheduleSuppUpdate() {
    go func (quit <- chan bool, suppWaiter *sync.WaitGroup, suppDB storage.Nosql, cfg *config.Config, lastUpdated *time.Time) {
        var (
            isWorking atomic.Value
            updateTicker = time.NewTicker(time.Minute * 60)
        )

        // time.Duration(cfg.SuppUpdateInterval)
        suppWaiter.Add(1)
        isWorking.Store(false)

        for {
            select {
            case launch := <- updateTicker.C:
                if !isWorking.Load().(bool) {
                    log.Infof("[ScheduleSuppUpdate] %v should begin", launch)
                    //go updateSupplement(wg, &isWorking)
                }

            case <- quit: {
                    log.Info("[ScheduleSuppUpdate] time to quit...")
                    suppWaiter.Done()
                    updateTicker.Stop()
                    return
                }
            }
        }
    }(a.QuitSuppUpdate, &a.UpdateWait, a.SuppleDB, a.Config)
}
