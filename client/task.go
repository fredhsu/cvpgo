package cvpgo

import (
	"encoding/json"
)

type TasksSuccessResponse struct {
	Total int
	Data  []Task
}

// Task is a config change task in CVP
type Task struct {
	WorkOrderID string `json:"workOrderId"`
	Name        string
	//workOrderEscalation (object),
	WorkOrderState             string
	CurrentTaskName            string
	TemplateID                 string `json:"templateId"`
	Data                       TaskData
	CreatedBy                  string
	ExecutedOnInLongFormat     int
	ExecutedBy                 string
	CurrentTaskType            string
	WorkFlowDetailsID          string `json:"workFlowDetailsId"`
	Description                string
	WorkOrderUserDefinedStatus string
	WorkOrderDetails           WorkOrderDetails
	CreatedOnInLongFormat      int
	CompletedOnInLongFormat    int
	Note                       string
	TaskStatus                 string
	TaskStatusBeforeCancel     string
	NewParentContainerName     string
	FactoryID                  int    `json:"factoryId"`
	NewParentContainerID       string `json:"newParentContainerId`
	CcID                       string `json:"ccId`
	ID                         int    `json:"id"`
}

// TaskData provides additional data on the task
type TaskData struct {
	IS_CONFIG_PUSH_NEEDED    string
	CurrentparentContainerId string
	WORKFLOW_ACTION          string
	NewparentContainerId     string
	NETELEMENT_ID            string
	IS_ADD_OR_MOVE_FLOW      bool
	VIEW                     string
	ImageIdList              string
	Imagebundle              string
	ImageBundleId            string
	APP_SESSION_ID           string
}
type WorkOrderDetails struct {
	WorkOrderDetailsId string
	NetElementHostName string
	NetElementId       string
	IpAddress          string
	WorkOrderId        string
	SerialNumber       string
	FactoryId          int
	Id                 int
}

type WorkOrderLog struct {
	TaskId  string
	Message string
	Source  string
}

// Overrides string output of designed config to remove color and formatting
// func (dc DesignedConfig) String() string {
// 	return fmt.Sprintf("%s", dc.Command)
// }

// GetTasks fetches all the current tasks in CV
func (c *CvpClient) GetTasks(query string) (TasksSuccessResponse, error) {
	tsr := TasksSuccessResponse{}
	queryParam := ""
	if query != "" {
		queryParam = "queryparam=" + query + "&"
	}
	url := "/task/getTasks.do?" + queryParam + "startIndex=0&endIndex=0"
	resp, err := c.Get(url)
	if err != nil {
		return tsr, err
	}
	err = json.Unmarshal(resp, &tsr)
	return tsr, err
}

func (c *CvpClient) AddWorkOrderLog(log WorkOrderLog) error {
	_, err := c.Call(log, "/workflow/addWorkOrderLog.do")
	return err
}
