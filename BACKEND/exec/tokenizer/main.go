package main

import (
    "bytes"
    "strings"
    "io/ioutil"
    "regexp"

    "golang.org/x/net/html"
    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
)

func main() {

    readme, err := ioutil.ReadFile("./readme/deanwampler-spark-workshop.html")
    if err != nil {
        log.Fatal(trace.Wrap(err))
    }

    // http://stackoverflow.com/questions/37290693/how-to-remove-redundant-spaces-whitespace-from-a-string-in-golang
    cndnsLead   := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
    cndnsInside := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
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
    //data = strings.Replace(data,"\n"," ", -1)
    data = strings.Replace(data,"\t"," ", -1)
    data = cndnsInside.ReplaceAllString(data, " ")

    log.Info(data)
}