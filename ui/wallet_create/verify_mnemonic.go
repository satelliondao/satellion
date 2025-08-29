package wallet_create

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/ui/frame"
	"github.com/satelliondao/satellion/ui/frame/page"
)

const wordCount = 3

type verifyState struct {
	ctx      *frame.AppContext
	mnemonic *mnemonic.Mnemonic
	inputs   []textinput.Model
	focus    int
	indices  []int
	err      string
}

func NewVerify(ctx *frame.AppContext) frame.Page {
	m := verifyState{ctx: ctx}
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

func (m verifyState) Init() tea.Cmd {
	return textinput.Blink
}

func (m verifyState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, frame.Navigate(page.Passphrase)
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

func (m verifyState) isMnemonicValid() bool {
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

func (m verifyState) View() string {
	v := frame.NewViewBuilder()
	if m.mnemonic == nil {
		return "Verify your mnemonic\n\nMnemonic not found. Press Esc to go back."
	}
	v.Line("Verify your mnemonic")
	v.Line("Type the requested words to confirm.")
	for i := 0; i < 3; i++ {
		v.Line(m.inputs[i].View())
	}
	if m.err != "" {
		v.Line(m.err)
	}
	if strings.TrimSpace(m.inputs[2].Value()) != "" {
		v.Line("Press Enter to continue")
	}
	return v.Build()
}
