package util

import (
    "strings"
    "errors"
    "github.com/google/go-github/github"
    "time"
)

func SafeGetString(str *string) string {
    if str == nil {
        return ""
    }
    if len(*str) == 0 {
        return ""
    }
    return *str
}

func SafeGetInt(value *int) (int, error) {
    if value == nil {
        return 0, errors.New("Cannot read int value")
    }
    return *value, nil
}

func SafeGetBool(value *bool) bool {
    if value == nil {
        return false
    }
    return *value
}

func SafeStringJoin(params... string) string {
    return strings.Join(params, "")
}

func SafeGetTimestamp(value *github.Timestamp) time.Time {
    if value == nil {
        return time.Time{}
    }
    return value.Time
}

func SafeGetTime(value *time.Time) time.Time {
    if value == nil {
        return time.Time{}
    }
    return *value
}