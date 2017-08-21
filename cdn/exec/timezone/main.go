package main

import (
    "time"
    "io/ioutil"
    "errors"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
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
        log.Fatal(trace.Wrap(err))
    }
    log.Info(tz.String())
}