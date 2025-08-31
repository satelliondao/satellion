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

func NewViewBuilder() *ViewBuilder {
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

func (b *ViewBuilder) Line(s string) *ViewBuilder {
	b.v += fmt.Sprintf("%s\n", s)
	return b
}

func (b *ViewBuilder) WithHelpText(s string) *ViewBuilder {
	b.helpText += fmt.Sprintln(s)
	return b
}

func (b *ViewBuilder) WithQuitText() *ViewBuilder {
	b.quitText = true
	return b
}

func (b *ViewBuilder) WithErrText(s string) *ViewBuilder {
	b.errText = s
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
