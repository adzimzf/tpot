package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/adzimzf/tpot/editor"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
)

var (
	// Dir is the path where tpot store the configuration & cache
	// Dir will be overridden by flag -D
	Dir = os.Getenv("HOME") + "/.tpot/"

	// ErrValidateConfig is an error to indicate config is invalid
	ErrValidateConfig = errors.New("config is invalid")
)

// configFileName we'll only support YAML file
const configFileName = "config.yaml"

// Config is a config for tpot
type Config struct {

	// Editor is the editor to edit configuration
	// the default editor is nano
	Editor string `json:"editor"  yaml:"editor"`

	// Proxies is list of proxy configuration
	Proxies []*Proxy `json:"proxies" yaml:"proxies"`
}

// NewConfig load config from the file and create it if no exist
func NewConfig(isDev bool) (*Config, error) {
	if isDev {
		path, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		Dir = path + "/dev/"
	}
	if err := addConfigDirExist(); err != nil {
		return nil, err
	}
	config, err := getConfig()
	if errors.Is(err, os.ErrNotExist) {
		config = &Config{
			Editor: editor.DefaultEditor,
		}
		if err := config.save(); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return config, nil
}

func getConfig() (*Config, error) {
	bytes, err := ioutil.ReadFile(Dir + configFileName)
	if errors.Is(err, os.ErrNotExist) {
		// we're on migration from JSON file to YAML
		// to ensure backward compatibility we'll read JSON if exists
		// then convert the YAML file
		bytes, err := ioutil.ReadFile(Dir + "config.json")
		if err != nil {
			return nil, err
		}

		var config Config
		err = json.Unmarshal(bytes, &config)
		if err != nil {
			return nil, err
		}

		if err = config.save(); err != nil {
			return nil, err
		}

		return &config, nil
	}
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Add adds a new proxy configuration
func (c *Config) Add() (string, error) {
	return c.AddPlain(proxyTemplate)
}

// AddPlain adds a new proxy configuration by plain configuration
func (c *Config) AddPlain(configPlain string) (string, error) {
	result, err := editor.Edit(configPlain, "add_proxy*.yaml")
	if err != nil {
		return "", err
	}

	if result == proxyTemplate {
		return "", fmt.Errorf("there's no proxy was added")
	}

	var tmpConfig Config
	if err := yaml.Unmarshal([]byte(result), &tmpConfig); err != nil {
		return result, err
	}

	if l := len(tmpConfig.Proxies); l != 1 {
		return result, fmt.Errorf("need one proxy confugration, find %d", l)
	}

	if err := tmpConfig.Proxies[0].Validate(); err != nil {
		return result, fmt.Errorf("failed to validate %v", err)
	}

	// ensure the environment name is not exist
	_, err = c.FindProxy(tmpConfig.Proxies[0].Env)
	if err != ErrEnvNotFound {
		return result, fmt.Errorf("environment %s is already exist", tmpConfig.Proxies[0].Env)
	}

	c.Proxies = append(c.Proxies, tmpConfig.Proxies[0])
	return result, c.save()
}

// save saves the config into YAML file
func (c *Config) save() error {
	bytes, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(Dir+configFileName, bytes, permission)
}

// Edit edits a particular proxy with match env name
func (c *Config) Edit(envName string) (string, error) {
	proxy, err := c.FindProxy(envName)
	if err != nil {
		return "", fmt.Errorf("proxy %s is not found", envName)
	}

	tmpConfig := Config{
		Proxies: []*Proxy{
			proxy,
		},
	}
	marshal, err := yaml.Marshal(tmpConfig)
	if err != nil {
		return "", fmt.Errorf("proxy confugration is invalid")
	}
	return c.EditPlain(envName, string(marshal))
}

// EditPlain edit specific proxy configuration by config plain
func (c *Config) EditPlain(envName, configPlain string) (string, error) {
	result, err := editor.Edit(configPlain, "edit_proxy*.yaml")
	if err != nil {
		return "", err
	}

	var tmpConfig Config
	if err := yaml.Unmarshal([]byte(result), &tmpConfig); err != nil {
		return result, err
	}

	if l := len(tmpConfig.Proxies); l != 1 {
		return result, fmt.Errorf("need one proxy confugration, find %d", l)
	}
	newProxy, err := tmpConfig.FindProxy(envName)
	if err != nil {
		// If the proxy doesn't not found means user changes the proxy env name
		newProxy = tmpConfig.Proxies[0]
	}

	if err := newProxy.Validate(); err != nil {
		return result, fmt.Errorf("failed to validate %v", err)
	}

	for i, proxy := range c.Proxies {
		if proxy.Env == envName {
			c.Proxies[i] = newProxy
		}
	}
	return result, c.save()
}

// EditAll edit all the proxy configuration
func (c *Config) EditAll() (string, error) {
	marshal, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}
	plain, err := c.EditAllPlain(string(marshal))
	if err != nil {
		return plain, err
	}
	return plain, nil
}

func (c *Config) EditAllPlain(configPlain string) (string, error) {
	result, err := editor.Edit(configPlain, "add_proxy*.yaml")
	if err != nil {
		return "", err
	}
	marshal, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}

	if result == string(marshal) {
		return "", fmt.Errorf("there's no proxy was added")
	}

	var tmpConfig Config
	if err := yaml.Unmarshal([]byte(result), &tmpConfig); err != nil {
		return result, err
	}

	if l := len(tmpConfig.Proxies); l < 1 {
		return result, fmt.Errorf("need one proxy confugration, find %d", l)
	}

	var tmp2Config = Config{
		Editor: tmpConfig.Editor,
	}
	for _, proxy := range tmpConfig.Proxies {
		if err := proxy.Validate(); err != nil {
			return result, fmt.Errorf("failed to validate environment %s, error: %v", proxy.Env, err)
		}

		// ensure the environment name is not exist
		_, err = tmp2Config.FindProxy(proxy.Env)
		if err != ErrEnvNotFound {
			return result, fmt.Errorf("environment %s is already exist", tmpConfig.Proxies[0].Env)
		}

		tmp2Config.Proxies = append(tmp2Config.Proxies, proxy)
	}

	return result, tmp2Config.save()
}

// FindProxy finds the proxy by environment name
func (c *Config) FindProxy(env string) (*Proxy, error) {
	for _, p := range c.Proxies {
		if p.Env == env {
			return p, nil
		}
	}
	return nil, ErrEnvNotFound
}

// String return the beauty configuration
// maybe in the future we can colorize the result
func (c *Config) String() (string, error) {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func addConfigDirExist() error {
	_, err := os.Stat(Dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(Dir, os.ModePerm)
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
