package tsh

import (
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/adzimzf/tpot/config"
)

type TSH struct {
	proxy   *config.Proxy
	dstHost string
}

func (t *TSH) Run() error {

	addr, err := t.cleanAddress()
	if err != nil {
		return err
	}

	cmd := exec.Command("tsh", "ssh",
		"--proxy="+addr,
		"--user="+t.proxy.UserName,
		"root@"+t.getIP(),
	)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (t *TSH) cleanAddress() (string, error) {
	u, err := url.Parse(t.proxy.Address)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}

func (t *TSH) getIP() string {
	for _, i2 := range t.proxy.Node.Items {
		if i2.Hostname == t.dstHost {
			return strings.Split(i2.Hostname, ":")[0]
		}
	}
	return ""
}

func NewTSH(p *config.Proxy, dstHost string) *TSH {
	return &TSH{
		proxy:   p,
		dstHost: dstHost,
	}
}
