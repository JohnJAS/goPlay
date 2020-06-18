package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	cdfCommon "autoUpgrade/common"
)

type Release struct {
	TargetVersion string         `json:"targetVersion"`
	FromVersion   string         `json:"fromVersion"`
	MajorRelease  string         `json:"majorRelease"`
	Versionless   string         `json:"versionless"`
	CommandCheck  []CommandCheck `json:"commandCheck,omitempty"`
	Steps         []Step         `json:"steps"`
}

type CommandCheck struct {
	Name   string `json:"name"`
	Action string `json:"action"`
}

type Step struct {
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

func GetUpgradeChain(path string) (result []string, err error) {
	var autoUpgradeJson []Release
	var data []byte

	data, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &autoUpgradeJson)
	if err != nil {
		return
	}

	for _, member := range autoUpgradeJson {
		result = append(result, member.TargetVersion)
	}

	return
}

func GetIfMajor(path string, targeVersion string) (isMajor bool, err error) {
	var autoUpgradeJson []Release
	var data []byte

	data, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &autoUpgradeJson)
	if err != nil {
		return
	}

	for _, member := range autoUpgradeJson {
		if member.TargetVersion == targeVersion {
			if member.MajorRelease == "true" {
				return true, nil
			} else {
				return false, nil
			}
		}
	}

	return false, errors.New("fail to find in json file")
}

func GetIfVersionless(path string, targeVersion string) (isMajor bool, err error) {
	var autoUpgradeJson []Release
	var data []byte

	data, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &autoUpgradeJson)
	if err != nil {
		return
	}

	for _, member := range autoUpgradeJson {
		if member.TargetVersion == targeVersion {
			if member.Versionless == "true" {
				return true, nil
			} else {
				return false, nil
			}
		}
	}

	return false, errors.New("fail to find in json file")
}

func GetAutoUpgradeJsonObj(path string) (autoUpgradeJson []Release, err error) {
	var data []byte
	data, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &autoUpgradeJson)
	if err != nil {
		return
	}

	return
}

func GetReleaseJsonObj(autoUpgradeJsonObj []Release, version string) (release Release, err error) {
	for _, release := range autoUpgradeJsonObj {
		if release.TargetVersion == version {
			return release, nil
		}
	}
	return
}

func GetStepObj(steps []Step, order string) (stepObj Step, err error) {
	for _, step := range steps {
		if step.Order == order {
			return step, nil
		}
	}
	return
}

func Test() {

	// Open our jsonFile
	jsonFile, err := os.Open(cdfCommon.AutoUpgradeJSON)
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
