package util

import (
    "log/syslog"
    "io/ioutil"

    log "github.com/Sirupsen/logrus"
    logrusSyslog "github.com/Sirupsen/logrus/hooks/syslog"
)

func InitSysLogger() {
    log.SetLevel(log.DebugLevel)
    // clear existing hooks:
    log.StandardLogger().Hooks = make(log.LevelHooks)
    log.SetFormatter(&log.TextFormatter{})

    hook, err := logrusSyslog.NewSyslogHook("", "", syslog.LOG_DEBUG, "")
    if err != nil {
        // syslog not available
        log.Warn("syslog not available. reverting to stderr")
    } else {
        // ... and disable stderr:
        log.AddHook(hook)
        log.SetOutput(ioutil.Discard)
    }
}

