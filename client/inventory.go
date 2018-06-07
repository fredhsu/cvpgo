package cvpgo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

type GetInventory struct {
	Total          int          `json:"total"`
	NetElementList []NetElement `json:"netElementList"`
}

type GetTempInventory struct {
	Total              int              `json:"total"`
	TempNetElementList []TempNetElement `json:"tempNetElement"`
}

type GetContainer struct {
	Total         int         `json:"total"`
	ContainerList []Container `json:"containerList"`
}

type Container struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type NetElement struct {
	ModeName         string `json:"modelName"`
	InternalVersion  string `json:"internalVersion"`
	SystemMacAddress string `json:"systemMacAddress"`
	SerialNumber     string `json:"serialNumber"`
	Version          string `json:"version"`
	Fqdn             string `json:"fqdn"`
	Key              string `json:"key"`
	IPAddress        string `json:"ipAddress"`
}

type TempNetElement struct {
	ModeName         string `json:"modelName"`
	InternalVersion  string `json:"internalVersion"`
	SystemMacAddress string `json:"systemMacAddress"`
	SerialNumber     string `json:"serialNumber"`
	Version          string `json:"version"`
	Fqdn             string `json:"fqdn"`
	Key              string `json:"key"`
	IPAddress        string `json:"ipAddress"`
	Status           string `json:"status"`
}

