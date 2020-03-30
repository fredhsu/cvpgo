package cvpgo

import (
	"testing"
)

type CVPInfo struct {
	IPAddress string
	Username  string
	Password  string
	Container string
}

type InventoryTestData struct {
	CVP             *CVPInfo
	DeviceIP        string
	DeviceContainer string
	DeviceHostname  string
	DeviceMAC       string
}

func buildTestData() *InventoryTestData {
	result := InventoryTestData{}
	result.CVP = &CVPInfo{
		IPAddress: "localhost",
		Username:  "cvpadmin",
		Password:  "cvpadmin1",
		Container: "Test"}
	result.DeviceIP = "192.168.100.1"
	result.DeviceContainer = "Tenant"
	result.DeviceHostname = "Device-A"
	result.DeviceMAC = "02:42:ac:2a:8f:7d"
	return &result
}

func TestAddContainerToRoot(t *testing.T) {
	testdata := buildTestData()
	cvpInfo := *testdata.CVP
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	err := cvp.AddContainer(cvpInfo.Container, "Tenant")
	if err != nil {
		t.Errorf("%+v", err)
	}
}

func TestGetContainer(t *testing.T) {
	testdata := buildTestData()
	cvpInfo := *testdata.CVP
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	_, err := cvp.GetContainerByName(cvpInfo.Container)
	if err != nil {
		t.Errorf("%+v", err)
	}
}

func TestAddDevices(t *testing.T) {
	//cvpInfo := CVPInfo{IPAddress: "10.90.224.178", Username: "cvpadmin", Password: "arista123", Container: "CoreSite"}
	testdata := buildTestData()
	cvpInfo := *testdata.CVP
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	err := cvp.AddDevice(testdata.DeviceIP, testdata.DeviceContainer)
	if err != nil {
		t.Errorf("%+v", err)
	}
}

func TestSearchTempInventory(t *testing.T) {
	//cvpInfo := CVPInfo{IPAddress: "10.90.224.178", Username: "cvpadmin", Password: "arista123", Container: "CoreSite"}
	testdata := buildTestData()
	cvpInfo := *testdata.CVP
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	dev, err := cvp.SearchInventory(testdata.DeviceIP)
	if err != nil {
		t.Errorf("%+v", err)
	}
	if dev == nil {
		t.Errorf("Did not retreive any devices\n")
	}
	t.Logf("Retrieved container status is %s", dev.Status)
}

func TestSaveCommit(t *testing.T) {
	//cvpInfo := CVPInfo{IPAddress: "10.90.224.178", Username: "cvpadmin", Password: "arista123", Container: "CoreSite"}
	testdata := buildTestData()
	cvpInfo := *testdata.CVP
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	err := cvp.SaveCommit(testdata.DeviceIP, 5)
	if err != nil {
		t.Errorf("%+v", err)
	} else {
		t.Logf("Device safely saved")
	}
}

func TestSaveInventory(t *testing.T) {
	//cvpInfo := CVPInfo{IPAddress: "10.90.224.178", Username: "cvpadmin", Password: "arista123", Container: "CoreSite"}
	testdata := buildTestData()
	cvpInfo := *testdata.CVP
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	err := cvp.SaveInventory()
	if err != nil {
		t.Errorf("%+v", err)
	}
}

func TestGetDevices(t *testing.T) {
	//cvpInfo := CVPInfo{IPAddress: "10.90.224.178", Username: "cvpadmin", Password: "arista123", Container: "CoreSite"}
	testdata := buildTestData()
	cvpInfo := *testdata.CVP
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	dev, err := cvp.GetDevice(testdata.DeviceHostname)
	if err != nil {
		t.Errorf("%+v", err)
	}
	if dev == nil {
		t.Errorf("Did not retreive any devices\n")
	}
	t.Logf("Retrieved %+v", dev)
}

func TestRemoveDevices(t *testing.T) {
	//cvpInfo := CVPInfo{IPAddress: "10.90.224.178", Username: "cvpadmin", Password: "arista123", Container: "CoreSite"}
	testdata := buildTestData()
	cvpInfo := *testdata.CVP
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	err := cvp.RemoveDevice(testdata.DeviceMAC)
	if err != nil {
		t.Errorf("%+v", err)
	}
}

func TestDeleteContainer(t *testing.T) {
	testdata := buildTestData()
	cvpInfo := *testdata.CVP
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	err := cvp.DeleteContainer(cvpInfo.Container, "Tenant")
	if err != nil {
		t.Errorf("%+v", err)
	}
}
