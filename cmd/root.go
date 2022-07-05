/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Subnet struct {
	IP      string `json:"ip"`
	Netmask string `json:"netmask"`
	Network string `json:"network"`
	Cidr    string `json:"cidr"`
	Class   string `json:"class"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "subnetcalc",
	Short: "A small utility to calculate subnet mask and related information",
	Long: `    _____       _                _            _      
   / ____|     | |              | |          | |     
  | (___  _   _| |__  _ __   ___| |_ ___ __ _| | ___ 
   \___ \| | | | '_ \| '_ \ / _ \ __/ __/ _\ | |/ __|
   ____) | |_| | |_) | | | |  __/ || (_| (_| | | (__ 
  |_____/ \__,_|_.__/|_| |_|\___|\__\___\__,_|_|\___|						 
 `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func Halt(e error) {
	if e != nil {
		fmt.Println("Error:", e)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.subnetcalc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().StringP("ip", "i", "", "IP address")
}
