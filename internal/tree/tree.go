package tree

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	bottomLeft string = " └──"

	white  = lipgloss.Color("#ffffff")
	black  = lipgloss.Color("#000000")
	grey = lipgloss.Color("#7c7980")
)

type Styles struct {
	Shapes     lipgloss.Style
	Selected   lipgloss.Style
	Unselected lipgloss.Style
}

func defaultStyles() Styles {
	return Styles{
		Shapes:     lipgloss.NewStyle().Margin(0, 0, 0, 0).Foreground(grey),
		Selected:   lipgloss.NewStyle().Margin(0, 0, 0, 0).Background(grey),
		Unselected: lipgloss.NewStyle().Margin(0, 0, 0, 0).Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}),
	}
}

type Node struct {
	Value    string
	Desc     string
	Children []Node
}

type Model struct {
	KeyMap KeyMap
	Styles Styles

	nodes  []Node
	cursor int


}

func New(nodes []Node) Model {
	return Model{
		KeyMap: DefaultKeyMap(),
		Styles: defaultStyles(),

		nodes:  nodes,

	}
}

// KeyMap holds the key bindings for the table.
type KeyMap struct {
	Bottom      key.Binding
	Top         key.Binding
	SectionDown key.Binding
	SectionUp   key.Binding
	Down        key.Binding
	Up          key.Binding
	Quit        key.Binding
}

// DefaultKeyMap is the default key bindings for the table.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Bottom: key.NewBinding(
			key.WithKeys("bottom"),
			key.WithHelp("end", "bottom"),
		),
		Top: key.NewBinding(
			key.WithKeys("top"),
			key.WithHelp("home", "top"),
		),
		SectionDown: key.NewBinding(
			key.WithKeys("secdown"),
			key.WithHelp("secdown", "section down"),
		),
		SectionUp: key.NewBinding(
			key.WithKeys("secup"),
			key.WithHelp("secup", "section up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
	}
}

func (m Model) Nodes() []Node {
	return m.nodes
}

func (m *Model) SetNodes(nodes []Node) {
	m.nodes = nodes
}

func (m *Model) NumberOfNodes() int {
	count := 0

	var countNodes func([]Node)
	countNodes = func(nodes []Node) {
		for _, node := range nodes {
			count++
			if node.Children != nil {
				countNodes(node.Children)
			}
		}
	}

	countNodes(m.nodes)

	return count

}


func (m Model) Cursor() int {
	return m.cursor
}

func (m *Model) SetCursor(cursor int) {
	m.cursor = cursor
}



func (m *Model) NavUp() {
	m.cursor--

	if m.cursor < 0 {
		m.cursor = 0
		return
	}

}

func (m *Model) NavDown() {
	m.cursor++

	if m.cursor >= m.NumberOfNodes() {
		m.cursor = m.NumberOfNodes() - 1
		return
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Up):
			m.NavUp()
		case key.Matches(msg, m.KeyMap.Down):
			m.NavDown()
		}	
	}

	return m, nil
}

func (m Model) View() string {
	nodes := m.Nodes()
	count := 0 // This is used to keep track of the index of the node we are on (important because we are using a recursive function)
	str := m.renderTree(m.nodes, 0, &count)

	if len(nodes) == 0 {
		return "No data"
	}
	return str
}

func (m *Model) renderTree(remainingNodes []Node, indent int, count *int) string {
	var b strings.Builder

	for _, node := range remainingNodes {

		var str string

		// If we aren't at the root, we add the arrow shape to the string
		if indent > 0 {
			shape := strings.Repeat(" ", (indent-1)*2) + m.Styles.Shapes.Render(bottomLeft) + " "
			str += shape
		}

		// Generate the correct index for the node
		idx := *count
		*count++

		// Format the string with fixed width for the value and description fields
		valueWidth := 12 + indent
		// descWidth := 20
		valueStr := fmt.Sprintf("%-*s", valueWidth, node.Value)
		// Calculate the fixed starting position for descStr
		const startDescAt = 12
		paddingNeeded := startDescAt - len(valueStr) - len("\t\t") // Adjust based on actual space taken by valueStr and tabs

		// Ensure padding is not negative
		if paddingNeeded < 0 {
			paddingNeeded = 0
		}

		// Use spaces to adjust the starting position of descStr
		descStr := fmt.Sprintf("%s%s", strings.Repeat(" ", paddingNeeded), node.Desc)
		// descStr := fmt.Sprintf("%s", node.Desc)

		// If we are at the cursor, we add the selected style to the string
		if m.cursor == idx {
			str += fmt.Sprintf("%s\t\t%s\n", m.Styles.Selected.Render(valueStr), m.Styles.Selected.Render(descStr))
		} else {
			str += fmt.Sprintf("%s\t\t%s\n", m.Styles.Unselected.Render(valueStr), m.Styles.Unselected.Render(descStr))
		}

		b.WriteString(str)

		if node.Children != nil {
			childStr := m.renderTree(node.Children, indent+1, count)
			b.WriteString(childStr)
		}
	}

	return b.String()
}

func (m *Model) currentCursorNode(remaningNodes []Node, indent int, count *int) (Node, bool) {
	for _, node := range remaningNodes {
		idx := *count
		*count++

		if m.cursor == idx {
			return node, true
		}

		if node.Children != nil {
			childNode, ok := m.currentCursorNode(node.Children, indent+1, count)
			if ok {
				return childNode, true
			}
		}
	}

	return Node{}, false
}

// GetNodeAtBFSIndex returns the node at a specific index in a breadth-first search traversal.
func (m Model) GetNodeAtCurrentCursor() (Node, bool) {
	count := 0
	return m.currentCursorNode(m.nodes, 0, &count)
}
func (m Model) ShortHelp() []key.Binding {
	kb := []key.Binding{
		m.KeyMap.Up,
		m.KeyMap.Down,
	}

	return kb
}

func (m Model) FullHelp() [][]key.Binding {
	kb := [][]key.Binding{{
		m.KeyMap.Up,
		m.KeyMap.Down,
	}}

	return kb
}


