package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type Cell struct {
	X, Y int
}

type Grid struct {
	Size  int
	Cells []Cell
}

type model struct {
	grid      Grid
	filled    map[Cell]bool
	cursor    Cell
	paintMode bool
	ready     bool
	viewport  viewport.Model
	content   string
}

func initialModel() model {
	return model{
		grid: Grid{
			Size:  30,
			Cells: generateCells(30),
		},
		filled: make(map[Cell]bool),
	}
}

func generateCells(size int) []Cell {
	cells := make([]Cell, size*size)

	for y := range size {
		for x := range size {
			cells[y*size+x] = Cell{
				X: x,
				Y: y,
			}
		}
	}

	return cells
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor.Y > 0 {
				m.cursor.Y--
			} else if m.cursor.Y == 0 {
				m.cursor.Y = m.grid.Size - 1
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor.Y == m.grid.Size-1 {
				m.cursor.Y = 0
			} else if m.cursor.Y < m.grid.Size-1 {
				m.cursor.Y++
			}

		// The "left" and "h" keys move the cursor left
		case "left", "h":
			if m.cursor.X > 0 {
				m.cursor.X--
			} else if m.cursor.X == 0 {
				m.cursor.X = m.grid.Size - 1
			}

		// The "right" and "l" keys move the cursor right
		case "right", "l":
			if m.cursor.X == m.grid.Size-1 {
				m.cursor.X = 0
			} else if m.cursor.X < m.grid.Size-1 {
				m.cursor.X++
			}

		case "p":
			m.paintMode = !m.paintMode

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.filled[m.cursor]
			if ok {
				delete(m.filled, m.cursor)
			} else {
				m.filled[m.cursor] = true
			}
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		fmt.Printf("Received windowsize %d %d", headerHeight, footerHeight)

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	m.viewport, _ = m.viewport.Update(msg)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	count := 0
	m.content = ""

	// Iterate over our grid
	for _, cell := range m.grid.Cells {

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.filled[cell]; ok {
			checked = "X" // selected!
		}

		// Is the cursor pointing at this choice?
		cursor := "" // no cursor
		if m.cursor == cell {
			cursor = "X" // cursor!

			if m.paintMode {
				m.filled[cell] = true
			}
			checked = ""
		}

		newLine := ""

		count++
		if count == m.grid.Size {
			newLine = "\n"
			count = 0
		}

		// Render the row
		m.content += fmt.Sprintf("[%s%s]%s", cursor, checked, newLine)
	}
	m.viewport.SetContent(m.content)

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	return ""
}

func (m model) footerView() string {
	// The footer

	pm := ""
	if m.paintMode {
		pm = "X"
	}
	s := ""
	s += fmt.Sprintf("\nPainting [%s]\n", pm)
	s += "\nPress q to quit.\n"

	return s
}
