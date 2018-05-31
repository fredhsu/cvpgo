package cvpgo

import (
	"log"
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
	NetElementID string
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
		NetElementID: "02:42:ac:2a:8f:7d",
	}
}

func TestValidateCompare(t *testing.T) {
	data := buildConfigletTestData()
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	id := data.NetElementID
	_, err := cvp.ValidateCompareCfglt(id, []string{})
	if err != nil {
		t.Errorf("Error adding configlet : %s", err)
	}
	//t.Logf("Return response is %+v", resp)
}

func TestGenerateReconcile(t *testing.T) {
	data := buildConfigletTestData()
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	id := data.NetElementID
	resp, err := cvp.ValidateCompareCfglt(id, []string{})
	if err != nil {
		log.Printf("Error generating reconcile configlet %+v", err)
	}
	configlet := Configlet{
		Name:   resp.ReconciledConfig.Name,
		Config: resp.ReconciledConfig.Config,
	}
	err = cvp.UpdateReconcile(id, configlet)
	if err != nil {
		t.Errorf("Error adding configlet : %s", err)
	}
}

func TestPushReconcile(t *testing.T) {
	data := buildConfigletTestData()
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	device, _ := cvp.GetDevice(data.Device)
	reconcileName := "RECONCILE_" + device.IPAddress
	dev, _ := cvp.GetDevice(data.Device)
	t.Logf("Retrieved %+v", dev)
	cnl := []string{reconcileName}
	sdata, err := cvp.ApplyConfigletToDevice(dev.IPAddress, dev.Fqdn, dev.SystemMacAddress, cnl, true)
	if err != nil {
		t.Errorf("Error applying configlet : %s", err)
	}
	taskIds := sdata.Data.TaskIds
	if err = cvp.ExecuteTasks(taskIds); err != nil {
		t.Errorf("Error executing tasks : %s", err)
	}
	if err = cvp.CheckTasks(taskIds, 10); err != nil {
		t.Errorf("Some tasks are not completed : %s", err)
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
	dev, _ := cvp.GetDevice(data.Device)
	t.Logf("Retrieved %+v", dev)
	cnl := data.Cnls
	_, err := cvp.ApplyConfigletToDevice(dev.IPAddress, dev.Fqdn, dev.SystemMacAddress, cnl, true)
	if err != nil {
		t.Errorf("Error applying configlet : %s", err)
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
	_, err := cvp.RemoveConfigletFromDevice(dev.IPAddress, dev.Fqdn, dev.Key, []string{data.NewCfglet}, true)
	if err != nil {
		t.Errorf("Error removing configlet : %s", err)
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
