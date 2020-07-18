package config

import (
	"fmt"
	"net/url"
)

type Proxy struct {
	Address  string `json:"address"`
	UserName string `json:"user_name"`
	Env      string `json:"env"`
	TwoFA    bool   `json:"two_fa"`
	Node     Node
}

type ProxySetterStation interface {
	Execute(p *Proxy) error
}

type ProxySetter struct{}

func NewProxySetterStations() *ProxySetter {
	return &ProxySetter{}
}

func (ps *ProxySetter) Execute(p *Proxy) error {
	twoFAStt := NewSetTwoFAStation(nil)
	userNameStt := NewSetUserNameStation(twoFAStt)
	addrStt := NewSetAddressStation(userNameStt)
	envStt := NewSetEnvStation(addrStt)

	return envStt.Execute(p)
}

type SetEnvStation struct {
	next ProxySetterStation
}

func NewSetEnvStation(next ProxySetterStation) *SetEnvStation {
	return &SetEnvStation{next}
}

func (s *SetEnvStation) Execute(p *Proxy) error {
	var err error
	p.Env, err = prompt("Environment", func(env string) error {
		return nil
	})

	if err != nil {
		return err
	}

	return determineNext(s.next, p)
}

type SetAddressStation struct {
	next ProxySetterStation
}

func NewSetAddressStation(next ProxySetterStation) *SetAddressStation {
	return &SetAddressStation{next}
}

func (s *SetAddressStation) Execute(p *Proxy) error {
	var err error
	p.Address, err = prompt("Proxy Address (with http protocol)", func(address string) error {
		_, err := url.ParseRequestURI(address)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return determineNext(s.next, p)
}

type SetUserNameStation struct {
	next ProxySetterStation
}

func NewSetUserNameStation(next ProxySetterStation) *SetUserNameStation {
	return &SetUserNameStation{next}
}

func (s *SetUserNameStation) Execute(p *Proxy) error {
	var err error
	p.UserName, err = prompt("Username (teleport username)", func(userName string) error {
		return nil
	})

	if err != nil {
		return err
	}

	return determineNext(s.next, p)
}

type SetTowFAStation struct {
	next ProxySetterStation
}

func NewSetTwoFAStation(next ProxySetterStation) *SetTowFAStation {
	return &SetTowFAStation{next}
}

func (s *SetTowFAStation) Execute(p *Proxy) error {
	isTwoFA, err := prompt("Is Need 2FA (Y/y/N/n)", func(towFA string) error {
		if towFA == "Y" || towFA == "y" || towFA == "N" || towFA == "n" {
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

	return determineNext(s.next, p)
}

func determineNext(next ProxySetterStation, p *Proxy) error {
	if next != nil {
		return next.Execute(p)
	}

	return nil
}
