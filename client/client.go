package cvpgo

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// CvpClient provides a client to a CVP Host
type CvpClient struct {
	BaseURL string
	Client  *http.Client
	Cookies []*http.Cookie
}

// User defines a CVP user
type User struct {
	UserID   string `json:"userId"`
	Password string `json:"password"`
}

type authResp struct {
	UserName  string `json:"userName"`
	SessionID string `json:"sessionId"`
}

// New creates a new CVP Client pointing to host
func New(host string, user string, password string) CvpClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	c := CvpClient{BaseURL: "https://" + host + "/cvpservice", Client: client}
	c.authenticate(user, password)
	return c
}

func (c *CvpClient) authenticate(username string, password string) {
	authURL := "/login/authenticate.do"
	url := c.BaseURL + authURL
	user := User{UserID: username, Password: password}
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
	authresp := authResp{}
	json.Unmarshal(body, &authresp)
	c.Cookies = resp.Cookies()
}

// Call issues a POST to the svcurl with a JSON encoded obj
func (c *CvpClient) Call(obj interface{}, svcurl string) ([]byte, error) {
	jsonValue, err := json.Marshal(obj)
	url := c.BaseURL + svcurl
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	for _, c := range c.Cookies {
		req.AddCookie(c)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.Client.Do(req)
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
