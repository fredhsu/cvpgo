package cvpgo

import (
	"encoding/json"
	"fmt"
	"log"
)

type AddConfigletData struct {
	Data         AddConfigletRepsonse `json:"data"`
	ErrorCode    string               `json:"errorCode"`
	ErrorMessage string               `json:"errorMessage"`
}

type AddConfigletRepsonse struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Config string `json:"config"`
	User   string `json:"user"`
}

type AddConfigletResp struct {
	Data []struct {
		Key                  string `json:"key"`
		Name                 string `json:"name"`
		Reconciled           bool   `json:"reconciled"`
		Config               string `json:"config"`
		User                 string `json:"user"`
		Note                 string `json:"note"`
		ContainerCount       int    `json:"containerCount"`
		NetElementCount      int    `json:"netElementCount"`
		DateTimeInLongFormat int    `json:"dateTimeInLongFormat"`
		IsDefault            string `json:"isDefault"`
		IsAutoBuilder        string `json:"isAutoBuilder"`
		Type                 string `json:"type"`
		FactoryID            int    `json:"factoryId"`
		ID                   int    `json:"id"`
	} `json:"data"`
}

type ApplyConfigletData struct {
	Data []ApplyConfiglet `json:"data"`
}

type ApplyConfiglet struct {
	Info                            string   `json:"info"`
	InfoPreview                     string   `json:"infoPreview"`
	Action                          string   `json:"action"`
	NodeType                        string   `json:"nodeType"`
	NodeID                          string   `json:"nodeId"`
	ToID                            string   `json:"toId"`
	ToIDType                        string   `json:"toIdType"`
	FromID                          string   `json:"fromId"`
	NodeName                        string   `json:"nodeName"`
	FromName                        string   `json:"fromName"`
	ToName                          string   `json:"toName"`
	NodeIPAddress                   string   `json:"nodeIpAddress"`
	NodeTargetIPAddress             string   `json:"nodeTargetIpAddress"`
	ConfigletList                   []string `json:"configletList"`
	ConfigletNamesList              []string `json:"configletNamesList"`
	IgnoreConfigletList             []string `json:"ignoreConfigletList"`
	IgnoreConfigletNamesList        []string `json:"ignoreConfigletNamesList"`
	ConfigletBuilderList            []string `json:"configletBuilderList"`
	ConfigletBuilderNamesList       []string `json:"configletBuilderNamesList"`
	IgnoreConfigletBuilderList      []string `json:"ignoreConfigletBuilderList"`
	IgnoreConfigletBuilderNamesList []string `json:"ignoreConfigletBuilderNamesList"`
}

type Configlet struct {
	Config string `json:"config"`
	Name   string `json:"name"`
}

type ValidateRequest struct {
	NetElementID string   `json:"netElementId"`
	ConfigIDList []string `json:"configIdList"`
	PageType     string   `json:"pageType"`
}

func (c *CvpClient) AddConfiglet(configlet Configlet) (AddConfigletData, error) {
	addConfigletURL := "/configlet/addConfiglet.do"
	resp, err := c.Call(configlet, addConfigletURL)
	body := AddConfigletData{}
	err = json.Unmarshal(resp, &body)
	if err != nil {
		log.Printf("Error adding configlet %+v", err)
	}
	if body.ErrorCode != "" {
		log.Printf("Error from CVP: %s", body.ErrorMessage)
		return body, fmt.Errorf("CVP returned error code: %s, %s", body.ErrorCode, body.ErrorMessage)
	}
	return body, err
}

// ValidateConfiglet takes the netElementId (MAC Address) and a Configlet ID
// given as Key from adding a configlet, and validates it
func (c *CvpClient) ValidateConfiglet(netElementID, cfgletID string) error {
	url := "/ztp/v2/validateAndCompareConfiglet.do"
	req := ValidateRequest{
		NetElementID: netElementID,
		ConfigIDList: []string{cfgletID},
	}
	_, err := c.Call(req, url)
	return err
}

/* ApplyConfigletToDevice applies configlets to a device
   deviceIpAddress -- Ip address of the device (type: string)
   deviceFqdn -- Fully qualified domain name for device (type: string)
   deviceKey -- mac address of the device (type: string)
   cnl -- List of name of configlets to be applied
   (type: List of Strings)
   ckl -- Keys of configlets to be applied (type: List of Strings)
*/
func (c *CvpClient) ApplyConfigletToDevice(deviceIP, deviceName, deviceMac string, cnl, ckl []string) error {
	// func (c *CvpClient) ApplyConfigletToDevice(deviceName, deviceMac string, cnl, ckl, []string) error {
	// func (c *CvpClient) ApplyConfigletToDevice(deviceIP, deviceName, deviceMac string, cnl, ckl, cbnl, cbkl []string) error {
	applyCfglet := ApplyConfiglet{
		Info:                            "Configlet Assign to device: " + deviceName,
		InfoPreview:                     "<b>Configlet assign</b> to Device " + deviceName,
		Action:                          "associate",
		NodeIPAddress:                   deviceIP,
		NodeTargetIPAddress:             deviceIP,
		NodeType:                        "configlet",
		ToID:                            deviceMac,
		ToIDType:                        "netelement",
		ToName:                          deviceName,
		ConfigletList:                   ckl,
		ConfigletNamesList:              cnl,
		ConfigletBuilderList:            []string{},
		ConfigletBuilderNamesList:       []string{},
		IgnoreConfigletList:             []string{},
		IgnoreConfigletNamesList:        []string{},
		IgnoreConfigletBuilderList:      []string{},
		IgnoreConfigletBuilderNamesList: []string{},
	}
	log.Printf("Applying configlet : %+v", applyCfglet)
	return c.addTempAction(applyCfglet)
}

func (c *CvpClient) addTempAction(action ApplyConfiglet) error {
	url := "/ztp/addTempAction.do?format=topology&queryParam=&nodeId=root"
	dataArray := []ApplyConfiglet{action}
	data := ApplyConfigletData{
		Data: dataArray,
	}
	resp, err := c.Call(data, url)
	log.Printf("Response from temp action %+s", resp)
	return err
}
