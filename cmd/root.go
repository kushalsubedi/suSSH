package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Define the green ASCII art for suSSH
const greenAsciiArt = `

         	@@@@@@@@@@@@@@@            
            @@@@@@@%%%%%@@@@@@@@        
          @@@@@#############%@@@@@      
         @@@@############%%%%%%@@@@     
         @@@%######%%%%%%%%%%%%%%%%%%%  
        @@@%%#####%%%*=------:...:=*%%%%
        @@@%#####%%%*=-------....   -%%%
  @@@@@@@@@%#####%%#+=---------------*%% @@@@@@@@@@
 @@@%%%%@@@%#####%%#+++=-------------*%% @ SUSSH @@
@@@%%%%%@@@%#####%%%*+++=-----------*%%% @@@@@@@@@@
@@@%###%@@%%#####%%%+++++++++++++++++#%%
@@@%%%%%@@%%#####%%%%#++++++++++++*%%%%%
@@@%%%%@@@%%#######%%%%%%%%%%%%%%%%%%@  
@@@%%%%@@@%%%##########%%%%%%%%%###@@@  
@@@%%%%@@@%%%######################%@@  
@@@%%%%@@@%%%######################%@@  
@@@%%%%@@@%%%%#####################@@@  
@@@%%%%@@@%%%%%###################%@@@  
@@@%%%%@@@%%%%%%#################%%@@@  
@@@%%%%@@@%%%%%%%%%###########%%%%%@@@  
@@@%%%%@@@%%%%%%%%%%%%%%%%%%%%%%%%%@@@  
 @@@@%%@@@%%%%%%%%%%%%%%%%%%%%%%%%@@@   
  @@@@@@@@%%%%%%%%%@@@@@@@@@@@%%%%@@@   
     @@@@@@%%%%%%%%@@@@@@@%%%%%%%%@@@   
        @@@%%%%%%%%@@@ @@@%%%%%%%@@@@   
        @@@%%%%%%%%@@@ @@@%%%%%%%@@@    
        @@@%%%%%%%%@@@  @@@@@@@@@@@     
        @@@@%%%%%@@@@@  @@@@@@@@@@      
         @@@@@@@@@@@                    

                                    

`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "suSSH",
	Short: "suSSH is a simple utility to manage your SSH keys",
	Long:  greenAsciiArt + "suSSH is a simple utility to manage your SSH keys. It allows you to add, list, and remove SSH profiles.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(greenAsciiArt)
		// Add your default action here if needed
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Define any flags or configuration settings here
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
