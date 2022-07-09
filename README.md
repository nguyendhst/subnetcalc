<pre>
    _____       _                _            _      
   / ____|     | |              | |          | |     
  | (___  _   _| |__  _ __   ___| |_ ___ __ _| | ___ 
   \___ \| | | | '_ \| '_ \ / _ \ __/ __/ _\ | |/ __|
   ____) | |_| | |_) | | | |  __/ || (_| (_| | | (__ 
  |_____/ \__,_|_.__/|_| |_|\___|\__\___\__,_|_|\___|

Usage:
  subnetcalc [flags]
  subnetcalc [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  web         Start web server

Flags:
  -h, --help          help for subnetcalc
  -4, --ipv4 string   IPv4 address with mask (192.168.1.1/24)
  -6, --ipv6 string   IPv6 address with prefix (2001:db8:85a3:0000:0000:8a2e:370:7334/64)

Use "subnetcalc [command] --help" for more information about a command.		
</pre>
## Set up dev environment
### Clone the repo and add dependencies
```sh
git clone https://github.com/nguyendhst/subnetcalc.git
go mod tidy
```
