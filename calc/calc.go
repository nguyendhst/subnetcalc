package calc

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	max32 uint32 = 0b11111111111111111111111111111111
)

// IPInput is the interface for all input types
type IPInput interface {
	Process() (IPResult, error)
}

// IPResult is the interface for all result types
type IPResult interface {
	String()
}

// ProcessInput processes the given IP input
func ProcessInput(inp IPInput) (IPResult, error) {
	res, err := inp.Process()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// IPv4 is the struct for IPv4 input
type IPv4 struct {
	Addr string
}

// IPv6 is the struct for IPv6 input
type IPv6 struct {
	Addr string
}

// IPv4Result is the struct for IPv4 result
type IPv4Result struct {
	Address        string `json:"address"`
	Prefix         int    `json:"prefix"`
	Mask           string `json:"mask"`
	MaskBinary     string `json:"mask_binary"`
	Wildcard       string `json:"wildcard"`
	WildcardBinary string `json:"wildcard_binary"`
	Lower          string `json:"lower"`
	Upper          string `json:"upper"`
}

// IPv6Result is the struct for IPv6 result
type IPv6Result struct {
	Address        string `json:"address"`
	NetworkAddress string `json:"network_address"`
}

// Process processes the given IPv4 input
func (ipv4 *IPv4) Process() (IPResult, error) {
	ok, err := VerifyIPv4(ipv4.Addr)
	if err != nil {
		fmt.Println(err)
		return &IPv4Result{}, err
	} else if !ok {
		return &IPv4Result{}, fmt.Errorf("invalid IPv4 address")
	}
	addr, k := parse(ipv4.Addr)
	if k <= 0 {
		return &IPv4Result{}, fmt.Errorf("problem with parsing")
	}
	return calcV4(addr, k), nil
}

func parse(addr string) (string, int) {
	re, err := regexp.Compile(`\/[0-9]*`)
	if err != nil {
		fmt.Println(err)
	}
	mask := re.FindStringSubmatch(addr)
	if mask == nil {
		return "", -1
	}
	prefix := mask[0][1:]
	prefixInt, err := strconv.Atoi(prefix)
	if err != nil {
		fmt.Println(err)
	}
	return addr[:len(addr)-len(prefix)-1], prefixInt
}

// Process processes the given IPv6 input
func (ipv6 *IPv6) Process() (IPResult, error) {
	ok, err := VerifyIPv6(ipv6.Addr)
	if err != nil {
		fmt.Println(err)
		return &IPv6Result{}, err
	} else if !ok {
		return &IPv6Result{}, fmt.Errorf("invalid IPv6 address")
	}
	addr, k := parse(ipv6.Addr)
	if k <= 0 {
		return &IPv6Result{}, fmt.Errorf("problem with parsing")
	}
	return calcV6(addr, k), nil
}

// VerifyIPv4 checks if the given IP is a valid IPv4 address with mask
func VerifyIPv4(address string) (bool, error) {
	ipv4, err := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\/((0?[0-9])|([1-2][0-9])|(3[0-2]))$`)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return ipv4.MatchString(address), nil
}

// VerifyIPv6 checks if the given IP is a valid IPv6 address with mask
func VerifyIPv6(address string) (bool, error) {
	ipv6, err := regexp.Compile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))\/(64|([1-5][0-9])|(6[0-4])|(0?[1-9]))$`)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return ipv6.MatchString(address), nil
}

func calcV4(address string, prefix int) IPResult {
	maskResult := maskV4(address, prefix)
	return &IPv4Result{
		Address:        address,
		Prefix:         prefix,
		Mask:           maskResult[0],
		MaskBinary:     maskResult[1],
		Wildcard:       maskResult[2],
		WildcardBinary: maskResult[3],
		Lower:          maskResult[4],
		Upper:          maskResult[5],
	}
}

