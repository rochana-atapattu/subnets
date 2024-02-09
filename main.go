package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/rochana-atapattu/subnets/internal/subnet"
)

const (
	columnKeyAddr    = "addr"
	columnKeyMask    = "mask"
	columnKeyAddrs   = "addrs"
	columnKeyUseable = "useable"
	columnKeyLabels = "labels"
	columnKeyHosts   = "hosts"
)

var (
	styleSubtle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))

	styleBase = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a7a")).
			BorderForeground(lipgloss.Color("#a38")).
			Align(lipgloss.Left)
)

type model struct {
	table             table.Model
	subnet            *subnet.Subnet
	currentSubnet     *subnet.Subnet
	lastSelectedEvent table.UserEventRowSelectToggled
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "d":
			addr := subnet.InetAton(strings.Split(m.table.HighlightedRow().Data[columnKeyAddr].(string), "/")[0])
			maskLen := subnet.MaskLen(subnet.InetAton(m.table.HighlightedRow().Data[columnKeyMask].(string)))
			m.subnet.Find(addr, maskLen).Divide()
		case "j":
			addr := subnet.InetAton(strings.Split(m.table.HighlightedRow().Data[columnKeyAddr].(string), "/")[0])
			maskLen := subnet.MaskLen(subnet.InetAton(m.table.HighlightedRow().Data[columnKeyMask].(string)))
			m.subnet.Find(addr, maskLen).Join()
		case "s":
			subnet.SaveTree(m.subnet, "subnets.json")
		case "l":
			subnet, err := subnet.LoadTree("subnets.json")
			if err != nil {
				fmt.Println("Error loading subnet tree:", err)
				return m, nil
			}
			m.subnet = subnet
		}
	}
	m.rows()
	return m, cmd
}

func (m *model) rows() {
	rows := []table.Row{}
	m.subnet.Iterate(func(n *subnet.Subnet) {
		s := subnet.NetworkAddress(n.Address, n.MaskLen)
		lastAddress := subnet.SubnetLastAddress(s, n.MaskLen)
		netmask := subnet.SubnetNetmask(n.MaskLen)
		rows = append(rows, table.NewRow(table.RowData{
			columnKeyAddr:    subnet.InetNtoa(n.Address) + "/" + fmt.Sprint(n.MaskLen),
			columnKeyMask:    subnet.InetNtoa(netmask),
			columnKeyAddrs:   subnet.InetNtoa(s+1) + " - " + subnet.InetNtoa(lastAddress),
			columnKeyUseable: subnet.InetNtoa(s+1) + " - " + subnet.InetNtoa(lastAddress-1),
			columnKeyHosts:   fmt.Sprint(subnet.SubnetAddresses(n.MaskLen)),
		}))
	})

	m.table = m.table.WithRows(rows)
		
}

func (m model) View() string {
	view := m.table.View() + "\n"
	return lipgloss.NewStyle().MarginLeft(1).Render(view)
}

func main() {
	columns := []table.Column{
		table.NewColumn(columnKeyAddr, "Subnet address", 15),
		table.NewColumn(columnKeyMask, "Netmask", 15),
		table.NewColumn(columnKeyAddrs, "Range of addressess", 33),
		table.NewColumn(columnKeyUseable, "Useable IPs", 33),
		table.NewColumn(columnKeyLabels, "Labels", 33),
		table.NewColumn(columnKeyHosts, "Hosts", 7),
	}

	m := model{}

	m.subnet = &subnet.Subnet{
		Address: subnet.InetAton("10.2.0.0"),
		MaskLen: 16,
	}

	t := table.New(columns)

	m.table = t.BorderRounded().
		WithBaseStyle(styleBase).
		WithPageSize(15).
		Focused(true)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
