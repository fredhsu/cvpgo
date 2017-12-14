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

func TestGetDevices(t *testing.T) {
	cvpInfo := CVPInfo{IPAddress: "10.90.224.178", Username: "cvpadmin", Password: "arista123", Container: "CoreSite"}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	dev, err := cvp.GetDevice("veos-35-167-79-72")
	if err != nil {
		t.Errorf("%+v", err)
	}
	if dev == nil {
		t.Errorf("Did not retreive any devices\n")
	}
	t.Logf("Retrieved %+v", dev)
}