func maskV4(addr string, prefix int) []string {
	if prefix == 32 {
		return []string{
			"255.255.255.255",
			fmt.Sprintf("%08b. %08b. %08b. %08b", 255, 255, 255, 255),
			"0",
			fmt.Sprintf("%08b. %08b. %08b. %08b", 0, 0, 0, 0),
		}
	}
	mask := max32 << uint32(32-prefix)
	first := mask >> 24 & 0xff
	second := mask >> 16 & 0xff
	third := mask >> 8 & 0xff
	fourth := mask & 0xff

	arr := ip4ToSlice(addr)

	return []string{
		fmt.Sprintf("%03d.%03d.%03d.%03d", first, second, third, fourth),
		fmt.Sprintf("%08b. %08b. %08b. %08b", first, second, third, fourth),
		fmt.Sprintf("%03d.%03d.%03d.%03d", ^uint8(first), ^uint8(second), ^uint8(third), ^uint8(fourth)),
		fmt.Sprintf("%08b. %08b. %08b. %08b", ^uint8(first), ^uint8(second), ^uint8(third), ^uint8(fourth)),
		fmt.Sprintf("%d.%d.%d.%d", (uint8(arr[0]) & uint8(first)), (uint8(arr[1]) & uint8(second)), (uint8(arr[2]) & uint8(third)), (uint8(arr[3]) & uint8(fourth))),
		fmt.Sprintf("%d.%d.%d.%d", ((uint8(arr[0]) & uint8(first)) + ^uint8(first)), ((uint8(arr[1]) & uint8(second)) + ^uint8(second)), ((uint8(arr[2]) & uint8(third)) + ^uint8(third)), ((uint8(arr[3]) & uint8(fourth)) + ^uint8(fourth))),
	}
}

func ip4ToSlice(address string) []uint64 {
	ip := strings.Split(address, ".")
	result := make([]uint64, 4)
	for i, v := range ip {
		result[i], _ = strconv.ParseUint(v, 10, 8)
	}
	return result
}

func calcV6(address string, prefix int) IPResult {
	arr := func(address string) []uint64 {
		ip := strings.Split(address, ":")
		result := make([]uint64, 8)
		for i, v := range ip {
			result[i], _ = strconv.ParseUint(v, 16, 16)
		}
		return result
	}(address)

	block := 7 - ((127 - prefix) / 16)
	bitNum := 15 - ((127 - prefix) % 16)

	res := func(arr []uint64, block int, bit int) []string {
		temp := uint64((arr[block] >> uint64(bit)) << uint64(bit))
		fullip := fmt.Sprintf("%04x:%04x:%04x:%04x:%04x:%04x:%04x:%04x", arr[0], arr[1], arr[2], arr[3], arr[4], arr[5], arr[6], arr[7])
		arr[block] = temp
		for i := 7; i > block; i-- {
			arr[i] = 0
		}
		return []string{
			fullip,
			fmt.Sprintf("%04x:%04x:%04x:%04x:%04x:%04x:%04x:%04x", arr[0], arr[1], arr[2], arr[3], arr[4], arr[5], arr[6], arr[7]),
		}
	}(arr, block, bitNum)

	return &IPv6Result{
		Address:        res[0],
		NetworkAddress: res[1],
	}
}

func (ipv4 *IPv4Result) String() {
	fmt.Printf("IPv4: \t\t\t%s/%d\nMask: \t\t\t%s\nMask Binary: \t\t%s\nWildcard: \t\t%s\nWildcard Binary: \t%s\nLower Bound: \t\t%s\nUpper Bound: \t\t%s\n", ipv4.Address, ipv4.Prefix, ipv4.Mask, ipv4.MaskBinary, ipv4.Wildcard, ipv4.WildcardBinary, ipv4.Lower, ipv4.Upper)
}

func (ipv6 *IPv6Result) String() {
	fmt.Printf("IPv6: \t\t%s\nNetwork: \t%s\n", ipv6.Address, ipv6.NetworkAddress)
}
