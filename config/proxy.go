package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

const permission = 0600
const proxyTemplate = `
---
proxies:
  # environment name that will be use for accessing proxy
  # the recommendation is simple & easy to remember
  # if you set as staging you can access by tpot staging
- env: staging

  # proxy address example https://teleport.mine.com
  address: ""
  
  # proxy username example john.doe@mycomp.com or adzimzf
  # if you're using auth_connector it can be empty
  user_name: ""

  # if your proxy server using auth connector such as gsuite, facebook & okta
  auth_connector: ""

  # is your proxy server need two factor authentication
  two_fa: false

  # specified the tsh binary if your proxy has different tsh version
  # relative path is not supported yet
  # example /usr/bin/tsh-2
  # default it'll use your OS PATH
  tsh_path: ""

  # port forwarding configuration
  forwarding:
	interval: 30
	nodes:
      - host: teleport-127.0.0.1
        user_login: root
		local_port: 12345
		remote_port: 9000
		remote_host: "localhost"
`

type Proxy struct {
	Address  string `yaml:"address"        json:"address"`
	UserName string `yaml:"user_name"      json:"user_name"`
	Env      string `yaml:"env"            json:"env"`
	TwoFA    bool   `yaml:"two_fa"         json:"two_fa"`

	// For using OAUTH like GMAIL, Facebook etc
	// empty means using username & password
	AuthConnector string `yaml:"auth_connector" json:"auth_connector"`

	// TSHPath is the location of TSH binary
	// by default it'll use your PATH location
	TSHPath string `yaml:"tsh_path"       json:"tsh_path"`

	// Node contains the node information from teleport server
	Node Node `yaml:"node,omitempty" json:"node"`

	Forwarding Forwarding `yaml:"forwarding"`
}

// Validate validates the proxy configuration the node will be ignored
func (p *Proxy) Validate() error {
	_, err := url.ParseRequestURI(p.Address)
	if err != nil {
		return fmt.Errorf("address is invalid, error:%v", err)
	}

	if p.AuthConnector == "" && p.UserName == "" {
		return fmt.Errorf("auth_connector or user_name must not empty")
	}

	// TODO: need to support relative path such as ~/bin
	_, err = os.Stat(p.TSHPath)
	if err != nil && p.TSHPath != "" {
		return fmt.Errorf("tsh_path is invalid")
	}

	return nil
}

// ProxyStatus contains data about proxy status
type ProxyStatus struct {
	// LoginAs is the username logged
	LoginAs string `json:"login_as"`

	// Roles is a list of teleport role
	Roles []string `json:"roles"`

	// UserLogins is a list of user login
	UserLogins []string `json:"user_logins"`
}

type Node struct {
	Status *ProxyStatus `json:"status"`
	Items  []Item       `json:"items"`
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

// ListHostname return the list of hostname
func (n *Node) ListHostname() (res []string) {
	for _, n := range n.Items {
		res = append(res, n.Hostname)
	}
	return
}

type Item struct {
	Hostname string `json:"hostname"`
	Address  string `json:"addr"`
}

var ErrEnvNotFound = fmt.Errorf("env not found")

// AppendNode append the n to the proxy node list
func (p *Proxy) AppendNode(n Node) (Node, error) {
	pNode, err := p.GetNode()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
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
	nodeBytes, err := ioutil.ReadFile(Dir + "node_" + p.Env + ".json")
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
	return p.save(bytes)
}

func (p *Proxy) save(date []byte) error {
	return ioutil.WriteFile(Dir+"node_"+p.Env+".json", date, permission)
}

type Forwarding struct {
	Interval int               `yaml:"interval"`
	Nodes    []*ForwardingNode `yaml:"nodes"`
}

type ForwardingNode struct {
	Host       string `yaml:"host"`
	ListenPort string `yaml:"listen_port"`
	RemotePort string `yaml:"remote_port"`
	RemoteHost string `yaml:"remote_host"`
	UserLogin  string `yaml:"user_login"`
	Status     bool   `yaml:"-"`
	Error      string `yaml:"-"`
}

func (b *ForwardingNode) ViewName() string {
	return fmt.Sprintf("%s_%s_%s_%s", b.Host, b.ListenPort, b.RemotePort, b.RemoteHost)
}

func (b *ForwardingNode) Address() string {
	return fmt.Sprintf("%s:%s:%s", b.ListenPort, b.RemoteHost, b.RemotePort)
}
