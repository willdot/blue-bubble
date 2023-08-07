package input

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Will prompt the user to input some text. If true is returned it means the user quit instead of selecting an option.
// The provided message can be used to display text to the user as part of the prompt, but is optional.
func PromptUserForText(message string) (string, bool) {
	model := initialTextInputModel(message)
	p := tea.NewProgram(model)
	res, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	if m, ok := res.(textInputModel); ok {
		if m.quit {
			return "", true
		}

		return m.textInput.Value(), false
	}

	return "", false
}

type textInputModel struct {
	textInput textinput.Model
	err       error
	quit      bool
}

func initialTextInputModel(question string) textInputModel {
	ti := textinput.New()
	ti.Placeholder = question
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return textInputModel{
		textInput: ti,
		err:       nil,
	}
}

func (m textInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quit = true
			return m, tea.Quit
		case tea.KeyEnter:
			return m, tea.Quit
		}
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m textInputModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	)
}
