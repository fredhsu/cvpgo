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
	IpAddress string
	BaseURL   string
	Client    *http.Client
	Cookies   []*http.Cookie
}

type User struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

type AuthResp struct {
	UserName  string `json:"userName"`
	SessionId string `json:"sessionId"`
}

// New creates a new CVP Client to host
func New(host string, user string, password string) CvpClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
	c := CvpClient{IpAddress: host, BaseURL: "https://" + host + "/cvpservice", Client: client}
	c.authenticate(user, password)
	return c
}

func (c *CvpClient) authenticate(username string, password string) {
	authURL := "/login/authenticate.do"
	url := c.BaseURL + authURL
	user := User{UserId: username, Password: password}
	jsonValue, err := json.Marshal(user)
	if err != nil {
		fmt.Println("error marshalling")
	}
	resp, err := c.Client.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	authresp := AuthResp{}
	json.Unmarshal(body, &authresp)
	c.Cookies = resp.Cookies()
}

func (c *CvpClient) Call(obj interface{}, svcurl string) ([]byte, error) {
	jsonValue, err := json.Marshal(obj)
	url := c.BaseURL + svcurl
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	for _, c := range c.Cookies {
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
