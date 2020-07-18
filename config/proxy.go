package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

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

func NewProxy(env string) (*Proxy, error) {
	bytes, err := ioutil.ReadFile(configDir + "config.json")
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	for _, p := range config.Proxies {
		if p.Env == env {
			return p, nil
		}
	}

	return nil, ErrEnvNotFound
}

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
