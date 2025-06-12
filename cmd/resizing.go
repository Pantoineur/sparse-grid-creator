package main

import tea "github.com/charmbracelet/bubbletea"

type ResizingModel struct {
	m model
}

func (m *ResizingModel) Init() tea.Cmd {
	return nil
}

func (m *ResizingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m *ResizingModel) View() tea.Model {
	return m.m
}
