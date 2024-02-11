package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rochana-atapattu/subnets/table"
	"github.com/rochana-atapattu/subnets/internal/subnet"
)

type sessionState int

const (
	stateList sessionState = iota
	stateEditing
)

const (
	columnKeyAddr    = "addr"
	columnKeyMask    = "mask"
	columnKeyAddrs   = "addrs"
	columnKeyUseable = "useable"
	columnKeyLabels  = "labels"
	columnKeyHosts   = "hosts"
	columnKeyJoin	   = "join"
)

var (
	styleSubtle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))

	styleBase = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a7a")).
			BorderForeground(lipgloss.Color("#a38")).
			Align(lipgloss.Left)
)

type model struct {
	state     sessionState
	listModel listModel
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch m.state {
	case stateList:
		newList, newCmd := m.listModel.Update(msg)
		listModel, ok := newList.(listModel)
		if !ok {
			panic("unexpected model type")
		}
		m.listModel = listModel
		cmd = newCmd
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	view := m.listModel.View() + "\n"
	return lipgloss.NewStyle().MarginLeft(1).Render(view)
}



type listModel struct {
	subnet *subnet.Subnet
	table  table.Model
}

func (m listModel) Init() tea.Cmd {
	
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.table, cmd = m.table.Update(msg)
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
func (m *listModel) rows() {
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

func (m listModel) View() string {
	return m.table.View()
}

func main() {
	
	if os.Getenv("HELP_DEBUG") != "" {
		f, err := tea.LogToFile("tmp/debug.log", "help")
		if err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		}
		defer f.Close() // nolint:errcheck
	}

	m := model{
		state: stateList,
	}
	columns := []table.Column{
		table.NewColumn(columnKeyAddr, "Subnet address", 15),
		table.NewColumn(columnKeyMask, "Netmask", 15),
		table.NewColumn(columnKeyAddrs, "Range of addressess", 33),
		table.NewColumn(columnKeyUseable, "Useable IPs", 33),
		table.NewColumn(columnKeyHosts, "Hosts", 7),
		table.NewColumn(columnKeyJoin, "Join", 7),
	}
	m.listModel.subnet = &subnet.Subnet{
		Address: subnet.InetAton("10.2.0.0"),
		MaskLen: 16,
	}

	t := table.New(columns)

	m.listModel.table = t.BorderRounded().
		WithBaseStyle(styleBase).
		Focused(true).WithHeaderVisibility(true)

	

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
