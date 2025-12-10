package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type Node struct {
	Name     string
	Attr     map[string]string
	Children []Node
	Content  string
}

func (n Node) HasAttr(attrName string) bool {
	_, exists := n.Attr[attrName]
	return exists
}

func DecodeXMLFile(filename string) (Node, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Node{}, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	var root Node
	stack := []*Node{}

	for {
		tok, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return Node{}, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			node := Node{
				Name: t.Name.Local,
				Attr: map[string]string{},
			}
			for _, a := range t.Attr {
				node.Attr[a.Name.Local] = a.Value
			}

			if len(stack) == 0 {
				root = node
				stack = append(stack, &root)
			} else {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, node)
				stack = append(stack, &parent.Children[len(parent.Children)-1])
			}

		case xml.EndElement:
			stack = stack[:len(stack)-1]

		case xml.CharData:
			content := strings.TrimSpace(string(t))
			if content != "" && len(stack) > 0 {
				curr := stack[len(stack)-1]
				if curr.Content == "" {
					curr.Content = content
				} else {
					curr.Content += " " + content
				}
			}
		}
	}

	return root, nil
}

func PrintNode(n Node, indent string) {

	fmt.Printf("%s<%s", indent, n.Name)
	for k, v := range n.Attr {
		fmt.Printf(" %s=\"%s\"", k, v)
	}
	fmt.Println(">")
	if n.Content != "" {
		fmt.Printf("%s  %s\n", indent, n.Content)
	}
	for _, c := range n.Children {
		PrintNode(c, indent+"  ")
	}
	fmt.Printf("%s</%s>\n", indent, n.Name)
}

func (n Node) String() string {
	var sb strings.Builder
	n.buildString(&sb, "")
	return sb.String()
}

func (n Node) buildString(sb *strings.Builder, indent string) {
	sb.WriteString(fmt.Sprintf("%s<%s", indent, n.Name))
	for k, v := range n.Attr {
		sb.WriteString(fmt.Sprintf(" %s=\"%s\"", k, v))
	}
	sb.WriteString(">\n")

	if n.Content != "" {
		sb.WriteString(fmt.Sprintf("%s  %s\n", indent, n.Content))
	}

	for _, c := range n.Children {
		c.buildString(sb, indent+"  ")
	}

	sb.WriteString(fmt.Sprintf("%s</%s>\n", indent, n.Name))
}
