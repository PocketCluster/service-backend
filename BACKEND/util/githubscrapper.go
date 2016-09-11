package util

import (
    "log"
    "strings"
    "io/ioutil"

    "github.com/PuerkitoBio/goquery"
)

func GithubReadmeScrap(location string, filename string) {
    const anchorIcon string = "<svg aria-hidden=\"true\" class=\"octicon octicon-link\" version=\"1.1\" viewBox=\"0 0 8 8\" height=\"12\" width=\"12\"><path d=\"M5.88.03c-.18.01-.36.03-.53.09-.27.1-.53.25-.75.47a.5.5 0 1 0 .69.69c.11-.11.24-.17.38-.22.35-.12.78-.07 1.06.22.39.39.39 1.04 0 1.44l-1.5 1.5c-.44.44-.8.48-1.06.47-.26-.01-.41-.13-.41-.13a.5.5 0 1 0-.5.88s.34.22.84.25c.5.03 1.2-.16 1.81-.78l1.5-1.5c.78-.78.78-2.04 0-2.81-.28-.28-.61-.45-.97-.53-.18-.04-.38-.04-.56-.03zm-2 2.31c-.5-.02-1.19.15-1.78.75l-1.5 1.5c-.78.78-.78 2.04 0 2.81.56.56 1.36.72 2.06.47.27-.1.53-.25.75-.47a.5.5 0 1 0-.69-.69c-.11.11-.24.17-.38.22-.35.12-.78.07-1.06-.22-.39-.39-.39-1.04 0-1.44l1.5-1.5c.4-.4.75-.45 1.03-.44.28.01.47.09.47.09a.5.5 0 1 0 .44-.88s-.34-.2-.84-.22z\"/></svg>&nbsp;"

    doc, err := goquery.NewDocument(location)
    if err != nil {
        log.Fatal(err)
    }

    readme := doc.Find("#readme").Clone()

    // remove "README" tag
    readme.Find("h3").EachWithBreak(func(i int, s *goquery.Selection) bool {
        if strings.Contains(strings.ToLower(s.Text()), "readme") {
            s.Remove()
            return false
        }
        return true
    })

    // remove all anchors svg
    readme.Find("a[href].anchor svg").Each(func(i int, s *goquery.Selection) {
        s.ReplaceWithHtml(anchorIcon)
    })

    // fix all img sources
    // TODO: fix <a href> with <img> tag
    readme.Find("img").Each(func(_ int, s *goquery.Selection) {
        src, exists := s.Attr("src")
        if !strings.HasPrefix(src, "http") && exists {
            s.SetAttr("src", "https://github.com/" + src)
        }
    })

    // img a.tag w/ image fix
    /*
        readme.Find("a[href] img").Each(func(_ int, s *goquery.Selection) {
            href, _ := s.Parent().Attr("href")
            if err == nil && !strings.HasPrefix(href,"http") {
                if strings.HasPrefix(href, "/") {
                    s.Parent().SetAttr("href", location + href)
                } else {
                    s.Parent().SetAttr("href", location + "/" + href)
                }
            }
        });
    */
    readme.Find("a[href]").Not(".anchor").Each(func(_ int, s *goquery.Selection) {
        href, _ := s.Attr("href")
        if err == nil && !strings.HasPrefix(href,"http") {
            if strings.HasPrefix(href, "/") {
                s.Parent().SetAttr("href", location + href)
            } else {
                s.Parent().SetAttr("href", location + "/" + href)
            }
        }
    });

    // read html
    html, err := readme.Html()
    if err != nil {
        log.Panic("Cannot read HTML")
    }

    // save to file
    err = ioutil.WriteFile("readme/" + filename, []byte(html), 0664)
    if err != nil {
        log.Panic("Cannot save HTML readme " + err.Error())
    }
}
