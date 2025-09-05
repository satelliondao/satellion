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
	"github.com/satelliondao/satellion/ui/framework"
	"github.com/satelliondao/satellion/ui/router"
)

const wordCount = 3

type State struct {
	ctx        *framework.AppContext
	mnemonic   *mnemonic.Mnemonic
	inputs     []textinput.Model
	focus      int
	indices    []int
	err        string
	walletName string
}

func New(ctx *framework.AppContext, params interface{}) framework.Page {
	m := State{ctx: ctx}
	if props, ok := params.(*router.VerifyMnemonicProps); ok {
		m.walletName = props.WalletName
		m.mnemonic = props.Mnemonic
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
			return m, router.Passphrase(m.walletName, m.mnemonic)
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
	v := framework.View()
	if m.mnemonic == nil {
		return "Verify your mnemonic\n\nMnemonic not found. Press Esc to go back."
	}
	v.L("Verify your mnemonic")
	for i := 0; i < 3; i++ {
		v.L(m.inputs[i].View())
	}
	if m.err != "" {
		v.L(m.err)
	}
	if strings.TrimSpace(m.inputs[2].Value()) != "" {
		v.Help("Press Enter to continue")
	}
	return v.Build()
}
