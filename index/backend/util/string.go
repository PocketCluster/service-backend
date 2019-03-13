package util

import (
    "regexp"
    "strings"
)

func RegexCapitalize(input string) string {
    // Function replacing words (assuming lower case input)
    replace := func(word string) string {
        switch word {
        case "with", "in", "a", "an", "to", "on", "the":
            return word
        }
        return strings.Title(word)
    }

    r := regexp.MustCompile(`\w+`)
    return r.ReplaceAllStringFunc(strings.ToLower(input), replace)
}

func Capitalize(input string) string {
    var (
        words []string    = strings.Fields(input)
        smallwords string = " with in a an to on the "
    )
    for index, word := range words {
        if strings.Contains(smallwords, " "+word+" ") {
            words[index] = word
        } else {
            words[index] = strings.Title(word)
        }
    }
    return strings.Join(words, " ")
}