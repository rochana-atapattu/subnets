package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rochana-atapattu/subnets/internal/subnet"
	"github.com/rochana-atapattu/subnets/internal/tree"
	"golang.org/x/term"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var (
	styleDoc  = lipgloss.NewStyle().Padding(1)
	styleHelp = lipgloss.NewStyle().Margin(0, 0, 0, 0).Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
)

type model struct {
	subnet *subnet.Subnet
	tree   tree.Model

	width  int
	height int

	Help     help.Model
	KeyMap   KeyMap
	showHelp bool
}

// KeyMap holds the key bindings for the table.
type KeyMap struct {
	Divide key.Binding
	Join   key.Binding
	Save   key.Binding
	Load   key.Binding
	Quit   key.Binding

	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding
}

// DefaultKeyMap is the default key bindings for the table.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Divide: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "divide"),
		),
		Join: key.NewBinding(
			key.WithKeys("j"),
			key.WithHelp("j", "join"),
		),
		Save: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "save"),
		),
		Load: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "load"),
		),
		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close help"),
		),

		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.Divide):
			node, ok := m.tree.GetNodeAtCurrentCursor()
			if !ok {
				fmt.Println("No node found")
				return m, nil
			}
			// fmt.Println("Node found", node.Value)
			s := strings.Split(node.Value, "/")
			addr := subnet.InetAton(s[0])
			// maskInt is s[1] in uint32
			maskInt, err := strconv.Atoi(s[1])
			if err != nil {
				fmt.Println("Error converting mask to int:", err)
				return m, nil
			}
			maskLen := uint32(maskInt)
			m.subnet.Find(addr, maskLen).Divide()
		case key.Matches(msg, m.KeyMap.Join):
			node, ok := m.tree.GetNodeAtCurrentCursor()
			if !ok {
				fmt.Println("No node found")
				return m, nil
			}
			// fmt.Println("Node found", node.Value)
			s := strings.Split(node.Value, "/")
			addr := subnet.InetAton(s[0])
			// maskInt is s[1] in uint32
			maskInt, err := strconv.Atoi(s[1])
			if err != nil {
				fmt.Println("Error converting mask to int:", err)
				return m, nil
			}
			maskLen := uint32(maskInt)
			m.subnet.Find(addr, maskLen).Join()
		case key.Matches(msg, m.KeyMap.Save):
			subnet.SaveTree(m.subnet, "subnets.json")
		case key.Matches(msg, m.KeyMap.Load):
			subnet, err := subnet.LoadTree("subnets.json")
			if err != nil {
				fmt.Println("Error loading subnet tree:", err)
				return m, nil
			}
			m.subnet = subnet
		case key.Matches(msg, m.KeyMap.ShowFullHelp):
			fallthrough
		case key.Matches(msg, m.KeyMap.CloseFullHelp):
			m.Help.ShowAll = !m.Help.ShowAll
		}
	}
	m.rows()
	m.tree, cmd = m.tree.Update(msg)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	availableHeight := m.height

	var sections []string

	var help string
	if m.showHelp {
		help = m.helpView()
		availableHeight -= lipgloss.Height(help)
	}
	sections = append(sections, lipgloss.NewStyle().Height(availableHeight).Render(m.tree.View(), help))
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *model) rows() {

	m.tree.SetNodes([]tree.Node{toNodeTree(m.subnet)})

}

func (m *model) SetShowHelp() bool {
	return m.showHelp
}
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: subnets <IP address> <mask length>")
		os.Exit(1)
	}
	ipAddr := os.Args[1]
	maskLengthStr := os.Args[2]

	// Validate the provided IP address
	if !subnet.IsValidIPAddress(ipAddr) {
		fmt.Println("Invalid IP address:", ipAddr)
		os.Exit(1)
	}

	// Convert the mask length from string to integer
	maskLength, err := strconv.Atoi(maskLengthStr)
	if err != nil || maskLength < 0 || maskLength > 32 {
		fmt.Println("Invalid mask length:", maskLengthStr)
		os.Exit(1)
	}

	if os.Getenv("HELP_DEBUG") != "" {
		f, err := tea.LogToFile("tmp/debug.log", "help")
		if err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		}
		defer f.Close() // nolint:errcheck
	}

	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 80
		h = 24
	}
	top, right, bottom, left := styleDoc.GetPadding()
	w = w - left - right
	h = h - top - bottom - 10

	// Use the provided IP address and mask length
	m := model{
		showHelp: true,
		Help:     help.New(),
		KeyMap:   DefaultKeyMap(),
		height:   h,
		width:    w,
	}
	m.subnet = &subnet.Subnet{
		Address: subnet.InetAton(ipAddr),
		MaskLen: uint32(maskLength),
	}

	nodes := []tree.Node{toNodeTree(m.subnet)}
	m.tree = tree.New(nodes)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
func toNodeTree(n *subnet.Subnet) tree.Node {
	// Convert the subnet's address and mask length to a string representation.
	// This will be the node's value.
	value := fmt.Sprintf("%s/%s", subnet.InetNtoa(n.Address), fmt.Sprint(n.MaskLen))

	// For the description, you might want to add additional information from the Subnet,
	// such as its labels or whether it's a left or right child.
	// This example simply joins the labels into a single string.
	s := subnet.NetworkAddress(n.Address, n.MaskLen)
	lastAddress := subnet.SubnetLastAddress(s, n.MaskLen)
	netmask := subnet.SubnetNetmask(n.MaskLen)
	columnKeyMask := subnet.InetNtoa(netmask)
	columnKeyAddrs := subnet.InetNtoa(s+1) + " - " + subnet.InetNtoa(lastAddress)
	columnKeyUseable := subnet.InetNtoa(s+1) + " - " + subnet.InetNtoa(lastAddress-1)
	columnKeyHosts := fmt.Sprint(subnet.SubnetAddresses(n.MaskLen))
	desc := fmt.Sprintf("| Netmask: %s | Range of addressess %s | Useable IPs %s | Hosts %s |", columnKeyMask, columnKeyAddrs, columnKeyUseable, columnKeyHosts)

	// Initialize the Node with the value and description.
	node := tree.Node{
		Value: value,
		Desc:  desc,
	}

	// Recursively convert the Subnet's children to Nodes and add them to the current Node's children.
	children := []tree.Node{}
	if n.Left != nil {
		children = append(children, toNodeTree(n.Left))
	}
	if n.Right != nil {
		children = append(children, toNodeTree(n.Right))
	}
	node.Children = children

	return node
}

func (m model) helpView() string {
	return styleHelp.Render(m.Help.View(m))
}

func (m model) ShortHelp() []key.Binding {
	sh := m.tree.ShortHelp()
	kb := []key.Binding{
		m.KeyMap.Divide,
		m.KeyMap.Join,
		m.KeyMap.Quit,
	}

	return append(kb,
		sh...,
	)
}

func (m model) FullHelp() [][]key.Binding {
	fh := m.tree.FullHelp()
	kb := [][]key.Binding{{
		m.KeyMap.Divide,
		m.KeyMap.Join,
		m.KeyMap.Quit,

		m.KeyMap.CloseFullHelp,
	}}

	return append(kb,
		fh...)
}
