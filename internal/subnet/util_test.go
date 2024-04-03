package subnet

import "testing"

func TestIsValidIPAddress(t *testing.T) {
	testCases := []struct {
		ip       string
		expected bool
	}{
		{"192.168.1.1", true},
		{"255.255.255.255", true},
		{"0.0.0.0", true},
		{"256.0.0.0", false},
		{"192.168.1", false},
		{"192.168.1.256", false},
		{"192.168.1.-1", false},
		{"192.168.1.01", true},
		{"192.168.1.1.1", false},
		{"192.168..1", false},
		{"abc.def.ghi.jkl", false},
	}

	for _, testCase := range testCases {
		actual := IsValidIPAddress(testCase.ip)
		if actual != testCase.expected {
			t.Errorf("IsValidIPAddress(%q) = %v, expected %v", testCase.ip, actual, testCase.expected)
		}
	}
}
