package fmk

type Utils
var Utils Utils

func (utils Utils) parseURL(url string) []string {
        urlBytes := []byte(url)
        wordBuilder := make([]byte, 0)
        result := make([]string, 0)
        for _, char := range urlBytes {
                if char == FRONT_SLASH {
                        result = append(result, string(wordBuilder))
                        wordBuilder = make([]byte, 0)
                } else {
                        wordBuilder = append(wordBuilder, char)
                }
        }
        result = append(result, string(wordBuilder))
        return result
}

