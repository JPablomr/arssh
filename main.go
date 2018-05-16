package main

import (
	"os"
)

func main() {
	// Get what we can from the command-line
	sshInfo := parseArgs()
	_, forceRefresh := os.LookupEnv("CACHEBUST")
	allInstances := getInstances(forceRefresh)
	setupSSH(sshInfo, allInstances)
}

func setupSSH(sshInfo sshInfo, instanceList []*InstanceData) {

	selectedInstance := selectInstance(instanceList, sshInfo.host)
	if selectedInstance == nil {
		os.Exit(0)
	}
	if sshInfo.user == "" {
		sshInfo.user = getDefaultUser(selectedInstance.Os)
	}
	sshInfo.ip = selectedInstance.PrivateIP

	execSSH(sshInfo)
}
