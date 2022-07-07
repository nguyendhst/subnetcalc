package calc

import (
	"fmt"
	_ "math/big"
	"regexp"
	"strconv"
	"strings"
)

const (
	max32 uint32 = 0b11111111111111111111111111111111
)

type IPInput interface {
	Process() (IPResult, error)
}

type IPResult interface {
	String()
}

func ProcessInput(inp IPInput) (IPResult, error) {
	res, err := inp.Process()
	if err != nil {
		return nil, err
	}
	return res, nil
}

type IPv4 struct {
	Addr string
}

type IPv6 struct {
	Addr string
}

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

type IPv6Result struct {
	Address     string `json:"address"`
	FullAddress string `json:"full_address"`
	TotalHost   int    `json:"total_host"`
}

func (ipv4 *IPv4) Process() (IPResult, error) {
	ok, err := VerifyIPv4(ipv4.Addr)
	if err != nil {
		fmt.Println(err)
		return IPv4Result{}, err
	} else if !ok {
		return IPv4Result{}, fmt.Errorf("invalid IPv4 address")
	}
	addr, k := Parse(ipv4.Addr)
	if k < 0 {
		return IPv4Result{}, fmt.Errorf("problem with parsing")
	}
	return CalcV4(addr, k), nil
}

func Parse(addr string) (string, int) {
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

func (ipv6 *IPv6) Process() (IPResult, error) {
	ok, err := VerifyIPv6(ipv6.Addr)
	if err != nil {
		fmt.Println(err)
		return IPv6Result{}, err
	} else if !ok {
		return IPv6Result{}, fmt.Errorf("invalid IPv6 address")
	}
	addr, k := Parse(ipv6.Addr)
	if k < 0 {
		return IPv6Result{}, fmt.Errorf("problem with parsing")
	}
	return CalcV6(addr, k), nil
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

func CalcV4(address string, prefix int) IPResult {
	maskResult := maskV4(address, prefix)
	return IPv4Result{
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

	arr := ipToSlice(addr)

	return []string{
		fmt.Sprintf("%03d.%03d.%03d.%03d", first, second, third, fourth),
		fmt.Sprintf("%08b. %08b. %08b. %08b", first, second, third, fourth),
		fmt.Sprintf("%03d.%03d.%03d.%03d", ^uint8(first), ^uint8(second), ^uint8(third), ^uint8(fourth)),
		fmt.Sprintf("%08b. %08b. %08b. %08b", ^uint8(first), ^uint8(second), ^uint8(third), ^uint8(fourth)),
		fmt.Sprintf("%d.%d.%d.%d", (uint8(arr[0]) & uint8(first)), (uint8(arr[1]) & uint8(second)), (uint8(arr[2]) & uint8(third)), (uint8(arr[3]) & uint8(fourth))),
		fmt.Sprintf("%d.%d.%d.%d", ((uint8(arr[0]) & uint8(first)) + ^uint8(first)), ((uint8(arr[1]) & uint8(second)) + ^uint8(second)), ((uint8(arr[2]) & uint8(third)) + ^uint8(third)), ((uint8(arr[3]) & uint8(fourth)) + ^uint8(fourth))),
	}
}

func ipToSlice(address string) []int {
	ip := strings.Split(address, ".")
	result := make([]int, 4)
	for i, v := range ip {
		result[i], _ = strconv.Atoi(v)
	}
	return result
}

func CalcV6(address string, prefix int) IPResult {
	return IPv6Result{
		Address: address,
	}
}

func (ipv4 IPv4Result) String() {
	fmt.Printf("IPv4: \t\t\t%s/%d\n", ipv4.Address, ipv4.Prefix)
	fmt.Printf("Mask: \t\t\t%s\n", ipv4.Mask)
	fmt.Printf("Mask Binary: \t\t%s\n", ipv4.MaskBinary)
	fmt.Printf("Wildcard: \t\t%s\n", ipv4.Wildcard)
	fmt.Printf("Wildcard Binary: \t%s\n", ipv4.WildcardBinary)
	fmt.Printf("Lower: \t\t\t%s\n", ipv4.Lower)
	fmt.Printf("Upper: \t\t\t%s\n", ipv4.Upper)
}

func (ipv6 IPv6Result) String() {
	fmt.Printf("IPv6: %s/%s\n", ipv6.Address, ipv6.FullAddress)
}