type InventoryData struct {
	Total      int `json:"total"`
	Containers struct {
		HierarchyNetElementCount int    `json:"hierarchyNetElementCount"`
		ChildNetElementCount     int    `json:"childNetElementCount"`
		Key                      string `json:"key"`
		Name                     string `json:"name"`
		Type                     string `json:"type"`
		ChildContainerCount      int    `json:"childContainerCount"`
		ParentContainerID        string `json:"parentContainerId"`
		Mode                     string `json:"mode"`
		DeviceStatus             string `json:"deviceStatus"`
		ChildTaskCount           int    `json:"childTaskCount"`
		ChildNetElementList      string `json:"childNetElementList"`
		TempAction               string `json:"tempAction"`
		TempEvent                string `json:"tempEvent"`
		ChildContainerList       []struct {
		} `json:"childContainerList"`
	} `json:"containers"`
	TempContainer []struct {
		Name             string `json:"name"`
		ParentID         string `json:"parentId"`
		Key              string `json:"key"`
		ChildContainerID string `json:"childContainerId"`
		Type             string `json:"type"`
		UserID           string `json:"userId"`
		FactoryID        int    `json:"factoryId"`
		ID               int    `json:"id"`
	} `json:"tempContainer"`
	Dashboard []struct {
		Count     int    `json:"count"`
		Status    string `json:"status"`
		FactoryID int    `json:"factoryId"`
		ID        int    `json:"id"`
	} `json:"dashboard"`
	TempNetElement []struct {
		ClassID          string   `json:"classId"`
		ModelName        string   `json:"modelName"`
		InternalVersion  string   `json:"internalVersion"`
		SystemMacAddress string   `json:"systemMacAddress"`
		MemTotal         string   `json:"memTotal"`
		BootupTimeStamp  string   `json:"bootupTimeStamp"`
		MemFree          string   `json:"memFree"`
		Architecture     string   `json:"architecture"`
		InternalBuildID  string   `json:"internalBuildId"`
		HardwareRevision string   `json:"hardwareRevision"`
		Fqdn             string   `json:"fqdn"`
		IPAddress        string   `json:"ipAddress"`
		TaskIDList       []string `json:"taskIdList"`
		FactoryID        string   `json:"factoryId"`
		ZtpMode          string   `json:"ztpMode"`
		IsDANZEnabled    string   `json:"isDANZEnabled"`
		ContainerID      string   `json:"containerId"`
		StatusMessage    string   `json:"statusMessage"`
		ContainerName    string   `json:"containerName"`
		UserID           string   `json:"userId"`
		Version          string   `json:"version"`
		Key              string   `json:"key"`
		ID               string   `json:"id"`
		Type             string   `json:"type"`
		SerialNumber     string   `json:"serialNumber"`
		Status           string   `json:"status"`
		IsMLAGEnabled    string   `json:"isMLAGEnabled"`
		UserName         string   `json:"userName"`
		Passowrd         string   `json:"passowrd"`
	} `json:"tempNetElement"`
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

// AddDevice adds a device into CVP's inventory
func (c *CvpClient) AddDevice(ipAddr string, cn string) error {
	container, err := c.GetContainerByName(cn)
	if err != nil {
		return err
	}
	element := []AddInventoryElement{
		AddInventoryElement{
			ContainerName: cn,
			ContainerId:   container.Key,
			ContainerType: "Existing",
			IpAddress:     ipAddr,
			ContainerList: []ContainerListElement{},
		},
	}
	addInventory := AddInventory{
		Data: element,
	}
	// In case the device is already in temp Inventory, we're purging it from there
	if _, err := c.SearchInventory(ipAddr); err != nil {
		c.CancelTempInventory()
	}

	addInventoryURL := "/inventory/add/addToInventory.do?startIndex=0&endIndex=15"
	_, err = c.Call(addInventory, addInventoryURL)
	return err
}

// SaveInventory saves all connected devices from CVP's temp into normal inventory
func (c *CvpClient) SaveInventory() error {
	saveInventoryURL := "/inventory/v2/saveInventory.do"
	_, err := c.Call("", saveInventoryURL)
	return err
}

// SearchInventory searches for a device CVP's temp inventory
func (c *CvpClient) SearchInventory(ip string) (*TempNetElement, error) {
	searchInventoryURL := "/inventory/add/searchInventory.do?queryparam=" + ip + "&startIndex=0&endIndex=0"
	respbody, err := c.Get(searchInventoryURL)
	respDevice := GetTempInventory{}
	err = json.Unmarshal(respbody, &respDevice)
	if err != nil {
		log.Printf("Error decoding getdevice :%s\n", err)
		return nil, err
	}
	if len(respDevice.TempNetElementList) == 0 {
		return nil, fmt.Errorf("No devices returned")
	}
	return &respDevice.TempNetElementList[0], err
}

// SaveCommit tries to save a device into CVP's inventory if it's connected
func (c *CvpClient) SaveCommit(ip string, seconds int) error {
	timeout := time.After(time.Duration(seconds) * time.Second)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("Device is still not connected")
		default:
		}
		dev, err := c.SearchInventory(ip)
		if err != nil {
			return err
		}
		if dev.Status == "Connected" {
			return c.SaveInventory()
		}
		time.Sleep(1 * time.Second)
		continue
	}
}

// Removes device from CVP's inventory
func (c *CvpClient) RemoveDevice(mac string) error {
	data := struct {
		Data []string `json:"data"`
	}{[]string{mac}}
	RemoveInventoryURL := "/inventory/deleteDevices.do?"
	_, err := c.Call(data, RemoveInventoryURL)
	return err
}

// Cancels temp device inventory
func (c *CvpClient) CancelTempInventory() error {
	CancelInventoryURL := "/inventory/add/cancelInventory.do"
	_, err := c.Get(CancelInventoryURL)
	return err
}

// GetDevice uses the unique ID of a device to lookup the full entry in CVP
// and returns the full NetElement entry
func (c *CvpClient) GetDevice(id string) (*NetElement, error) {
	getDeviceURL := "/inventory/getInventory.do?queryparam=" + url.QueryEscape(id) + "&startIndex=0&endIndex=0"
	respbody, err := c.Get(getDeviceURL)
	respDevice := GetInventory{}
	err = json.Unmarshal(respbody, &respDevice)
	if err != nil {
		log.Printf("Error decoding getdevice :%s\n", err)
		return nil, err
	}
	if len(respDevice.NetElementList) == 0 {
		return nil, fmt.Errorf("No devices returned")
	}
	return &respDevice.NetElementList[0], err
}

