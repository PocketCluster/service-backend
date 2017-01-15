package main

import (
    "bytes"
    "strings"
    "io/ioutil"
    "regexp"

    "golang.org/x/net/html"
    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "time"
    "path"
)

func main() {

    const targetDir string = "./readme"

    // open the directory
    dirEntries, err := ioutil.ReadDir(targetDir)
    if err != nil {
        log.Fatal(trace.Wrap(err))
        return
    }

    // walk the directory entries for indexing
    // http://stackoverflow.com/questions/37290693/how-to-remove-redundant-spaces-whitespace-from-a-string-in-golang
    cndnsLead   := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
    cndnsInside := regexp.MustCompile(`[\s\p{Zs}]{2,}`)

    log.Printf("Tokenizing...")
    count := 0
    startTime := time.Now()

    for _, dirEntry := range dirEntries {
        filename := dirEntry.Name()
        // read the bytes
        readme, err := ioutil.ReadFile(path.Join(targetDir, filename))
        if err != nil {
            log.Error(trace.Wrap(err))
            continue
        }

        var buffer bytes.Buffer
        content := html.NewTokenizer(strings.NewReader(string(readme)))
        for {
            // token type
            tokenType := content.Next()
            if tokenType == html.ErrorToken {
                break
            }
            token := content.Token()
            switch tokenType {
            case html.StartTagToken: // <tag>
            // type Token struct {
            //     Type     TokenType
            //     DataAtom atom.Atom
            //     Data     string
            //     Attr     []Attribute
            // }
            //
            // type Attribute struct {
            //     Namespace, Key, Val string
            // }
            case html.TextToken: // text between start and end tag
                buffer.WriteString(token.Data)

            case html.EndTagToken: // </tag>
            case html.SelfClosingTagToken: // <tag/>
            }
        }

        //data := strings.TrimSpace(buffer.String())
        data := cndnsLead.ReplaceAllString(buffer.String(), "")
        // newline should be preserved
        //data = strings.Replace(data,"\n"," ", -1)
        data = strings.Replace(data,"\t"," ", -1)
        data = cndnsInside.ReplaceAllString(data, " ")

        ioutil.WriteFile(path.Join(targetDir, strings.Replace(filename, ".html", ".txt", -1)), []byte(data), 0600)

        count++
        if count%100 == 0 {
            indexDuration := time.Since(startTime)
            indexDurationSeconds := float64(indexDuration) / float64(time.Second)
            timePerDoc := float64(indexDuration) / float64(count)
            log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
        }
    }

    indexDuration := time.Since(startTime)
    indexDurationSeconds := float64(indexDuration) / float64(time.Second)
    timePerDoc := float64(indexDuration) / float64(count)
    log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
}