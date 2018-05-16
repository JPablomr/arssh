package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

// This will work just fine for the time being.
// TODO: Add default in case that env var is missing
var dir = os.Getenv("HOME") + "/.arssh"
var filename = dir + "/" + cacheFilename() + ".json"

// Fetches instances from cache
// Sometimes refreshes the cache
func getInstances(forceRefresh bool) []*InstanceData {
	createCacheFolder()
	rand.Seed(time.Now().UnixNano())
	randRefresh := (rand.Intn(100) > 80)

	if _, err := os.Stat(filename); os.IsNotExist(err) || forceRefresh || randRefresh {
		fmt.Println("Refreshing cache...")
		instanceData := getInstanceData()
		writeInstanceCache(instanceData)
		return instanceData
	}
	instanceData := readInstanceCache()
	return instanceData
}

func createCacheFolder() {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0750)
		if err != nil {
			fmt.Println("Could not create dir! ", err)
		}
	}
}

func cacheFilename() string {
	profileName, ok := os.LookupEnv("AWS_PROFILE")
	if !ok {
		return "default"
	}
	return profileName
}

func writeInstanceCache(instances []*InstanceData) {

	marshalled, err := json.Marshal(instances)
	if err != nil {
		fmt.Println("Error creating JSON", err)
		return
	}
	ioutil.WriteFile(filename, marshalled, 0644)
}

func readInstanceCache() []*InstanceData {
	data, _ := ioutil.ReadFile(filename)
	var result []*InstanceData
	json.Unmarshal(data, &result)
	return result
}
