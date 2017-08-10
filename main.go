package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	//"github.com/fredhsu/cvpgo/configlet"
	"github.com/fredhsu/cvpgo/inventory"
	"io/ioutil"
	"net/http"
)

type User struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

type AuthResp struct {
	UserName  string `json:"userName"`
	SessionId string `json:"sessionId"`
}

type AddInventory struct {
	Data []AddInventoryElement `json:"data"`
}

type AddInventoryElement struct {
	ContainerName string                 `json:"containerName"`
	ContainerId   string                 `json:"containerId"`
	ContainerType string                 `json:"containerType"`
	IpAddress     string                 `json:"ipAddress"`
	ContainerList []ContainerListElement `json:"containerList"`
}

type ContainerListElement struct {
	Name             string `json:"name"`
	Key              string `json:"key"`
	ChildContainerId string `json:"childContainerId"`
	Type             string `json:"type"`
}

func AddDevice(ipAddr string, cookies []*http.Cookie, client *http.Client) error {
	addDevice := AddInventoryElement{
		ContainerName: "Tenant",
		ContainerId:   "root",
		ContainerType: "Existing",
		IpAddress:     ipAddr,
		ContainerList: []ContainerListElement{},
	}
	addInventory := AddInventory{
		Data: []AddInventoryElement{addDevice},
	}
	jsonValue, err := json.Marshal(addInventory)
	addInventoryUrl := "/inventory/add/addToInventory.do?startIndex=0&endIndex=15"
	url := baseurl + addInventoryUrl
	fmt.Println(string(jsonValue))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	for _, c := range cookies {
		req.AddCookie(c)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("response Body:", string(body))
	return nil
}

var baseurl string

func main() {
	fmt.Println("vim-go")
	// Set connection to CVP
	cvpIp := "172.28.169.180"
	authUrl := "/login/authenticate.do"
	baseurl = "https://" + cvpIp + "/cvpservice"
	url := baseurl + authUrl
	fmt.Println(url)
	user := User{UserId: "cvpadmin", Password: "arista123"}
	jsonValue, err := json.Marshal(user)
	if err != nil {
		fmt.Println("error marshalling")
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonValue))
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
	fmt.Printf("Response JSON %s", authresp)
	cookies := resp.Cookies()

	inventory.AddDevice("10.10.10.2", cookies, client)

}
