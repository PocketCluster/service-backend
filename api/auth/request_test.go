package auth

import (
    "runtime"
    "testing"
    "path/filepath"
)

func testCSVFilename() string {
    var (
        _, testfile, _, _= runtime.Caller(0)
    )
    return filepath.Join(filepath.Dir(testfile), "request_test.csv")
}

func testEmptyCSVFilename() string {
    var (
        _, testfile, _, _= runtime.Caller(0)
    )
    return filepath.Join(filepath.Dir(testfile), "empty_request_test.csv")
}

func Test_CSV_Reading_Fail(t *testing.T) {
    if _, err := readRequestCSV(""); err == nil {
        t.Error("absent file should generate error")
    } else {
        t.Log(err.Error())
    }
}

func Test_Empty_CSV_Reading_Fail(t *testing.T) {
    if req, err := readRequestCSV(testEmptyCSVFilename()); err == nil {
        t.Errorf("empty file should generate error | %v", req)
    } else {
        t.Log(err.Error())
    }
}

func Test_CSV_Reading(t *testing.T) {
    csvfile := testCSVFilename()
    req, err := readRequestCSV(csvfile)
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }
    t.Logf("%v | %v", csvfile, req)
}
