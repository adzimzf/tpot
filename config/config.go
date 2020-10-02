package config

import (
	"encoding/json"
	"io/ioutil"
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

// AppendNode append the n to the proxy node list
func (p *Proxy) AppendNode(n Node) (Node, error) {

	pNode, err := p.GetNode()
	if err != nil {
		return pNode, err
	}
	for _, pn := range n.Items {
		var found bool
		for _, ni := range pNode.Items {
			if ni.Hostname == pn.Hostname {
				found = true
			}
		}
		if !found {
			pNode.Items = append(pNode.Items, pn)
		}
	}
	return pNode, nil
}

// GetNode get the node from proxy cache
func (p *Proxy) GetNode() (Node, error) {
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
