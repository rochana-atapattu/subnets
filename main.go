package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rochana-atapattu/subnets/internal/subnet"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table  table.Model
	subnet *subnet.Subnet
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "d":
			addr := subnet.InetAton(strings.Split(m.table.SelectedRow()[0], "/")[0])
			maskLen := subnet.MaskLen(subnet.InetAton(m.table.SelectedRow()[1]))
			m.subnet.Find(addr, maskLen).Divide()
		case "j":
			addr := subnet.InetAton(strings.Split(m.table.SelectedRow()[0], "/")[0])
			maskLen := subnet.MaskLen(subnet.InetAton(m.table.SelectedRow()[1]))
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
	m.table, cmd = m.table.Update(msg)
	m.rows()
	return m, cmd
}

func (m *model) rows() {
	var rows []table.Row
	m.subnet.Iterate(func(n *subnet.Subnet) {
		s := subnet.NetworkAddress(n.Address, n.MaskLen)
		lastAddress := subnet.SubnetLastAddress(s, n.MaskLen)
		netmask := subnet.SubnetNetmask(n.MaskLen)
		rows = append(rows, table.Row{
			subnet.InetNtoa(n.Address) + "/" + fmt.Sprint(n.MaskLen),
			subnet.InetNtoa(netmask),
			subnet.InetNtoa(s+1) + " - " + subnet.InetNtoa(lastAddress),
			subnet.InetNtoa(s+1) + " - " + subnet.InetNtoa(lastAddress-1),
			fmt.Sprint(subnet.SubnetAddresses(n.MaskLen)),
		})
	})
	m.table.SetRows(rows)
}
func (m model) View() string {

	return baseStyle.Render(m.table.View()) + "\n"
}

func main() {
	columns := []table.Column{
		{Title: "Subnet address", Width: 15},
		{Title: "Netmask", Width: 15},
		{Title: "Range of addressess", Width: 33},
		{Title: "Useable IPs", Width: 33},
		{Title: "Hosts", Width: 7},
	}

	m := model{}

	m.subnet = &subnet.Subnet{
		Address: subnet.InetAton("10.2.0.0"),
		MaskLen: 16,
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m.table = t

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
