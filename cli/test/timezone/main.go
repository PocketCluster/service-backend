package main

import (
    "time"
    "io/ioutil"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
)

func localTimezone() (*time.Location, error) {
    systz, err := ioutil.ReadFile("/etc/timezone")
    if len(systz) == 0 {
        return nil, errors.New("empty system timezone string")
    }
    if err != nil {
        return nil, err
    }
    localtz, err := time.LoadLocation(string(systz))
    if err != nil {
        return nil, err
    }
    return localtz, nil
}

func main() {
    tz, err := localTimezone()
    if err != nil {
        log.Fatal(errors.WithStack(err))
    }
    log.Info(tz.String())
}