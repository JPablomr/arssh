package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type sshInfo struct {
	user string
	host string
	ip   string
	args []string
}

func main() {
	_, forceRefresh := os.LookupEnv("CACHEBUST")
	allInstances := getInstances(forceRefresh)
	setupSSH(allInstances)
}

func setupSSH(instanceList []*InstanceData) {
	// Get what we can from the command-line
	sshInfo := parseArgs()
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

func execSSH(sshInfo sshInfo) {
	sshBin, err := exec.LookPath("ssh")
	if err != nil {
		panic(err)
	}
	env := os.Environ()
	hostString := fmt.Sprintf("%s@%s", sshInfo.user, sshInfo.ip)
	sshArgs := []string{sshBin, hostString}
	sshArgs = append(sshArgs, sshInfo.args...)
	err = syscall.Exec(sshBin, sshArgs, env)
	if err != nil {
		panic(err)
	}
}

// Initial Parsing of the args
func parseArgs() sshInfo {
	hostString := os.Args[1]
	result := sshInfo{}
	// Here we know the user and the host substring
	if strings.Contains(hostString, "@") {
		splitData := strings.Split(hostString, "@")
		result.user = splitData[0]
		result.host = splitData[1]
	} else { // We just know the host substring
		result.host = hostString
	}
	result.args = os.Args[2:]

	return result
}

func prettyPrint(instances []*InstanceData) {
	for counter, instance := range instances {
		fmt.Printf("%d) %s - %s - %s (%s) - %s\n",
			counter+1,
			instance.Name,
			instance.InstanceID,
			instance.PrivateIP,
			instance.Az,
			instance.LaunchTime.Format(time.RFC3339Nano),
		)
	}
}

func instanceSearch(instances []*InstanceData, searchTerm string) []*InstanceData {
	// Avoid allocations with this one little trick!
	// This works because this empty slice shares the same array as instances
	results := instances[:0]
	for _, instance := range instances {
		if strings.Contains(instance.Name, searchTerm) {
			results = append(results, instance)
		}
	}
	return results
}

// Select Instance from results
func selectInstance(instances []*InstanceData, searchTerm string) *InstanceData {
	candidates := instanceSearch(instances, searchTerm)
	reader := bufio.NewReader(os.Stdin)
	if len(candidates) == 1 {
		return candidates[0]
	}
	fmt.Println("Matching Instances:")
	prettyPrint(candidates)
	fmt.Print("Select an instance number to ssh to it or anything else to cancel: ")
	rawInput, _ := reader.ReadString('\n')
	option, err := strconv.Atoi(rawInput[:len(rawInput)-1])
	if err != nil || option < 1 {
		fmt.Println("Invalid input! Exiting.", err)
		return nil
	}
	return candidates[option-1]
}
