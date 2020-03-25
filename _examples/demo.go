package main

import (
	"log"

	swc "github.com/randlabs/server-watchdog-go"
)

func main() {
	swdClient, err := swc.Create(swc.ClientOptions{
		Host: "127.0.0.1",
		Port: 3004,
		ApiKey: "set-some-key",
		DefaultChannel: "default",
	})
	if err != nil {
		log.Fatalf("Error creating client [%v]\n", err.Error())
	}

	log.Printf("Sending an error message thru channel '%v'...\n", swdClient.GetDefaultChannel())
	err = swdClient.Error("This is a sample error message from the Server Watchdog test application", "")
	if err != nil {
		log.Fatalf("Error sending message [%v]\n", err.Error())
	}

	//--------

	log.Printf("Sending a warning message thru channel '%v'...\n", swdClient.GetDefaultChannel())
	err = swdClient.Warn("This is a sample warning message from the Server Watchdog test application", "")
	if err != nil {
		log.Fatalf("Error sending message [%v]\n", err.Error())
	}

	//--------

	log.Printf("Sending an information message thru channel '%v'...\n", swdClient.GetDefaultChannel())
	err = swdClient.Info("This is a sample information message from the Server Watchdog test application", "")
	if err != nil {
		log.Fatalf("Error sending message [%v]\n", err.Error())
	}

	//--------

	log.Printf("Start monitoring ourselves thru channel '%v'...\n", swdClient.GetDefaultChannel())
	err = swdClient.ProcessWatch(0, "Server Watcher Go SDK Demo", "", "")
	if err != nil {
		log.Fatalf("Error sending message [%v]\n", err.Error())
	}

	//--------

	log.Printf("Stop monitoring ourselves on channel '%v'...\n", swdClient.GetDefaultChannel())
	err = swdClient.ProcessUnwatch(0, "")
	if err != nil {
		log.Fatalf("Error sending message [%v]\n", err.Error())
	}

	log.Printf("Done!\n")
}
