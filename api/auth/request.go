package auth

import (
    "encoding/csv"
    "io"
    "os"

    "github.com/pkg/errors"
)

func readRequestCSV(filename string) ([]string, error) {
    csvf, err := os.Open(filename)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    defer csvf.Close()

    var (
        csvr = csv.NewReader(csvf)
        entries []string = nil
        headerSkipped = false
    )
    for {
        entry, err := csvr.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, errors.WithStack(err)
        }
        if !headerSkipped {
            headerSkipped = true
            continue
        }
        entries = append(entries, entry[1])
    }
    if len(entries) == 0 {
        return nil, errors.New("empty request")
    }
    return entries, nil
}