package staff

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

// Choice represents a selectable option
type Choice struct {
	Label string
	Value interface{}
}

// ChoiceResult represents the result of handling a key message
type ChoiceResult struct {
	Action   string  // "navigation", "selection", "none"
	Selected *Choice // the selected choice (nil for navigation)
	Consumed bool    // whether the key message was consumed by the selector
}

// ChoiceResult constants
const (
	ActionNavigation = "navigation"
	ActionSelection  = "selection"
	ActionNone       = "none"
)

// ChoiceSelector handles selection from a list of choices
type ChoiceSelector struct {
	choices []Choice
	cursor  int
}

// NewChoiceSelector creates a new choice selector with the given choices
func NewChoiceSelector(choices []Choice) *ChoiceSelector {
	return &ChoiceSelector{
		choices: choices,
		cursor:  0,
	}
}

// NewChoiceSelectorFromStrings creates a choice selector from string labels
func NewChoiceSelectorFromStrings(labels []string) *ChoiceSelector {
	choices := make([]Choice, len(labels))
	for i, label := range labels {
		choices[i] = Choice{Label: label, Value: i}
	}
	return NewChoiceSelector(choices)
}

// SetChoices updates the choices in the selector
func (s *ChoiceSelector) SetChoices(choices []Choice) {
	s.choices = choices
	if s.cursor >= len(choices) {
		s.cursor = 0
	}
}

// SetChoicesFromStrings updates the choices from string labels
func (s *ChoiceSelector) SetChoicesFromStrings(labels []string) {
	s.SetChoices(make([]Choice, len(labels)))
	for i, label := range labels {
		s.choices[i] = Choice{Label: label, Value: i}
	}
	if s.cursor >= len(labels) {
		s.cursor = 0
	}
}

// Next moves the cursor to the next choice
func (s *ChoiceSelector) Next() {
	if len(s.choices) == 0 {
		return
	}
	if s.cursor < len(s.choices)-1 {
		s.cursor++
	} else {
		s.cursor = 0
	}
}

// Previous moves the cursor to the previous choice
func (s *ChoiceSelector) Previous() {
	if len(s.choices) == 0 {
		return
	}
	if s.cursor > 0 {
		s.cursor--
	} else {
		s.cursor = len(s.choices) - 1
	}
}

// Selected returns the currently selected choice
func (s *ChoiceSelector) Selected() *Choice {
	if len(s.choices) == 0 || s.cursor >= len(s.choices) {
		return nil
	}
	return &s.choices[s.cursor]
}

// SelectedIndex returns the index of the currently selected choice
func (s *ChoiceSelector) SelectedIndex() int {
	return s.cursor
}

// SelectIndex sets the cursor to a specific index
func (s *ChoiceSelector) SelectIndex(index int) {
	if index >= 0 && index < len(s.choices) {
		s.cursor = index
	}
}

// Update handles key messages and returns a ChoiceResult indicating what action was taken
func (s *ChoiceSelector) Update(msg tea.Msg) ChoiceResult {
	switch v := msg.(type) {
	case tea.KeyMsg:
		switch v.String() {
		case "up", "k":
			s.Previous()
			return ChoiceResult{
				Action:   ActionNavigation,
				Selected: nil,
				Consumed: true,
			}
		case "down", "j":
			s.Next()
			return ChoiceResult{
				Action:   ActionNavigation,
				Selected: nil,
				Consumed: true,
			}
		case "enter":
			selected := s.Selected()
			if selected != nil {
				return ChoiceResult{
					Action:   ActionSelection,
					Selected: selected,
					Consumed: true,
				}
			}
		}
	}

	return ChoiceResult{
		Action:   ActionNone,
		Selected: nil,
		Consumed: false,
	}
}

// Count returns the number of choices
func (s *ChoiceSelector) Count() int {
	return len(s.choices)
}

// IsEmpty returns true if there are no choices
func (s *ChoiceSelector) IsEmpty() bool {
	return len(s.choices) == 0
}

// Render generates the display string for the choice selector
func (s *ChoiceSelector) Render(highlightColor ...color.Attribute) string {
	var result string

	// Default to cyan highlight
	highlight := color.FgHiCyan
	if len(highlightColor) > 0 {
		highlight = highlightColor[0]
	}

	for i, choice := range s.choices {
		cursor := " "
		if s.cursor == i {
			cursor = color.New(highlight).Sprintf(">")
		}
		result += fmt.Sprintf("%s %s\n", cursor, choice.Label)
	}

	return result
}

// RenderWithPrefix generates the display string with a custom prefix
func (s *ChoiceSelector) RenderWithPrefix(prefix string, highlightColor ...color.Attribute) string {
	var result string

	// Default to cyan highlight
	highlight := color.FgHiCyan
	if len(highlightColor) > 0 {
		highlight = highlightColor[0]
	}

	for i, choice := range s.choices {
		cursor := " "
		if s.cursor == i {
			cursor = color.New(highlight).Sprintf(prefix)
		}
		result += fmt.Sprintf("%s %s\n", cursor, choice.Label)
	}

	return result
}

// GetChoiceByValue finds a choice by its value
func (s *ChoiceSelector) GetChoiceByValue(value interface{}) *Choice {
	for _, choice := range s.choices {
		if choice.Value == value {
			return &choice
		}
	}
	return nil
}

// GetChoiceByIndex returns a choice at the specified index
func (s *ChoiceSelector) GetChoiceByIndex(index int) *Choice {
	if index < 0 || index >= len(s.choices) {
		return nil
	}
	return &s.choices[index]
}
