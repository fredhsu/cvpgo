package main

import (
	cvpgo "github.com/fredhsu/cvpgo/client"
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

func main() {
	// Set connection to CVP
	cvpIP := "10.90.224.178"
	cvp := cvpgo.New(cvpIP, "cvpadmin", "arista123")
	cvp.AddDevice("10.10.10.2")

	// configlet.AddConfiglet(newconfiglet, cookies, client)
}
