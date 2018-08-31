package config

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path"

	mlog "github.com/exfly/manageme/log"
	"github.com/spf13/viper"
)

var ErrNotFound = errors.New("Not Found")

func LoadConfig(filename string) error {
	file, err := os.Open(path.Join(".", filename))
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(file)
	// mlog.DEBUG("\n%v", string(content))
	if err != nil {
		file.Close()
		return err
	}
	file.Close()
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(content))
	mlog.INFO("loaded configfile %s", filename)
	return nil
}
