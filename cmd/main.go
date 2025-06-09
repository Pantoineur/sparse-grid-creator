package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var boldStyle = lipgloss.NewStyle().
	Bold(true)
var cursorStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#00FF00"))

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

const (
	PaintPath = iota
	PaintObstacle
	Eraser
)

var typeName = map[int]string{
	PaintPath:     "Path",
	PaintObstacle: "Obstacle",
	Eraser:        "Eraser",
}

type model struct {
	grid              Grid
	filled            map[Cell]int
	cursor            Cell
	additionalCursors map[Cell]bool
	isPainting        bool
	paintType         int
	ready             bool
	viewport          viewport.Model
	content           string
	showCells         bool
	mouseEvent        tea.MouseEvent
}

func initialModel() model {
	return model{
		grid: Grid{
			Size:  30,
			Cells: generateCells(30),
		},
		filled:            make(map[Cell]int),
		additionalCursors: make(map[Cell]bool),
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
			m.isPainting = !m.isPainting

		case "m":
			m.paintType++
			if m.paintType > 2 {
				m.paintType = 0
			}

			// TODO
		//case "alt+right":
		//	m.additionalCursors[Cell{X: m.cursor.X + 1, Y: m.cursor.Y}] = true
		//
		//case "alt+left":
		//	m.additionalCursors[Cell{X: m.cursor.X - 1, Y: m.cursor.Y}] = true
		//
		//case "alt+up":
		//	m.additionalCursors[Cell{X: m.cursor.X, Y: m.cursor.Y - 1}] = true
		//
		//case "alt+down":
		//	m.additionalCursors[Cell{X: m.cursor.X, Y: m.cursor.Y + 1}] = true

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.filled[m.cursor]
			if ok {
				delete(m.filled, m.cursor)
			} else {
				m.filled[m.cursor] = m.paintType
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

		// TODO
		//case tea.MouseMsg:
		//	return m, tea.Printf("(X: %d, Y: %d) %s", msg.X, msg.Y, tea.MouseEvent(msg))
	}

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

		checked := " "
		if cellType, ok := m.filled[cell]; ok {
			var c string
			switch cellType {
			case PaintPath:
				c = "X"
			case PaintObstacle:
				c = "O"
			case Eraser:
				delete(m.filled, cell)
			}

			checked = boldStyle.Render(c)
		}

		cursor := ""
		if m.cursor == cell {
			cursor = cursorStyle.Render("O")

			if m.isPainting {
				m.filled[cell] = m.paintType
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
	if m.isPainting {
		pm = "X"
	}
	s := ""
	s += fmt.Sprintf("\nPainting [%s] Type [%s]\n", boldStyle.Render(pm), boldStyle.Render(typeName[m.paintType]))
	s += "\nPress q to quit.\n"

	return s
}
