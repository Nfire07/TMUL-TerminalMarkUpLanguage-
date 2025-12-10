package parser

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var defaultStyles = map[string]lipgloss.Style{
	"text":  lipgloss.NewStyle(),
	"title": lipgloss.NewStyle(),
	"list":  lipgloss.NewStyle(),
	"table": lipgloss.NewStyle(),
}

type NodeModel struct {
	Node Node
}

func (node Node) XmlToModel() tea.Model {
	return NodeModel{
		Node: node,
	}
}

func (m NodeModel) Init() tea.Cmd {
	return nil
}

func (m NodeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m NodeModel) View() string {
	return m.nodeToView(m.Node)
}

func applyAttrsToStyle(n Node, base lipgloss.Style) lipgloss.Style {
	style := base
	for k, v := range n.Attr {
		switch strings.ToLower(k) {
		case "bold":
			style = style.Bold(v == "true")
		case "italic":
			style = style.Italic(v == "true")
		case "faint":
			style = style.Faint(v == "true")
		case "underline":
			style = style.Underline(v == "true")
		case "blink":
			style = style.Blink(v == "true")
		case "reverse":
			style = style.Reverse(v == "true")
		case "strikethrough":
			style = style.Strikethrough(v == "true")
		case "color", "foreground":
			style = style.Foreground(lipgloss.Color(v))
		case "bg", "background":
			style = style.Background(lipgloss.Color(v))
		case "padding":
			style = style.Padding(1)
		case "paddingtop":
			if n, err := strconv.Atoi(v); err == nil {
				style = style.PaddingTop(n)
			}
		case "paddingleft":
			if n, err := strconv.Atoi(v); err == nil {
				style = style.PaddingLeft(n)
			}
		case "paddingright":
			if n, err := strconv.Atoi(v); err == nil {
				style = style.PaddingRight(n)
			}
		case "paddingbottom":
			if n, err := strconv.Atoi(v); err == nil {
				style = style.PaddingBottom(n)
			}
		case "margin":
			style = style.Margin(1)
		case "margintop":
			if n, err := strconv.Atoi(v); err == nil {
				style = style.MarginTop(n)
			}
		case "marginbottom":
			if n, err := strconv.Atoi(v); err == nil {
				style = style.MarginBottom(n)
			}
		case "marginleft":
			if n, err := strconv.Atoi(v); err == nil {
				style = style.MarginLeft(n)
			}
		case "marginright":
			if n, err := strconv.Atoi(v); err == nil {
				style = style.MarginRight(n)
			}
		case "border":

			var b lipgloss.Border
			switch strings.ToLower(v) {
			case "normal":
				b = lipgloss.NormalBorder()
			case "rounded":
				b = lipgloss.RoundedBorder()
			case "double":
				b = lipgloss.DoubleBorder()
			case "thick":
				b = lipgloss.ThickBorder()
			case "hidden":
				b = lipgloss.HiddenBorder()
			default:
				b = lipgloss.NormalBorder()
			}
			style = style.Border(b)

			if fg, ok := n.Attr["borderforeground"]; ok {
				style = style.BorderForeground(lipgloss.Color(fg))
			}

			if bg, ok := n.Attr["borderbackground"]; ok {
				style = style.BorderBackground(lipgloss.Color(bg))
			}
		case "borderforeground":
			style = style.BorderForeground(lipgloss.Color(v))
		case "width":
			if n, err := strconv.Atoi(v); err == nil {
				style = style.Width(n)
			}
		case "height":
			if n, err := strconv.Atoi(v); err == nil {
				style = style.Height(n)
			}
		case "align":
			switch strings.ToLower(v) {
			case "left":
				style = style.Align(lipgloss.Left)
			case "center":
				style = style.Align(lipgloss.Center)
			case "right":
				style = style.Align(lipgloss.Right)
			}
		case "maxwidth":
			if n, err := strconv.Atoi(v); err == nil {
				style = style.MaxWidth(n)
			}
		case "maxheight":
			if n, err := strconv.Atoi(v); err == nil {
				style = style.MaxHeight(n)
			}
		}
	}
	return style
}

func (m NodeModel) nodeToView(n Node) string {
	baseStyle, exists := defaultStyles[n.Name]
	if !exists {
		baseStyle = lipgloss.NewStyle()
	}

	style := applyAttrsToStyle(n, baseStyle)

	var sb strings.Builder

	switch n.Name {
	case "text":
		sb.WriteString(style.Render(strings.TrimSpace(n.Content)))

	case "box":
		var content strings.Builder
		for _, c := range n.Children {
			content.WriteString(m.nodeToView(c))
			content.WriteString("\n")
		}
		sb.WriteString(style.Render(content.String()))

	case "list":
		for i, child := range n.Children {
			sb.WriteString(style.Render("• " + m.nodeToView(child)))
			if i != len(n.Children)-1 {
				sb.WriteString("\n")
			}
		}

	default:
		for _, c := range n.Children {
			sb.WriteString(m.nodeToView(c) + "\n")
		}
	}

	return sb.String()
}
