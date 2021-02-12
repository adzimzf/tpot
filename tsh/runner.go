package tsh

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/adzimzf/tpot/config"
)

type TSH struct {
	proxy              *config.Proxy
	userLogin, dstHost string
}

// tshBinary is the `tsh` binary where we depends
const tshBinary = "tsh"

// tshVersion is the supported tsh binary version
const tshVersion = "v4.1.11"

// SSH run the `tsh ssh` commands
func (t *TSH) SSH(username, host string) error {
	args, err := t.getProxyFlags()
	if err != nil {
		return err
	}

	args = append(args, t.authFlags()...)

	ipAddress, ok := t.proxy.Node.LookUpIPAddress(host)
	if !ok {
		return fmt.Errorf("couldn't find IP address")
	}

	args = append(args, username+"@"+ipAddress)

	cmd := exec.Command(t.tshBinary(), append([]string{"ssh"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ListNodes get the list nodes from proxy
func (t *TSH) ListNodes() (config.Node, error) {

	if err := t.login(); err != nil {
		return config.Node{}, err
	}

	args, err := t.getProxyFlags()
	if err != nil {
		return config.Node{}, err
	}

	cmd := exec.Command(t.tshBinary(), append([]string{"ls"}, args...)...)
	var stdOut, stdErr = &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout = stdOut
	cmd.Stdin = os.Stdin
	cmd.Stderr = stdErr
	if err = cmd.Run(); err != nil {
		return config.Node{}, err
	}
	if errStr := stdErr.String(); errStr != "" {
		return config.Node{}, errors.New(errStr)
	}

	return parseNodesFromString(stdOut.String()), nil
}

func (t *TSH) login() error {
	args, err := t.getProxyFlags()
	if err != nil {
		return err
	}

	args = append(args, t.authFlags()...)

	cmd := exec.Command(t.tshBinary(), append([]string{"login"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stdin
	return cmd.Run()
}

func (t *TSH) getProxyFlags() ([]string, error) {
	proxyAddress, err := t.cleanAddress()
	if err != nil {
		return nil, err
	}

	return []string{"--proxy=" + proxyAddress}, nil
}

// authFlags return the authentication flags
func (t *TSH) authFlags() []string {
	var args []string
	if t.proxy.AuthConnector != "" {
		args = append(args, "--auth="+t.proxy.AuthConnector)
	} else {
		args = append(args, "--user="+t.proxy.UserName)
	}
	return args
}

func (t *TSH) cleanAddress() (string, error) {
	u, err := url.Parse(t.proxy.Address)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}

// tshBinary return the location of TSH binary
func (t *TSH) tshBinary() string {
	if t.proxy.TSHPath != "" {
		return t.proxy.TSHPath
	}
	return tshBinary
}

func parseNodesFromString(nodeStr string) config.Node {
	var nodeList []config.Item
	for _, line := range strings.Split(nodeStr, "\n") {

		// remove the header of node table
		// for now on the data will get in table formatting,
		// to support all `tsh` old version
		// because the JSON formatting is only supported by
		// newer TSH
		if strings.HasPrefix(line, "Node") || strings.HasPrefix(line, "---") || strings.HasPrefix(line, " ") {
			continue
		}
		lines := strings.Split(line, " ")

		// infoCount indicate that the node information we want to get has already fulfill
		var infoCount int
		var node config.Item
		for _, s := range lines {
			if s == "" {
				continue
			}
			if infoCount == 2 {
				break
			}
			if infoCount == 0 {
				node.Hostname = s
			} else {
				node.Address = s
			}
			infoCount++
		}
		// doesn't need to append an empty node
		if node != (config.Item{}) {
			nodeList = append(nodeList, node)
		}
	}

	return config.Node{
		Items: nodeList,
	}
}

// NewTSH creates a new TSH
func NewTSH(p *config.Proxy) *TSH {
	return &TSH{
		proxy: p,
	}
}
