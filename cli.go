// Handle all the CLI interaction:

// - Instance selection
// - exec into SSH

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

// Initial Parsing of the args
func parseArgs() sshInfo {
	if len(os.Args) == 1 {
		usage()
	}
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

func usage() {
	fmt.Println(`Usage:
	arssh [user@][Instance Name or tag] [optional SSH commands]
	`)
	os.Exit(0)
}

// run SSH with the host info
func execSSH(sshInfo sshInfo) {
	sshBin, err := exec.LookPath("ssh")
	if err != nil {
		panic(err)
	}
	env := os.Environ()
	hostString := fmt.Sprintf("%s@%s", sshInfo.user, sshInfo.ip)
	sshArgs := []string{sshBin, hostString}
	// Append is a variadic function, by adding ... to the slice it'll insert its contents
	sshArgs = append(sshArgs, sshInfo.args...)
	err = syscall.Exec(sshBin, sshArgs, env)
	if err != nil {
		panic(err)
	}
}

// Print the instances in a list
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

// Search instances by ID and Name Tag
func instanceSearch(instances []*InstanceData, searchTerm string) []*InstanceData {
	// Avoid allocations with this one little trick!
	// This works because this empty slice shares the same array as instances
	results := instances[:0]
	for _, instance := range instances {
		// Check for name tag
		if strings.Contains(instance.Name, searchTerm) || instance.InstanceID == searchTerm {
			results = append(results, instance)
		}
	}
	return results
}
