package verify_mnemonic

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/ui/page"
	"github.com/satelliondao/satellion/ui/staff"
)

const wordCount = 3

type State struct {
	ctx      *staff.AppContext
	mnemonic *mnemonic.Mnemonic
	inputs   []textinput.Model
	focus    int
	indices  []int
	err      string
}

func New(ctx *staff.AppContext) staff.Page {
	m := State{ctx: ctx}
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

func (m State) Init() tea.Cmd {
	return textinput.Blink
}

func (m State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, staff.Navigate(page.Passphrase)
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

func (m State) View() string {
	v := staff.NewViewBuilder()
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
