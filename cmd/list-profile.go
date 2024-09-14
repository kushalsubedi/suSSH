package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var listProfileCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all SSH profiles",
	Long:  "List all SSH profiles and allow selection using arrow keys. Selecting a profile will initiate an SSH login.",
	Run: func(cmd *cobra.Command, args []string) {
		ListProfiles()
	},
}

func ListProfiles() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error retrieving home directory:", err)
		return
	}

	filePath := filepath.Join(homeDir, ".ssh", "ssh-config.json")

	// Read the profiles file
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading profiles file:", err)
		return
	}

	var profiles []Profile // Reusing the Profile struct from add.go
	err = json.Unmarshal(file, &profiles)
	if err != nil {
		fmt.Println("Error parsing profiles:", err)
		return
	}

	if len(profiles) == 0 {
		fmt.Println("No profiles found.")
		return
	}

	options := make([]string, len(profiles))
	for i, profile := range profiles {
		options[i] = fmt.Sprintf("HostName: %s, HostIP: %s, KeyPath: %s", profile.HostName, profile.HostIP, profile.KeyPath)
	}

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

	var selectedProfile Profile
	for _, profile := range profiles {
		if fmt.Sprintf("HostName: %s, HostIP: %s, KeyPath: %s", profile.HostName, profile.HostIP, profile.KeyPath) == selectedOption {
			selectedProfile = profile
			break
		}
	}

	LoginToInstance(selectedProfile)
}

func LoginToInstance(profile Profile) {
	permissionCmd := exec.Command("chmod", "400", profile.KeyPath)
	fmt.Println("🔑 Changing file permission...")

	if err := permissionCmd.Run(); err != nil {
		fmt.Println("Error Changing in permission", err)
		return
	}
	time.Sleep(1 * time.Second)
	sshCmd := exec.Command("ssh", "-i", profile.KeyPath, fmt.Sprintf("ubuntu@%s", profile.HostIP))

	fmt.Println("🚀 Logging in to", profile.HostName)

	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr
	sshCmd.Stdin = os.Stdin

	if err := sshCmd.Run(); err != nil {
		fmt.Println("Login  TO Instance Succes")
	}
}
func init() {
	rootCmd.AddCommand(listProfileCmd)
}
