/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Ipv4 string
var Ipv6 string
var v4Flag bool
var v6Flag bool

type SubnetInput struct {
	Address string `json:"address"`
	Prefix  string `json:"prefix"`
}

var Input SubnetInput = SubnetInput{}

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
	Run: func(cmd *cobra.Command, args []string) {
		err := flagCheck()
		if err != nil {
			Halt(err)
		}

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

func Halt(e error) {
	if e != nil {
		fmt.Println("Error:", e)
		os.Exit(1)
	}
}

func flagCheck() error {
	if v4Flag && v6Flag {
		return fmt.Errorf("only one IP address can be specified")
	} else if !v4Flag && !v6Flag {
		return fmt.Errorf("one IP address must be specified")
	}
	return nil
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.subnetcalc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().StringP("ip", "i", "", "IP address")
	rootCmd.Flags().StringVarP(&Ipv4, "ipv4", "4", "", "IPv4 address with mask (192.168.1.1/24)")
	rootCmd.Flags().StringVarP(&Ipv6, "ipv6", "6", "", "IPv6 address with prefix (2001:db8:85a3::8a2e:370:7334/64)")
	v4Flag = rootCmd.Flags().Lookup("ipv4").Changed
	v6Flag = rootCmd.Flags().Lookup("ipv6").Changed
}
