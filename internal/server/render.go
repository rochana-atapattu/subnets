package server

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/rochana-atapattu/subnets/internal/subnet"
	"github.com/rochana-atapattu/subnets/internal/types"
	"github.com/rochana-atapattu/subnets/internal/view"
)

func render(c echo.Context, cmp templ.Component) error {
c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}


func viewRows(root *subnet.Subnet) templ.Component {
	var rows []types.TableRowData
	root.Iterate(func(n *subnet.Subnet) {
		s := subnet.NetworkAddress(n.Address, n.MaskLen)
		lastAddress := subnet.SubnetLastAddress(s, n.MaskLen)
		netmask := subnet.SubnetNetmask(n.MaskLen)
		rows = append(rows, types.TableRowData{
			SubnetAddress:    subnet.InetNtoa(n.Address),
			Netbits:          n.MaskLen,
			ParentSubnetAddress: subnet.InetNtoa(n.Parent.Address),
			ParentNetbits:    n.Parent.MaskLen,
			Netmask:          subnet.InetNtoa(netmask),
			RangeOfAddresses: subnet.InetNtoa(s+1) + " - " + subnet.InetNtoa(lastAddress),
			UseableIPs:       subnet.InetNtoa(s+1) + " - " + subnet.InetNtoa(lastAddress-1),
			Hosts:            fmt.Sprint(subnet.SubnetAddresses(n.MaskLen)),
			Divide:           "Divide",
			Join:             "Join",
		})
	})
	return view.RowComponent(view.Row(rows))
}

