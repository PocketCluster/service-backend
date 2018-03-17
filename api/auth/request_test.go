package auth

import (
    "os"
    "path/filepath"
    "reflect"
    "runtime"
    "testing"
    "io/ioutil"
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

func testInvRecFilename() string {
    var (
        _, testfile, _, _= runtime.Caller(0)
    )
    return filepath.Join(filepath.Dir(testfile), "invitation_record_test.csv")
}

func Test_CSV_Reading_Fail(t *testing.T) {
    if _, err := readRequestCSV(""); err == nil {
        t.Error("absent file should generate error")
        t.FailNow()
    } else {
        t.Log(err.Error())
    }
}

func Test_Empty_CSV_Reading_Fail(t *testing.T) {
    if req, err := readRequestCSV(testEmptyCSVFilename()); err == nil {
        t.Errorf("empty file should generate error | %v", req)
        t.FailNow()
    } else {
        t.Log(err.Error())
    }
}

func Test_CSV_Reading(t *testing.T) {
    reqs, err := readRequestCSV(testCSVFilename())
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }
    var expected = []string{"ldgoodbrod@gmail.com", "socjuan.jbt@gmail.com", "diegobarrioh@gmail.com"}
    if !reflect.DeepEqual(reqs, expected) {
        t.Errorf("email list does not match expected | actual %v, expected %v", reqs, expected)
        t.FailNow()
    }
}

func Test_Invitation_CSV_Record(t *testing.T) {
    orm, err := openTestOrm()
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }
    defer closeTestOrm(orm)

    reqs, err := readRequestCSV(testCSVFilename())
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }

    err = updateRequestRecord(reqs, orm)
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }

    err = recordInvitation(orm, testInvRecFilename())
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }

    actual, err := ioutil.ReadFile(testInvRecFilename())
    if err != nil {
        t.Error(err.Error())
        t.FailNow()
    }
    t.Log(string(actual))

    os.Remove(testInvRecFilename())
}
