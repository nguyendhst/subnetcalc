package cmd

import (
	"fmt"
	"os"

	"github.com/nguyendhst/subnetcalc/calc"
	"github.com/spf13/cobra"
)

var Ipv4 string
var Ipv6 string
var v4Flag bool
var v6Flag bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "subnetcalc",
	Short: "A small utility to calculate subnet mask and related information",
	Long:  logo,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var input calc.IPInput
		v4Flag = cmd.Flags().Lookup("ipv4").Changed
		v6Flag = cmd.Flags().Lookup("ipv6").Changed

		err := flagCheck()
		if err != nil {
			Halt(err)
		}

		if v4Flag {
			input = &calc.IPv4{Addr: Ipv4}
		} else {
			input = &calc.IPv6{Addr: Ipv6}
		}

		res, err := calc.ProcessInput(input)
		if err != nil {
			Halt(err)
		}
		res.String()
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

// Halt prints error and exits
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
	rootCmd.Flags().StringVarP(&Ipv6, "ipv6", "6", "", "IPv6 address with prefix (2001:db8:85a3:0000:0000:8a2e:370:7334/64)")
}

var logo = `  _____       _                _            _      
 / ____|     | |              | |          | |     
| (___  _   _| |__  _ __   ___| |_ ___ __ _| | ___ 
 \___ \| | | | '_ \| '_ \ / _ \ __/ __/ _\ | |/ __|
 ____) | |_| | |_) | | | |  __/ || (_| (_| | | (__ 
|_____/ \__,_|_.__/|_| |_|\___|\__\___\__,_|_|\___|						 
`
