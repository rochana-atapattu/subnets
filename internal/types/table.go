package types

type CalculationRequest struct {
	Network string `json:"network" form:"network"`
	Netbits uint32 `json:"netbits" form:"netbits"`
}

type TableRowData struct {
	SubnetAddress       string
	Netbits             uint32
	ParentSubnetAddress string
	ParentNetbits       uint32
	Netmask             string
	RangeOfAddresses    string
	UseableIPs          string
	Hosts               string
	Divide              string
	Join                string
}

