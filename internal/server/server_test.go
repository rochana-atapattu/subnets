package server

import (
	"testing"

	"github.com/rochana-atapattu/subnets/internal/subnet"
)

// TestInetAtonNtoa tests the conversion between string IP addresses and uint32 representation.
func TestInetAtonNtoa(t *testing.T) {
	testIP := "192.168.1.1"
	expectedUint := uint32(3232235777)
	resultUint := subnet.InetAton(testIP)
	if resultUint != expectedUint {
		t.Errorf("inetAton(%s) = %d; want %d", testIP, resultUint, expectedUint)
	}

	resultIP := subnet.InetNtoa(expectedUint)
	if resultIP != testIP {
		t.Errorf("inetNtoa(%d) = %s; want %s", expectedUint, resultIP, testIP)
	}
}

// TestSubnetCalculations tests the subnet mask, network address, and last address calculations.
func TestSubnetCalculations(t *testing.T) {
	ip := subnet.InetAton("10.2.0.0")
	maskLen := uint32(16)

	expectedNetmask := uint32(0xffff0000) // 255.255.0.0
	if result := subnet.SubnetNetmask(maskLen); result != expectedNetmask {
		t.Errorf("subnetNetmask(%d) = %x; want %x", maskLen, result, expectedNetmask)
	}

	expectedNetwork := uint32(0x0a020000) // 10.2.0.0
	if result := subnet.NetworkAddress(ip, maskLen); result != expectedNetwork {
		t.Errorf("networkAddress(%x, %d) = %x; want %x", ip, maskLen, result, expectedNetwork)
	}

	expectedLastAddress := uint32(0x0a02ffff) // 10.2.255.255
	if result := subnet.SubnetLastAddress(ip, maskLen); result != expectedLastAddress {
		t.Errorf("subnetLastAddress(%x, %d) = %x; want %x", ip, maskLen, result, expectedLastAddress)
	}
}

// TestSubnetDivision tests the division of subnets into two for various scenarios.
func TestSubnetDivision(t *testing.T) {
	cases := []struct {
		name              string
		address           string
		initialMaskLen    uint32
		expectedLeftAddr  string
		expectedLeftMask  uint32
		expectedRightAddr string
		expectedRightMask uint32
	}{
		{
			name:              "Dividing a /16 subnet",
			address:           "10.2.0.0",
			initialMaskLen:    16,
			expectedLeftAddr:  "10.2.0.0",
			expectedLeftMask:  17,
			expectedRightAddr: "10.2.128.0",
			expectedRightMask: 17,
		},
		{
			name:              "Dividing a /24 subnet",
			address:           "192.168.1.0",
			initialMaskLen:    24,
			expectedLeftAddr:  "192.168.1.0",
			expectedLeftMask:  25,
			expectedRightAddr: "192.168.1.128",
			expectedRightMask: 25,
		},
		{
			name:              "Edge case: Dividing a /31 subnet",
			address:           "192.168.1.0",
			initialMaskLen:    31,
			expectedLeftAddr:  "192.168.1.0",
			expectedLeftMask:  32,
			expectedRightAddr: "192.168.1.1",
			expectedRightMask: 32,
		},
		// {
		//     name:               "Edge case: Attempting to divide a /32 subnet",
		//     address:            "192.168.1.1",
		//     initialMaskLen:     32,
		//     expectedLeftAddr:   "192.168.1.1", // Should remain unchanged
		//     expectedLeftMask:   32,
		//     expectedRightAddr:  "192.168.1.1", // No right child should be created
		//     expectedRightMask:  32, // This indicates an invalid scenario, handled in the test
		// },
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := &subnet.Subnet{
				Address: subnet.InetAton(c.address),
				MaskLen: c.initialMaskLen,
			}
			root.Divide()

			leftAddr := subnet.InetNtoa(root.Left.Address)
			rightAddr := ""
			if root.Right != nil { // Check if right child exists (it should not for a /32 subnet)
				rightAddr = subnet.InetNtoa(root.Right.Address)
			}

			if leftAddr != c.expectedLeftAddr {
				t.Errorf("Left child Address = %s; want %s", leftAddr, c.expectedLeftAddr)
			}
			if root.Left.MaskLen != c.expectedLeftMask {
				t.Errorf("Left child MaskLen = %d; want %d", root.Left.MaskLen, c.expectedLeftMask)
			}
			if rightAddr != c.expectedRightAddr {
				t.Errorf("Right child Address = %s; want %s", rightAddr, c.expectedRightAddr)
			}
			if root.Right != nil && root.Right.MaskLen != c.expectedRightMask {
				t.Errorf("Right child MaskLen = %d; want %d", root.Right.MaskLen, c.expectedRightMask)
			}
		})
	}
}
