package main

import (
	"fmt"

	cvpgo "github.com/fredhsu/cvpgo/client"
)

func main() {
	// Set connection to CVP
	cvpIP := "10.90.224.175"
	cvp := cvpgo.New(cvpIP, "cvpadmin", "arista")
	//configlets, err := cvp.Get("/configlet/getConfiglets.do?startIndex=0&endIndex=0")
	// if err != nil {
	// 	fmt.Println("error!")
	// 	fmt.Println(err)
	// 	return
	// }
	// m := cvpgo.ConfigletList{}
	//err = json.Unmarshal(configlets, &m)
	//fmt.Println("Configets are : ", m)

	cft, err := cvp.GetConfigForTask("101")
	// fmt.Println(cft)
	// fmt.Println("ConfigForTask ", cft)
	// for _, dc := range cft.DesignedConfig {
	// 	fmt.Println(dc)
	// }

	fmt.Println(cft.GetDesignedConfig())
	//fmt.Println(configlets)
	//cvp.AddDevice("10.10.10.2", "CoreSite")
	// configlet.AddConfiglet(newconfiglet, cookies, client)
	if err != nil {
		fmt.Println("error!")
		fmt.Println(err)
		return
	}
}
