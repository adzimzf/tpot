package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/manifoldco/promptui"
)

var configDir = os.Getenv("HOME") + "/.tpot/"

type Config struct {
	Proxies []*Proxy `json:"proxies"`
}

func AddConfig() (err error) {
	if err := addConfigDirExist(); err != nil {
		return err
	}

	var cfg Config
	cfgPath := configDir + "config.json"
	file, err := ioutil.ReadFile(cfgPath)
	if err == nil {
		err := json.Unmarshal(file, &cfg)
		if err != nil {
			return err
		}
	}

	p := &Proxy{}
	NewProxySetterStations().Execute(p)
	cfg.Proxies = append(cfg.Proxies, p)

	bytes, err := json.MarshalIndent(cfg, "", "	")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cfgPath, bytes, permission)
}

func addConfigDirExist() error {
	_, err := os.Stat(configDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(configDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return err

}

func prompt(label string, validate func(string2 string) error) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	return prompt.Run()
}
