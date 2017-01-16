package util

import (
    "regexp"
    "strings"
    "io/ioutil"

    "github.com/PuerkitoBio/goquery"
)

func GithubReadmeScrap(location string, filename string) (string, error) {
    const anchorIcon string = "<svg version=\"1.1\" viewBox=\"0 0 8 8\" height=\"12\" width=\"12\"><path d=\"M5.88.03c-.18.01-.36.03-.53.09-.27.1-.53.25-.75.47a.5.5 0 1 0 .69.69c.11-.11.24-.17.38-.22.35-.12.78-.07 1.06.22.39.39.39 1.04 0 1.44l-1.5 1.5c-.44.44-.8.48-1.06.47-.26-.01-.41-.13-.41-.13a.5.5 0 1 0-.5.88s.34.22.84.25c.5.03 1.2-.16 1.81-.78l1.5-1.5c.78-.78.78-2.04 0-2.81-.28-.28-.61-.45-.97-.53-.18-.04-.38-.04-.56-.03zm-2 2.31c-.5-.02-1.19.15-1.78.75l-1.5 1.5c-.78.78-.78 2.04 0 2.81.56.56 1.36.72 2.06.47.27-.1.53-.25.75-.47a.5.5 0 1 0-.69-.69c-.11.11-.24.17-.38.22-.35.12-.78.07-1.06-.22-.39-.39-.39-1.04 0-1.44l1.5-1.5c.4-.4.75-.45 1.03-.44.28.01.47.09.47.09a.5.5 0 1 0 .44-.88s-.34-.2-.84-.22z\"/></svg>"
    var (
        httpCheck = regexp.MustCompile(`(^http://|^https://)`)
        repoName = strings.Replace(strings.ToLower(location), "https://github.com", "", 1)

        // walk the directory entries for indexing
        // http://stackoverflow.com/questions/37290693/how-to-remove-redundant-spaces-whitespace-from-a-string-in-golang
        cndnsLead   = regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
        cndnsInside = regexp.MustCompile(`[\s\p{Zs}]{2,}`)
        textContent string = ""
    )

    doc, err := goquery.NewDocument(location)
    if err != nil {
        return textContent, err
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

    // fix all img sources
    readme.Find("img").Each(func(_ int, s *goquery.Selection) {
        src, exists := s.Attr("src")
        if exists && !strings.HasPrefix(strings.ToLower(src), "http") {
            s.SetAttr("src", "https://github.com/" + src)
        }
    })

    // fix all link
    readme.Find("a").Each(func(_ int, s *goquery.Selection) {
        if s.HasClass("anchor") {
            // remove all anchors svg
            s.Find("svg").Each(func(i int, svg *goquery.Selection) {
                svg.ReplaceWithHtml(anchorIcon)
            })

            // remove user-content- header
            anchor, exists := s.Attr("id")
            if exists && strings.HasPrefix(strings.ToLower(anchor), "user-content-") {
                s.SetAttr("id", strings.Replace(anchor, "user-content-", "", -1))
            }
            return
        }

        // remove all links
        href, exists := s.Attr("href")
        if exists {
            link := strings.ToLower(href)

            // http, https check
            idxs := httpCheck.FindIndex([]byte(link))
            if len(idxs) != 0 && idxs[0] == 0 {
                return
            }

            // check if link starts with anchor: following two should not be switched
            if strings.HasPrefix(link, "#user-content-") {
                s.SetAttr("href", strings.Replace(link, "user-content-", "", -1))
                return
            }
            if strings.HasPrefix(link, "#") {
                return
            }

            // check if link starts with /user/reponame
            if strings.HasPrefix(link, repoName) {
                s.SetAttr("href", "https://github.com" + href)
                return
            }

            // check if link starts with user/reponame
            if strings.HasPrefix(link, repoName[1:]) {
                s.SetAttr("href", "https://github.com/" + href)
                return
            }
        }
    });
    // save readme html
    html, err := readme.Html()
    if err != nil {
        return textContent, err
    }
    // save to file
    err = ioutil.WriteFile(filename, []byte(html), 0444)
    if err != nil {
        return textContent, err
    }

    /* --- --- --- TOKENIZING HTML FOR SEARCH INDEXING --- --- --- */
    // remove all code for indexing
    readme.Find("pre").Each(func(_ int, s *goquery.Selection) {
        s.Empty().Remove()
    })
    //textContent = strings.TrimSpace(buffer.String())
    textContent = cndnsLead.ReplaceAllString(readme.Text(), "")
    // newline should be preserved
    //textContent = strings.Replace(data,"\n"," ", -1)
    textContent = strings.Replace(textContent,"\t"," ", -1)
    textContent = cndnsInside.ReplaceAllString(textContent, " ")
    return textContent, nil
}
