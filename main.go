package main

import (
	"log"
	"os"
	xml "xmlToTUI/parser"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return ""
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	defer f.Close()

	entries, err := os.ReadDir("./xml")
	if err != nil {
		log.Fatal(err)
	}

	xmls := make(map[string]xml.Node)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := "./xml/" + entry.Name()
		root, err := xml.DecodeXMLFile(path)
		if err != nil {
			log.Printf("Error in %s: %v", path, err)
			continue
		}

		if !root.HasAttr("type") || root.Attr["type"] != "tmul-1.0v" {
			continue
		}

		log.Printf("Loaded file ", entry.Name())
		log.Printf("\n", root.String())
		xmls[entry.Name()] = root
	}

	if _, exist := xmls["index.xml"]; !exist {
		log.Printf("Required valid file index.xml to continue")
		return
	}

	p := tea.NewProgram(xmls["index.xml"].XmlToModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
