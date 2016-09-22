package hunter

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v1"
)

type Config struct {
	Common struct {
		Version  string
		IsDebug  bool `yaml:"debug"`
		LogPath  string
		LogLevel string
	}

	Hunter struct {
		IsDaemon bool
		PullTime int
		Date     string
		OutPath  string
	}

	BaseFilter struct {
		SuffixType   string
		MaxLen       int
		ExcludeDelim bool
		CharType     int
	}

	AdvFilter struct {
		OccurChars int
	}
}

var Conf = &Config{}

func init() {
	data, err := ioutil.ReadFile("dropcatch.yaml")
	if err != nil {
		log.Fatal("read config error :", err)
	}

	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		log.Fatal("yaml decode error :", err)
	}

	log.Println(Conf)
}
