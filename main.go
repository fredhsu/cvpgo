package main

import (
	cvpgo "github.com/fredhsu/cvpgo/client"
)

func main() {
	// Set connection to CVP
	cvpIP := "10.90.224.178"
	cvp := cvpgo.New(cvpIP, "cvpadmin", "arista123")
	cvp.AddDevice("10.10.10.2")

	// configlet.AddConfiglet(newconfiglet, cookies, client)
}
