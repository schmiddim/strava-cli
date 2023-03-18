package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type EditorModel struct {
	nameInput        textinput.Model
	descriptionInput textarea.Model
	help             help.Model
	hasUpdate        bool
}

type EditorKeys struct {
	Save        key.Binding
	Quit        key.Binding
	SwitchInput key.Binding
}

func NewEditorModel(currentName string, currentDescription string) EditorModel {
	m := EditorModel{
		nameInput:        textinput.New(),
		descriptionInput: textarea.New(),
		help:             help.New(),
	}
	m.nameInput.SetValue(currentName)
	m.descriptionInput.SetValue(currentDescription)
	m.nameInput.Focus()
	return m
}

func (m EditorModel) Init() tea.Cmd {
	return nil
}

func (m EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, editorKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, editorKeys.Save):
			m.hasUpdate = true
			return m, tea.Quit
		case key.Matches(msg, editorKeys.SwitchInput):
			if m.nameInput.Focused() {
				m.nameInput.Blur()
				return m, m.descriptionInput.Focus()
			}
			if m.descriptionInput.Focused() {
				m.descriptionInput.Blur()
				return m, m.nameInput.Focus()
			}
			return m, m.nameInput.Focus()
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m EditorModel) View() string {
	return fmt.Sprintf(
		"Name\n%s\n\nDescription\n%s\n\n%s",
		m.nameInput.View(),
		m.descriptionInput.View(),
		m.help.View(editorKeys),
	)
}

func (m EditorModel) Name() string {
	return m.nameInput.Value()
}

func (m EditorModel) Description() string {
	return m.descriptionInput.Value()
}

func (m EditorModel) HasUpdate() bool {
	return m.hasUpdate
}

func (m *EditorModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 2)
	m.nameInput, cmds[0] = m.nameInput.Update(msg)
	m.descriptionInput, cmds[1] = m.descriptionInput.Update(msg)
	return tea.Batch(cmds...)
}

var editorKeys = EditorKeys{
	Save: key.NewBinding(
		key.WithKeys("ctrl+s", "cmd+s"),
		key.WithHelp("ctrl+s", "Save"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "cmd+c", "esc"),
		key.WithHelp("ctrl+c/esc", "Cancel"),
	),
	SwitchInput: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Switch input"),
	),
}

func (k EditorKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Save,
		k.Quit,
		k.SwitchInput,
	}
}

func (k EditorKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