// Returns a container that exactly matches the name.
func (c *CvpClient) GetContainerByName(query string) (*Container, error) {
	getContainerURL := "/provisioning/searchTopology.do?queryParam=" + query + "&startIndex=0&endIndex=0"
	respbody, err := c.Get(getContainerURL)
	respContainer := GetContainer{}
	err = json.Unmarshal(respbody, &respContainer)
	if err != nil {
		log.Printf("Error decoding getcontainer :%s\n", err)
		return nil, err
	}
	if len(respContainer.ContainerList) == 0 {
		return nil, fmt.Errorf("No container named \"%s\" found", query)
	}
	return &respContainer.ContainerList[0], err
}

// Returns a container name based on its ID
func (c *CvpClient) GetContainerNameById(query string) (string, error) {
	url := "/provisioning/getContainerInfoById.do?containerId=" + query
	respbody, err := c.Get(url)
	respContainer := struct {
		Name string `json:"name"`
	}{}
	err = json.Unmarshal(respbody, &respContainer)
	if err != nil {
		log.Printf("Error decoding getcontainerbyid :%s\n", err)
		return "", err
	}
	return respContainer.Name, nil
}

// GetInventory will return all the devices in CVP
func (c *CvpClient) GetInventory(query string) (*[]NetElement, error) {
	getDeviceURL := "/inventory/getInventory.do?queryparam=" + query + "&startIndex=0&endIndex=0"
	respbody, err := c.Get(getDeviceURL)
	respDevice := GetInventory{}
	err = json.Unmarshal(respbody, &respDevice)
	if err != nil {
		log.Printf("Error decoding getdevice :%s\n", err)
		return nil, err
	}
	if len(respDevice.NetElementList) == 0 {
		return nil, fmt.Errorf("No devices returned")
	}
	return &respDevice.NetElementList, err
}

func (c *CvpClient) AddContainerToRoot(new string) error {
	name, err := c.GetContainerNameById("root")
	if err != nil {
		log.Printf("Could not find 'root' container")
		return err
	}
	return c.AddContainer(new, name)
}

func (c *CvpClient) AddContainer(new, parent string) error {
	newC := &Container{
		Name: new,
	}
	parentC, err := c.GetContainerByName(parent)
	if err != nil {
		log.Printf("Parent container %s cannot be found", parent)
		return err
	}

	_, err = c.containerOp(newC, parentC, "add")
	if err != nil {
		log.Printf("Error applying configlet : %s", err)
		return err
	}

	return nil
}

func (c *CvpClient) DeleteContainer(name, parent string) error {
	currentC, err := c.GetContainerByName(name)
	if err != nil {
		log.Printf("Container %s cannot be found", name)
		return err
	}
	parentC, err := c.GetContainerByName(parent)
	if err != nil {
		log.Printf("Parent container %s cannot be found", parent)
		return err
	}
	_, err = c.containerOp(currentC, parentC, "delete")
	if err != nil {
		log.Printf("Error applying configlet : %s", err)
		return err
	}

	return nil
}

func (c *CvpClient) containerOp(container, parent *Container, op string) (sdata SaveData, err error) {
	info := "Performing " + op + " operation on container " + container.Name
	data := Action{
		Info:        info,
		InfoPreview: info,
		Action:      op,
		NodeType:    "container",
		NodeID:      container.Key,
		NodeName:    container.Name,
		ToIDType:    "container",
	}

	if op == "add" {
		data.ToID = parent.Key
		data.ToName = parent.Name
		data.NodeID = "new_container"
	} else if op == "delete" {
		data.FromID = parent.Key
		data.FromName = parent.Name
	}
	log.Printf("Operation data read: %+v", data)
	if err = c.addTempAction(data); err != nil {
		return sdata, err
	}

	return c.saveTopologyV2([]string{})
}
