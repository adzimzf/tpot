package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"strconv"

	"github.com/manifoldco/promptui"
)

var configDir = os.Getenv("HOME") + "/.tpot/"

type Config struct {
	Proxies []Proxy `json:"proxies"`
}

type Proxy struct {
	Address  string `json:"address"`
	UserName string `json:"user_name"`
	Env      string `json:"env"`
	TwoFA    bool   `json:"two_fa"`
	Node     Node
}

type Node struct {
	Items []Item `json:"items"`
}

type Item struct {
	Hostname string `json:"hostname"`
	Address  string `json:"addr"`
}

var ErrEnvNotFound = fmt.Errorf("env not found")

const permission = 775

func (p *Proxy) UpdateNode(n Node) error {
	bytes, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configDir+"node_"+p.Env+".json", bytes, permission)
}

func (p *Proxy) LoadNode() error {
	nodeBytes, err := ioutil.ReadFile(configDir + "node_" + p.Env + ".json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(nodeBytes, &p.Node)
	if err != nil {
		return err
	}
	return nil
}

func NewProxy(env string) (*Proxy, error) {
	bytes, err := ioutil.ReadFile(configDir + "config.json")
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(bytes, &config)
	for _, p := range config.Proxies {
		if p.Env == env {
			return &p, nil
		}
	}
	return nil, ErrEnvNotFound
}

func AddConfig() (err error) {

	if err := addConfigDirExist(); err != nil {
		return err
	}

	cfgPath := configDir + "config.json"

	var p Proxy
	p.Env, err = prompt("Environment", func(string2 string) error {
		return nil
	})
	if err != nil {
		return err
	}
	p.Address, err = prompt("Proxy Address (with http protocol)", func(string2 string) error {
		_, err := url.ParseRequestURI(string2)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	p.UserName, err = prompt("Username (teleport username)", func(string2 string) error {
		return nil
	})
	if err != nil {
		return err
	}
	isTwoFA, err := prompt("Is Need 2FA (Y/y/N/n)", func(string2 string) error {
		if string2 == "Y" || string2 == "y" || string2 == "N" || string2 == "n" {
			return nil
		}
		return fmt.Errorf("invalid formatting")
	})
	if err != nil {
		return err
	}

	if isTwoFA == "Y" || isTwoFA == "y" {
		p.TwoFA = true
	}

	var cfg Config
	file, err := ioutil.ReadFile(cfgPath)
	if err == nil {
		err := json.Unmarshal(file, &cfg)
		if err != nil {
			return err
		}
	}

	cfg.Proxies = append(cfg.Proxies, p)

	bytes, err := json.Marshal(cfg)
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

func isRoot() (bool, error) {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	i, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		return false, err
	}

	return i == 0, nil
}
