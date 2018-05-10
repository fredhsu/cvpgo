package cvpgo

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type JsonData struct {
	Data         interface{} `json:"data,omitempty"`
	ErrorCode    string      `json:"errorCode"`
	ErrorMessage string      `json:"errorMessage"`
}

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
	Key    string `json:"key,omitempty"`
}

type DeleteConfiglet struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type ConfigletList struct {
	List []Configlet `json:"configletList"`
}

type ValidateRequest struct {
	NetElementID string   `json:"netElementId"`
	ConfigIDList []string `json:"configIdList"`
	PageType     string   `json:"pageType"`
}

type ValidateConfigRequest struct {
	NetElementID string `json:"netElementId"`
	Config       string `json:"config"`
}

type ValidateConfigResponse struct {
	WarningCount string `json:"warningCount"`
	ErrorCount   string `json:"errorCount"`
}

func checkErrors(data JsonData) error {
	if data.ErrorCode != "" {
		log.Printf("Error from CVP: %s", data.ErrorMessage)
		return fmt.Errorf("CVP returned error code: %s, %s", data.ErrorCode, data.ErrorMessage)
	}
	return nil
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
// NOTE: despite what swagger says, this does not generate reconcile configlet
func (c *CvpClient) ValidateConfiglet(netElementID string, cfgletIDList []string) error {
	url := "/provisioning/v2/validateAndCompareConfiglet.do"
	req := ValidateRequest{
		NetElementID: netElementID,
		ConfigIDList: cfgletIDList,
	}
	_, err := c.Call(req, url)
	return err
}

func (c *CvpClient) ValidateConfig(netElementID, config string) error {
	url := "/configlet/validateConfig.do"
	req := ValidateConfigRequest{
		NetElementID: netElementID,
		Config:       config,
	}
	resp, err := c.Call(req, url)
	body := ValidateConfigResponse{}
	err = json.Unmarshal(resp, &body)
	if err != nil {
		log.Printf("Error validating config %+v", err)
	}
	if _, err := strconv.Atoi(body.ErrorCount); err != nil {
		return fmt.Errorf("Config validation produced errors")
	}
	return nil
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
	cfgletCurrent, err := c.GetConfigletByDeviceID(deviceMac)
	if err != nil {
		log.Printf("Error retrieving configlets from a device")
		return err
	}
	cfgletNew, err := c.getConfigletsByName(cnl)
	if err != nil {
		log.Printf("Error retrieving configlets by its name")
		return err
	}
	cfgletAll, _ := c.mergeCfglet(cfgletCurrent, cfgletNew)
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
		ConfigletList:                   getKeys(cfgletAll),
		ConfigletNamesList:              getNames(cfgletAll),
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

// GetConfigletByDeviceID gets list of configlets assigned to a device
func (c *CvpClient) GetConfigletByDeviceID(deviceMac string) ([]Configlet, error) {
	url := "/provisioning/getConfigletsByNetElementId.do?netElementId=" + deviceMac + "&queryParam=&startIndex=0&endIndex=15"
	respbody, err := c.Get(url)
	respConfiglet := ConfigletList{}
	err = json.Unmarshal(respbody, &respConfiglet)
	if err != nil {
		log.Printf("Error decoding GetConfigletByDeviceID :%s\n", err)
		return nil, err
	}
	return respConfiglet.List, err
}

func (c *CvpClient) GetConfigletByName(cfglet string) (Configlet, error) {
	url := "/configlet/getConfigletByName.do?name=" + cfglet
	respbody, err := c.Get(url)
	respConfiglet := Configlet{}
	err = json.Unmarshal(respbody, &respConfiglet)
	if err != nil {
		log.Printf("Error decoding getConfigletByName :%s\n", err)
		return respConfiglet, err
	}
	return respConfiglet, err
}

func (c *CvpClient) getConfigletsByName(cfglets []string) (result []Configlet, err error) {
	for _, cfgletName := range cfglets {
		cfglet, err := c.GetConfigletByName(cfgletName)
		if err != nil {
			return result, err
		}
		result = append(result, cfglet)
	}
	return result, nil
}

func contains(list []Configlet, elem Configlet) bool {
	for _, t := range list {
		if t == elem {
			return true
		}
	}
	return false
}

func (c *CvpClient) filterCfglet(all, remove []Configlet) (stay []Configlet, err error) {
	for _, cfglet := range all {
		if !contains(remove, cfglet) {
			stay = append(stay, cfglet)
		}
	}
	return stay, nil
}

func (c *CvpClient) mergeCfglet(current, new []Configlet) (all []Configlet, err error) {
	all = append(current, all...)
	for _, cfglet := range new {
		if !contains(current, cfglet) {
			all = append(all, cfglet)
		}
	}
	return all, nil
}

func getNames(cfglets []Configlet) []string {
	result := make([]string, 0)
	for _, cfglet := range cfglets {
		result = append(result, cfglet.Name)
	}
	return result
}

func getKeys(cfglets []Configlet) []string {
	result := make([]string, 0)
	for _, cfglet := range cfglets {
		result = append(result, cfglet.Key)
	}
	return result
}

// RemoveConfigletFromDevice removes configlets (list of strings) from device
func (c *CvpClient) RemoveConfigletFromDevice(deviceIP, deviceName, deviceMac string, cfgletRemoveNames []string, save bool) error {
	cfgletAll, err := c.GetConfigletByDeviceID(deviceMac)
	if err != nil {
		return err
	}
	cfgletRemove, err := c.getConfigletsByName(cfgletRemoveNames)
	if err != nil {
		return err
	}
	cfgletRemain, err := c.filterCfglet(cfgletAll, cfgletRemove)
	removeCfglet := ApplyConfiglet{
		Info:                            "Configlet Remove from device: " + deviceName,
		InfoPreview:                     "<b>Configlet remove</b> from Device " + deviceName,
		Action:                          "associate",
		NodeIPAddress:                   deviceIP,
		NodeTargetIPAddress:             deviceIP,
		NodeType:                        "configlet",
		ToID:                            deviceMac,
		ToIDType:                        "netelement",
		ToName:                          deviceName,
		ConfigletList:                   getKeys(cfgletRemain),
		ConfigletNamesList:              getNames(cfgletRemain),
		ConfigletBuilderList:            []string{},
		ConfigletBuilderNamesList:       []string{},
		IgnoreConfigletList:             getKeys(cfgletRemove),
		IgnoreConfigletNamesList:        getNames(cfgletRemove),
		IgnoreConfigletBuilderList:      []string{},
		IgnoreConfigletBuilderNamesList: []string{},
	}
	log.Printf("Removing configlet : %+v", removeCfglet)
	err = c.addTempAction(removeCfglet)
	if save {
		return c.saveTopologyV2(removeCfglet)
	}
	return err
}

func (c *CvpClient) saveTopologyV2(action ApplyConfiglet) error {
	url := "/provisioning/v2/saveTopology.do"
	dataArray := []ApplyConfiglet{action}
	data := ApplyConfigletData{
		Data: dataArray,
	}
	resp, err := c.Call(data, url)
	log.Printf("Response from temp action %+s", resp)
	return err
}

// DeleteConfiglet deletes configlet from CVP
func (c *CvpClient) DeleteConfiglet(cfgletName string) error {
	url := "/configlet/deleteConfiglet.do"
	cfglet, err := c.GetConfigletByName(cfgletName)
	// properties which are not allowed by the schema: ["config"]
	delCfgl := DeleteConfiglet{
		Key:  cfglet.Key,
		Name: cfgletName,
	}
	body := []DeleteConfiglet{delCfgl}
	resp, err := c.Call(body, url)
	responseBody := JsonData{}
	if err = json.Unmarshal(resp, &responseBody); err != nil {
		log.Printf("Error adding configlet %+v", err)
	}
	return checkErrors(responseBody)
}
