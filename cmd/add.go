package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// Profile struct represents an SSH profile
type Profile struct {
	HostName string `json:"hostname"`
	HostIP   string `json:"hostip"`
	KeyPath  string `json:"keypath"`
}

func AddProfile(profile Profile) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	filePath := filepath.Join(homeDir, ".ssh", "ssh-config.json")

	var profiles []Profile

	// Check if the file exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, read and unmarshal existing profiles
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read existing config file: %v", err)
		}
		if len(fileContent) > 0 {
			if err := json.Unmarshal(fileContent, &profiles); err != nil {
				return fmt.Errorf("failed to parse existing config file: %v", err)
			}
		}
	}

	// Check for duplicate hostname
	for _, existingProfile := range profiles {
		if existingProfile.HostName == profile.HostName {
			return fmt.Errorf("a profile with hostname '%s' already exists", profile.HostName)
		}
	}

	// Append the new profile
	profiles = append(profiles, profile)

	// Marshal and write back to file
	profilesJSON, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profiles to JSON: %v", err)
	}

	// Ensure .ssh directory exists
	if err := os.MkdirAll(filepath.Join(homeDir, ".ssh"), 0700); err != nil {
		return fmt.Errorf("failed to create .ssh directory: %v", err)
	}

	if err := os.WriteFile(filePath, profilesJSON, 0600); err != nil {
		return fmt.Errorf("failed to write profiles to config file: %v", err)
	}

	fmt.Println("✅ Profile added successfully!")
	return nil
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new SSH profile",
	Long:  `Add a new SSH profile by providing hostname, host IP, and path to the SSH key file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve flags
		hostName, _ := cmd.Flags().GetString("hostname")
		hostIP, _ := cmd.Flags().GetString("hostip")
		keyPath, _ := cmd.Flags().GetString("keypath")

		// Prompt for missing flags
		if hostName == "" {
			prompt := &survey.Input{
				Message: "Enter Hostname:",
			}
			survey.AskOne(prompt, &hostName, survey.WithValidator(survey.Required))
		}

		if hostIP == "" {
			prompt := &survey.Input{
				Message: "Enter Host IP Address:",
			}
			survey.AskOne(prompt, &hostIP, survey.WithValidator(survey.Required))
		}

		if keyPath == "" {
			prompt := &survey.Input{
				Message: "Enter Path to SSH Key File:",
			}
			survey.AskOne(prompt, &keyPath, survey.WithValidator(survey.Required))
		}

		if hostName == "" || hostIP == "" || keyPath == "" { // Validate flags
			fmt.Println("❌ All flags are required.")
			os.Exit(1)
		}
		profile := Profile{
			HostName: hostName,
			HostIP:   hostIP,
			KeyPath:  keyPath,
		}

		// Add the profile
		if err := AddProfile(profile); err != nil {
			fmt.Printf("❌ Error adding profile: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Define flags for the add command
	addCmd.Flags().StringP("hostname", "n", "", "Hostname for the SSH profile")
	addCmd.Flags().StringP("hostip", "i", "", "Host IP address for the SSH profile")
	addCmd.Flags().StringP("keypath", "k", "", "Path to the SSH key filepath")

	// Add the add command to the root command
	rootCmd.AddCommand(addCmd)
}
