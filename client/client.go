package cvpgo

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var baseurl string
var client *http.Client
var cookies []*http.Cookie

type CvpClient struct {
	Baseurl string
	client  *http.Client
	Cookies []*http.Cookie
}

func (c *CvpClient) NewCvpClient() {
}

func (c *CvpClient) SetHost(host string) {
	c.Baseurl = "https://" + host + "/cvpservice"
}

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
}

func SetHost(host string) {
	baseurl = "https://" + host + "/cvpservice"
}

func Authenticate() {
}

func Call(obj interface{}, svcurl string, cookies []*http.Cookie) ([]byte, error) {
	jsonValue, err := json.Marshal(obj)
	url := baseurl + svcurl
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	for _, c := range cookies {
		req.AddCookie(c)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("response Body:", string(body))
	return body, nil
}
