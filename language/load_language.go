package language

import (
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type MultiLocalizer map[string]Localizer

const DEFAULT_LANGUAGE = "en"

func (s MultiLocalizer) GetLocalizer(req *http.Request) Localizer {
	headerLocalizer := req.Header.Get("language")
	if len(headerLocalizer) == 0 {
		headerLocalizer = DEFAULT_LANGUAGE
	}
	return s[headerLocalizer]
}

func (s MultiLocalizer) Get(l string) Localizer {
	value, ok := s[l]
	if !ok {
		return s[DEFAULT_LANGUAGE]
	}
	return value
}

func (s MultiLocalizer) GetDefault() Localizer {
	return s[DEFAULT_LANGUAGE]
}

func LoadAllFileLanguage(path string) MultiLocalizer {
	res := make(MultiLocalizer)
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return res
	}

	for _, f := range files {
		loadLanguageFileYaml(res, path, f.Name())
	}
	return res
}

func loadLanguageFileYaml(multiLocalizer MultiLocalizer, path, language string) {
	yfile, err := os.ReadFile(path + language)

	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]string)

	err2 := yaml.Unmarshal(yfile, &data)

	if err2 != nil {
		log.Fatal(err2)
	}

	textLanguage := strings.Split(language, ".")[0]
	multiLocalizer[textLanguage] = data
}
