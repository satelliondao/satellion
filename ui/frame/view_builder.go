package frame

import (
	"fmt"

	"github.com/fatih/color"
)

type ViewBuilder struct {
	v        string
	helpText string
	errText  string
}

func NewViewBuilder() *ViewBuilder {
	b := &ViewBuilder{
		v: "",
	}
	b.withLogo()
	return b
}

func (b *ViewBuilder) withLogo() *ViewBuilder {
	title := "SATELLION WALLET"
	b.v = color.New(color.BgHiWhite).Sprintf("%s\n", color.New(color.Bold).Sprintf("%s", title))
	b.v += "\n"
	return b
}

func (b *ViewBuilder) Line(s string) *ViewBuilder {
	b.v += fmt.Sprintf("%s\n", s)
	return b
}

func (b *ViewBuilder) WithHelpText(s string) *ViewBuilder {
	b.helpText = s
	return b
}

func (b *ViewBuilder) WithQuitText() *ViewBuilder {
	b.v += color.New(color.FgWhite).Sprintf("\nctrl+c to exit")
	return b
}

func (b *ViewBuilder) WithErrText(s string) *ViewBuilder {
	b.errText = s
	return b
}

func (b *ViewBuilder) Build() string {
	v := b.v
	if b.errText != "" {
		v += color.New(color.FgHiRed).Sprintf("\n%s", b.errText)
	}
	if b.helpText != "" {
		v += color.New(color.FgHiYellow).Sprintf("\n%s", b.helpText)
	}
	return v
}
