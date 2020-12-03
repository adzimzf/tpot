package scapper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/adzimzf/tpot/config"
	"github.com/manifoldco/promptui"
)

type Scrapper struct {
	proxy  config.Proxy
	client http.Client
}

func NewScrapper(p config.Proxy) *Scrapper {
	return &Scrapper{
		proxy:  p,
		client: http.Client{Timeout: 30 * time.Second},
	}
}

func (s *Scrapper) getCSRF() (string, error) {
	resp, err := s.client.Get(s.proxy.Address + "/web/login")
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// it could simply solve by using regex
	for _, line := range strings.Split(string(respByte), "\n") {
		if strings.Contains(line, "grv_csrf_token") {
			r := regexp.MustCompile(".*content=\"")
			s1 := r.ReplaceAllString(line, "")
			r2 := regexp.MustCompile("\".*")
			return r2.ReplaceAllString(s1, ""), nil
		}
	}
	return "", fmt.Errorf("csrf not found")
}

func (s *Scrapper) GetNodes() (config.Node, error) {

	request, err := http.NewRequest(http.MethodGet, s.proxy.Address+"/v1/webapi/sites/main/nodes", nil)
	if err != nil {
		return config.Node{}, err
	}

	jwtToken, cookie, err := s.getJWTToken()
	if err != nil {
		return config.Node{}, err
	}
	request.Header.Add("Cookie", cookie)
	request.Header.Add("Authorization", "Bearer "+jwtToken)
	resp, err := s.client.Do(request)
	if err != nil {
		return config.Node{}, err
	}

	defer resp.Body.Close()
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return config.Node{}, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(string(respByte))
		return config.Node{}, fmt.Errorf("http error code: %d", resp.StatusCode)
	}

	var nodes config.Node
	err = json.Unmarshal(respByte, &nodes)
	if err != nil {
		return config.Node{}, err
	}
	return nodes, nil
}

func (s *Scrapper) getJWTToken() (string, string, error) {
	csrf, err := s.getCSRF()
	if err != nil {
		return "", "", err
	}

	req, err := s.getLoginReq(csrf)
	if err != nil {
		return "", "", err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("http code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	var res map[string]interface{}
	err = json.Unmarshal(respByte, &res)
	if err != nil {
		return "", "", err
	}
	token, ok := res["token"].(string)
	if !ok {
		return "", "", fmt.Errorf("failed to convert token to string")
	}
	return token, resp.Header.Get("Set-Cookie"), nil
}

func (s *Scrapper) getPassAndFactor() (string, string, error) {

	pass, err := s.prompt("Password", '*')
	if err != nil {
		return "", "", err
	}

	if !s.proxy.TwoFA {
		return pass, "", nil
	}

	twoFA, err := s.prompt("2FA Token", rune(0))
	if err != nil {
		return "", "", err
	}
	return pass, twoFA, nil

}

func (s *Scrapper) prompt(label string, mask rune) (string, error) {

	prompt := promptui.Prompt{
		Label: label,
		Mask:  mask,
	}

	return prompt.Run()
}

func (s *Scrapper) getLoginReq(csrfToken string) (*http.Request, error) {

	pass, token, err := s.getPassAndFactor()
	if err != nil {
		return nil, err
	}

	body := bytes.NewBufferString(fmt.Sprintf(`{"user":"%s","pass":"%s","second_factor_token":"%s"}`,
		s.proxy.UserName, pass, token))
	request, err := http.NewRequest("POST", s.proxy.Address+"/v1/webapi/sessions", bufio.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-CSRF-Token", csrfToken)
	request.Header.Add("Accept", "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,"+
		"application/signed-exchange;v=b3;q=0.9")
	request.Header.Add("Cookie", "grv_csrf="+csrfToken)
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")

	return request, nil
}
