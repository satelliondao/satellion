package framework

import (
	"fmt"

	"github.com/fatih/color"
)

type ViewBuilder struct {
	v        string
	helpText string
	errText  string
	quitText bool
}

func View() *ViewBuilder {
	b := &ViewBuilder{
		v: "",
	}
	b.withLogo()
	return b
}

func (b *ViewBuilder) HideLogo() *ViewBuilder {
	b.v = ""
	return b
}

func (b *ViewBuilder) withLogo() *ViewBuilder {
	title := "SATELLION WALLET"
	b.v = fmt.Sprintf("%s\n", color.New(color.Bold).Sprintln(title))
	return b
}

func (b *ViewBuilder) L(format string, args ...interface{}) *ViewBuilder {
	b.v += fmt.Sprintf(format+"\n", args...)
	return b
}

func (b *ViewBuilder) Warn(format string, args ...interface{}) *ViewBuilder {
	b.v += color.New(color.FgYellow).Sprintf(format+"\n", args...)
	return b
}

func (b *ViewBuilder) Help(s string) *ViewBuilder {
	b.helpText += fmt.Sprintln(s)
	return b
}

func (b *ViewBuilder) QuitHint() *ViewBuilder {
	b.quitText = true
	return b
}

func (b *ViewBuilder) Err(err string) *ViewBuilder {
	b.errText = err
	return b
}

func (b *ViewBuilder) Build() string {
	v := b.v
	if b.errText != "" {
		v += color.New(color.FgHiRed).Sprintln(b.errText)
	}
	if b.helpText != "" {
		v += color.New(color.FgHiBlack).Sprintln(b.helpText)
	}
	if b.quitText {
		v += color.New(color.FgHiBlack).Sprintln("Esc to go home. Ctrl+C to exit")
	}
	return v
}
