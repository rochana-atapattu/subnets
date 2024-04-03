package subnet

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// inetAton converts a dotted IP address string to a uint32.
func InetAton(ipAddr string) uint32 {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
}

// inetNtoa converts a uint32 to a dotted IP address string.
func InetNtoa(ipInt uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ipInt>>24), byte(ipInt>>16), byte(ipInt>>8), byte(ipInt))
}

// subnetNetmask calculates the subnet mask for a given mask length.
func SubnetNetmask(maskLen uint32) uint32 {
	return ^uint32(0) << (32 - uint(maskLen))
}

func MaskLen(subnetMask uint32) uint32 {
	// Count the number of leading 1s in the subnetMask.
	var maskLen uint32 = 0
	for subnetMask&0x80000000 != 0 {
		maskLen++
		subnetMask <<= 1 // Shift left to check the next bit.
	}
	return maskLen
}

// networkAddress calculates the network address for a given IP address and subnet mask length.
func NetworkAddress(ip, maskLen uint32) uint32 {
	mask := SubnetNetmask(maskLen)
	return ip & mask
}

// subnetAddresses calculates the number of addresses in a subnet based on the mask length.
func SubnetAddresses(maskLen uint32) uint32 {
	return 1 << (32 - uint(maskLen))
}

// subnetLastAddress calculates the last IP address in a subnet.
func SubnetLastAddress(subnet, maskLen uint32) uint32 {
	return subnet + SubnetAddresses(maskLen) - 1
}

// IsValidIPAddress checks if the given string is a valid IPv4 address.
func IsValidIPAddress(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return false
		}
	}
	return true
}
