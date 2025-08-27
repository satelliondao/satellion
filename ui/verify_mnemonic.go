package ui

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/mnemonic"
)

const wordCount = 3

type verifyMnemonicModel struct {
	ctx      *AppContext
	mnemonic *mnemonic.Mnemonic
	inputs   []textinput.Model
	focus    int
	indices  []int
	err      string
}

func NewVerifyMnemonic(ctx *AppContext) Page {
	m := verifyMnemonicModel{ctx: ctx}
	if ctx != nil {
		m.mnemonic = ctx.TempMnemonic
	}
	if m.mnemonic != nil {
		rand.Seed(time.Now().UnixNano())
		perm := rand.Perm(len(m.mnemonic.Words))
		m.indices = []int{perm[0], perm[1], perm[2]}
		sort.Ints(m.indices)
		m.inputs = make([]textinput.Model, wordCount)
		for i := 0; i < wordCount; i++ {
			in := textinput.New()
			in.Placeholder = fmt.Sprintf("Word #%d", m.indices[i]+1)
			in.CharLimit = 50
			in.Width = 24
			m.inputs[i] = in
		}
		m.inputs[0].Focus()
		m.focus = 0
	}
	return m
}

func (m verifyMnemonicModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m verifyMnemonicModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.mnemonic == nil || len(m.inputs) == 0 {
				return m, nil
			}
			last := len(m.inputs) - 1
			if m.focus < last {
				if strings.TrimSpace(m.inputs[m.focus].Value()) == "" {
					return m, nil
				}
				m.inputs[m.focus].Blur()
				m.focus++
				m.inputs[m.focus].Focus()
				return m, nil
			}
			// if !m.isMnemonicValid() {
			// 	m.err = "Words do not match. Try again."
			// 	return m, nil
			// }
			if err := m.ctx.Router.AddWallet(m.ctx.TempWalletName, m.mnemonic); err != nil {
				m.err = err.Error()
				return m, nil
			}
			m.ctx.TempWalletName = ""
			m.ctx.TempMnemonic = nil
			return m, Navigate(config.HomePage)
		}
	}
	if len(m.inputs) > 0 {
		cur := m.inputs[m.focus]
		cur, cmd = cur.Update(msg)
		m.inputs[m.focus] = cur
		return m, cmd
	}
	return m, cmd
}

func (m verifyMnemonicModel) isMnemonicValid() bool {
	valid := true
	for i := 0; i < wordCount; i++ {
		want := strings.ToLower(m.mnemonic.Words[m.indices[i]])
		got := strings.ToLower(strings.TrimSpace(m.inputs[i].Value()))
		if got != want {
			valid = false
			break
		}
	}
	return valid
}

func (m verifyMnemonicModel) View() string {
	if m.mnemonic == nil {
		return "Verify your mnemonic\n\nMnemonic not found. Press Esc to go back."
	}
	view := "Verify your mnemonic\n\nType the requested words to confirm.\n\n"
	for i := 0; i < 3; i++ {
		view += fmt.Sprintf("%s\n", m.inputs[i].View())
	}
	if m.err != "" {
		view += fmt.Sprintf("\n%s\n", m.err)
	}
	if strings.TrimSpace(m.inputs[2].Value()) != "" {
		view += "\nPress Enter to save"
	}
	return view
}
