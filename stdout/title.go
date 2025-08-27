package stdout

import (
	"strings"
)

func Title() string {
	text := "SATELLION WALLET"
	spaced := strings.Join(strings.Split(text, ""), " ")
	var b strings.Builder
	for _, ch := range spaced {
		b.WriteString(string(ch))
	}
	return b.String()
}
