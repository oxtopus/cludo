package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

type CludoConfig struct {
	BaseImage string
}

func readConfig(dir string) (CludoConfig, error) {
	var config CludoConfig
	var err error
	filename := path.Join(dir, "cludo.json")

	_, err = os.Stat(filename)
	if err == nil {
		file, e := ioutil.ReadFile(filename)

		if e != nil {
			return config, e
		}

		err := json.Unmarshal(file, &config)

		if err != nil {
			return config, err
		}

		return config, nil
	}

	return config, err
}

func makeDefaultConfig() CludoConfig {
	return CludoConfig{"cludo-base"}
}

func MakeConfig(wd string) (CludoConfig, error) {
	var config CludoConfig
	var err error

	config, err = readConfig(wd)

	if err != nil {
		usr, _ := user.Current()
		config, err = readConfig(usr.HomeDir)
		if err != nil {
			config = makeDefaultConfig()
		}
	}

	return config, err
}
