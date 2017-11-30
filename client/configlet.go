package cvpgo

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

type Configlet struct {
	Config string `json:"config"`
	Name   string `json:"name"`
}

func (c *CvpClient) AddConfiglet(configlet Configlet) error {
	addConfigletURL := "/configlet/addConfiglet.do"
	_, err := c.Call(configlet, addConfigletURL)
	return err
}
