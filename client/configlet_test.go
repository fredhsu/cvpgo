package cvpgo

import (
	"testing"
)

type TestData struct {
	CvpIP        string
	CvpUser      string
	CvpPwd       string
	CvpContainer string
	Device       string
	Cnls         []string
	NewConfig    string
	NewCfglet    string
}

func buildConfigletTestData() *TestData {
	return &TestData{
		CvpIP:        "localhost",
		CvpUser:      "cvpadmin",
		CvpPwd:       "cvpadmin1",
		CvpContainer: "Tenant",
		Cnls:         []string{"Test1"},
		NewConfig:    "username TESTCONFIG nopassword",
		NewCfglet:    "Test1",
		Device:       "Device-A",
	}
}

func TestAddConfiglet(t *testing.T) {
	data := buildConfigletTestData()
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	configlet := Configlet{
		Name:   data.NewCfglet,
		Config: data.NewConfig,
	}
	_, err := cvp.AddConfiglet(configlet)
	if err != nil {
		t.Errorf("Error adding configlet : %s", err)
	}
}

func TestGetConfiglet(t *testing.T) {
	data := buildConfigletTestData()
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	cfglet, err := cvp.GetConfigletByName(data.NewCfglet)
	if err != nil {
		t.Errorf("Error getting configlet : %s", err)
	}
	if cfglet.Key == "" {
		t.Logf("No Configlets were found")
	} else {
		t.Logf("Returned configlet with key %s", cfglet.Key)
	}
}

func TestApplyConfiglet(t *testing.T) {
	data := buildConfigletTestData()
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	dev, err := cvp.GetDevice(data.Device)
	//t.Logf("device found %+v", dev)
	if err != nil {
		t.Errorf("Erorr getting device info : %s", err)
	}
	cflts, err := cvp.getConfigletsByName(data.Cnls)
	if err != nil {
		t.Errorf("Erorr getting configlets from device : %s", err)
	}
	//t.Logf("Configlets found %+v", cflts)
	cnls := getNames(cflts)
	ckls := getKeys(cflts)
	t.Logf("Configlet Names:%s, Keys: %s", cnls, ckls)
	err = cvp.ApplyConfigletToDevice(dev.IPAddress, dev.Fqdn, dev.SystemMacAddress, cnls, ckls)
	if err != nil {
		t.Errorf("Erorr applying configlet : %s", err)
	}
}

func TestValidateConfig(t *testing.T) {
	data := buildConfigletTestData()
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	device, err := cvp.GetDevice(data.Device)
	if err != nil {
		t.Errorf("Error getting device : %s", data.Device)
	}
	err = cvp.ValidateConfig(device.Key, data.NewConfig)
	if err != nil {
		t.Errorf("Error getting configlet : %s", err)
	}
}

func TestGetConfigletByDeviceID(t *testing.T) {
	data := buildConfigletTestData()
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	device, err := cvp.GetDevice(data.Device)
	_, err = cvp.GetConfigletByDeviceID(device.Key)
	if err != nil {
		t.Errorf("Error getting configlet : %s", err)
	}
}

func TestRemoveConfiglet(t *testing.T) {
	data := buildConfigletTestData()
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	dev, _ := cvp.GetDevice(data.Device)
	t.Logf("Retrieved %+v", dev)
	err := cvp.RemoveConfigletFromDevice(dev.IPAddress, dev.Fqdn, dev.Key, []string{data.NewCfglet}, true)
	if err != nil {
		t.Errorf("Error applying configlet : %s", err)
	}
}

func TestDeleteConfiglet(t *testing.T) {
	data := buildConfigletTestData()
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	err := cvp.DeleteConfiglet(data.NewCfglet)
	if err != nil {
		t.Errorf("Erorr adding configlet : %s", err)
	}
}
