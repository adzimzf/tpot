package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type Proxy struct {
	Address  string `json:"address"`
	UserName string `json:"user_name"`
	Env      string `json:"env"`
	TwoFA    bool   `json:"two_fa"`

	// For using OAUTH like GMAIL, Facebook etc
	// empty means using username & password
	AuthConnector string `json:"auth_connector"`
	Node          Node
}

type Node struct {
	Items []Item `json:"items"`
}

// LookUpIPAddress lookup the IP address by host
func (n *Node) LookUpIPAddress(host string) (string, bool) {
	for _, i2 := range n.Items {
		if i2.Hostname == host {
			return strings.Split(i2.Hostname, ":")[0], true
		}
	}
	return "", false
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
		return Node{}, err
	}
	err = json.Unmarshal(nodeBytes, &p.Node)
	if err != nil {
		return p.Node, err
	}
	return p.Node, nil
}

// UpdateNode update the cache node
func (p *Proxy) UpdateNode(n Node) error {
	bytes, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configDir+"node_"+p.Env+".json", bytes, permission)
}
