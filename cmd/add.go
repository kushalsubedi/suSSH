package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Profile struct {
	HostName string
	HostIP   string
	KeyPath  string
}

func AddProfile(hostName string, hostIP string, keyPath string) {
	// Check if the file exists
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}
	filePath := filepath.Join(homeDir, ".ssh", "ssh-config.json")
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		// Create the file
		_, err = os.Create(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	// Read the file
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Parse the file
	var profiles []Profile
	err = json.Unmarshal(file, &profiles)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Check if the profile already exists
	for _, profile := range profiles {
		if profile.HostName == hostName {
			fmt.Println("Profile already exists")
			return
		}
	}
	// Add the new profile
	profiles = append(profiles, Profile{HostName: hostName, HostIP: hostIP, KeyPath: keyPath})
	// Write the new file
	profilesJSON, err := json.Marshal(profiles)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.WriteFile(filePath, profilesJSON, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Profile added")
}
