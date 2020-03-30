package cvpgo

import (
	"encoding/json"
	"fmt"
)

// ConfigForTask is a list of configurations attachedk to a task
type ConfigForTask struct {
	ReconciledConfig   string
	Reconcile          int
	New                int
	Configs            []Configlet
	DesignedConfig     []DesignedConfig
	Total              int
	RunningConfig      []DesignedConfig
	IsReconcileInvoked bool
	Mismatch           int
	Warning            []string
	Errors             string
	TargetIPAddress    string `json:"targetIpAddress"`
}

// DesignedConfig is a configuration to be added to or currently running on a device
type DesignedConfig struct {
	Command string `json:"command"`
	RowID   int    `json:"rowId"`
	Code    string `json:"code"`
	BlockID string `json:"blockId"`
}

// Overrides string output of designed config to remove color and formatting
func (dc DesignedConfig) String() string {
	return fmt.Sprintf("%s", dc.Command)
}

// GetConfigForTask takes a work order ID and returns the config assosciated with that task
func (c *CvpClient) GetConfigForTask(workOrderForID string) (ConfigForTask, error) {
	cft := ConfigForTask{}
	url := "/provisioning/getconfigfortask.do?workorderid=" + workOrderForID
	resp, err := c.Get(url)
	if err != nil {
		return cft, err
	}
	err = json.Unmarshal(resp, &cft)
	return cft, err
}

// GetDesignedConfig returns just the designed config without formatting as a string
func (cft *ConfigForTask) GetDesignedConfig() string {
	config := ""
	for _, dc := range cft.DesignedConfig {
		config += fmt.Sprintf("%s\n", dc.Command)
	}
	return config
}
