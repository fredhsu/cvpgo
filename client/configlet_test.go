package cvpgo

import (
	"os"
	"strings"
	"testing"
)

type TestData struct {
	CvpIP        string
	CvpUser      string
	CvpPwd       string
	CvpContainer string
	Device       string
	Ckls         []string
	Cnls         []string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvList(key string, fallback []string) []string {
	if value, ok := os.LookupEnv(key); ok {
		return strings.Split(value, ",")
	}
	return fallback
}

func updateFromEnvironment(data *TestData) {
	data.CvpIP = getEnv("CVP_IP", data.CvpIP)
	data.CvpUser = getEnv("CVP_USER", data.CvpUser)
	data.CvpPwd = getEnv("CVP_PWD", data.CvpPwd)
	data.CvpContainer = getEnv("CVP_CONT", data.CvpContainer)
	data.Device = getEnv("CVP_DEVICE", data.Device)
	data.Ckls = getEnvList("CVP_CKLS", data.Ckls)
	data.Cnls = getEnvList("CVP_CNLS", data.Cnls)
}

func TestApplyConfiglet(t *testing.T) {
	data := TestData{
		CvpIP:        "10.90.224.178",
		CvpUser:      "cvpadmin",
		CvpPwd:       "arista123",
		CvpContainer: "CoreSite",
		Device:       "veos-35-166-252-96",
		Ckls:         []string{"configlet_944_2883464286319224", "configlet_950_2883521004044029"},
		Cnls:         []string{"veos-35-166-252-96-tunnel", "RECONCILE_35.166.252.96"},
	}
	updateFromEnvironment(&data)
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	dev, _ := cvp.GetDevice(data.Device)
	t.Logf("Retrieved %+v", dev)
	ckl := data.Ckls
	cnl := data.Cnls
	err := cvp.ApplyConfigletToDevice(dev.IPAddress, dev.Fqdn, dev.SystemMacAddress, cnl, ckl)
	if err != nil {
		t.Errorf("Erorr applying configlet : %s", err)
	}
}
