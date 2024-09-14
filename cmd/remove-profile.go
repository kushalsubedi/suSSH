package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var RemoveProfileCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove an SSH profile",
	Long:  "Remove an SSH profile from the list of profiles.",
	Run: func(cmd *cobra.Command, args []string) {
		RemoveProfile()
	},
}

func RemoveProfile() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error retrieving home directory:", err)
		return
	}

	filePath := filepath.Join(homeDir, ".ssh", "ssh-config.json")

	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading profiles file:", err)
		return
	}

	var profiles []Profile
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
		Message: "Choose a profile to remove:",
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
		}
	}

	var newProfiles []Profile
	for _, profile := range profiles {
		if profile != selectedProfile {
			newProfiles = append(newProfiles, profile)
		}
	}

	profilesJSON, err := json.MarshalIndent(newProfiles, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling profiles to JSON:", err)
		return

	}

	if err := os.WriteFile(filePath, profilesJSON, 0600); err != nil {
		fmt.Println("Error writing profiles to config file:", err)
		return
	}

	fmt.Println("✅ Profile removed successfully!")

}

func init() {
	rootCmd.AddCommand(RemoveProfileCmd)
}
