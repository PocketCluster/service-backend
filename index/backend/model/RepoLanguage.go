package model

func MakeLanguageEntryKey(repoID string) []byte {
    return []byte("lang-" + repoID)
}

type RepoLanguage struct {
    // Programming Language
    Language          string        `msgpack:"lang"`
    // Percentage
    Percentage        float32       `msgpack:"percent"`
}

type ListLanguage []*RepoLanguage

func (slice ListLanguage) Len() int {
    return len(slice)
}

func (slice ListLanguage) Less(i, j int) bool {
    return slice[j].Percentage < slice[i].Percentage;
}

func (slice ListLanguage) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}
