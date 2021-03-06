package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Release struct {
	TargetVersion string         `json:"targetVersion"`
	FromVersion   string         `json:"fromVersion"`
	MajorRelease  string         `json:"majorRelease"`
	Versionless   string         `json:"versionless"`
	CommandCheck  []CommandCheck `json:"commandCheck,omitempty"`
	Steps         []Steps        `json:"steps"`
}

type CommandCheck struct {
	Name   string `json:"name"`
	Action string `json:"action"`
}

type Steps struct {
	Order       string `json:"order"`
	Action      string `json:"action"`
	Description string `json:"description"`
	Command     string `json:"command"`
	Args        []Args `json:"args,omitempty"`
}

type Args struct {
	Option      string `json:"option"`
	Type        string `json:"type"`
	Nullable    string `json:"nullable"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

func main() {

	// Open our jsonFile
	jsonFile, err := os.Open("autoUpgrade.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened autoUpgrade.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var autoUpgradeJson []Release

	var autoUpgradeJsonMap []map[string]interface{}

	json.Unmarshal(byteValue, &autoUpgradeJson)

	json.Unmarshal(byteValue, &autoUpgradeJsonMap)

	fmt.Println(autoUpgradeJson)

	fmt.Println(autoUpgradeJsonMap[0]["targetVersion"])

}
