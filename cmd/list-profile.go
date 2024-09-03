package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// ListProfileCmd represents the list-profile command
var listProfileCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all SSH profiles",
	Long:  "List all SSH profiles and allow selection using arrow keys. Selecting a profile will initiate an SSH login.",
	Run: func(cmd *cobra.Command, args []string) {
		ListProfiles()
	},
}

func ListProfiles() {
	// Retrieve the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error retrieving home directory:", err)
		return
	}

	// Define the file path for the profiles
	filePath := filepath.Join(homeDir, ".ssh", "ssh-config.json")

	// Read the profiles file
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading profiles file:", err)
		return
	}

	// Parse the profiles
	var profiles []Profile // Reusing the Profile struct from add.go
	err = json.Unmarshal(file, &profiles)
	if err != nil {
		fmt.Println("Error parsing profiles:", err)
		return
	}

	// If no profiles are found
	if len(profiles) == 0 {
		fmt.Println("No profiles found.")
		return
	}

	// Prepare options for selection
	options := make([]string, len(profiles))
	for i, profile := range profiles {
		options[i] = fmt.Sprintf("HostName: %s, HostIP: %s, KeyPath: %s", profile.HostName, profile.HostIP, profile.KeyPath)
	}

	// Prompt user to select a profile
	var selectedOption string
	prompt := &survey.Select{
		Message: "Choose a profile to login:",
		Options: options,
	}

	err = survey.AskOne(prompt, &selectedOption)
	if err != nil {
		fmt.Println("Error selecting profile:", err)
		return
	}

	// Find the selected profile
	var selectedProfile Profile
	for _, profile := range profiles {
		if fmt.Sprintf("HostName: %s, HostIP: %s, KeyPath: %s", profile.HostName, profile.HostIP, profile.KeyPath) == selectedOption {
			selectedProfile = profile
			break
		}
	}

	// Trigger login process
	LoginToInstance(selectedProfile)
}

func LoginToInstance(profile Profile) {
	// Construct SSH command
	permissionCMD := exec.Command("chmod", "400", profile.KeyPath)
	sshCmd := exec.Command("ssh", "-i", profile.KeyPath, fmt.Sprintf("ubuntu@%s", profile.HostIP))

	// add loading with emoji
	fmt.Println("🚀 Logging in to", profile.HostName)

	// Set the command to use the user's current terminal
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr
	sshCmd.Stdin = os.Stdin

	// Execute the command
	err := sshCmd.Run()
	if err != nil {
		fmt.Println("Failed to login:", err)
	}
	err = permissionCMD.Run()
	if err != nil {
		fmt.Println("Failed to change permission:", err)
	}
}

func init() {
	rootCmd.AddCommand(listProfileCmd)
}
