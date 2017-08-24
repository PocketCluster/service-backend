package cforigin

import (
    "net/http"

    "github.com/pkg/errors"
    "strings"
)

var (
    //https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
    allowedCountries map[string]bool = map[string]bool {
        "US": true,
        "CA": true,
        "MX": true,
    }
)

// these configuration values heavily depend on api.pocketcluster.io nginx configuration.
// When the configuration changes, please make proper adjustements here.
func IsOriginAllowedCountry(req *http.Request) error {
    var (
        origin = req.Header.Get("cf-ipcountry")
        realip = req.Header.Get("x-real-ip")
    )
    // this is debug connection
    if len(realip) != 0 && strings.TrimSpace(realip) == "198.199.115.209" {
        return nil
    }
    if len(origin) == 0 || len(origin) != 2 {
        return errors.Errorf("invalid country form or server is not covered with CloudFlare")
    }

    for c, _ := range allowedCountries {
        if c == origin {
            return nil
        }
    }

    return errors.Errorf("originated country %v is now allowed to access service", origin)
}

