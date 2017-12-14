package cvpgo

import (
	"testing"
)

func TestApplyConfiglet(t *testing.T) {
	cvpInfo := CVPInfo{IPAddress: "10.90.224.178", Username: "cvpadmin", Password: "arista123", Container: "CoreSite"}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	dev, _ := cvp.GetDevice("veos-35-166-252-96")
	t.Logf("Retrieved %+v", dev)
	ckl := []string{"configlet_944_2883464286319224", "configlet_950_2883521004044029"}
	cnl := []string{"veos-35-166-252-96-tunnel", "RECONCILE_35.166.252.96"}
	err := cvp.ApplyConfigletToDevice(dev.IPAddress, dev.Fqdn, dev.SystemMacAddress, cnl, ckl)
	if err != nil {
		t.Errorf("Erorr applying configlet : %s", err)
	}
}
