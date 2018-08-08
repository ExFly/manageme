package config

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"

	mlog "github.com/exfly/manageme/log"
	"github.com/spf13/viper"
)

var ErrNotFound = errors.New("Not Found")

func LoadConfig(filename string) {
	file, err := os.Open(path.Join(".", filename))
	content, err := ioutil.ReadAll(file)
	// mlog.DEBUG("\n%v", string(content))
	if err != nil {
		file.Close()
		log.Fatal(err)
	}
	file.Close()
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(content))
	mlog.INFO("loaded configfile %s", filename)
}
