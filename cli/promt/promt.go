package prompt

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/satelliondao/satellion/bip39"
	"github.com/satelliondao/satellion/cli/palette"
	"golang.org/x/term"
)

var validator = bip39.NewValidator()

// ProvideMnemonic prompts for a 12-word BIP39 mnemonic and returns it as bytes.
func ProvideMnemonic() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		palette.Question.Println("Enter 12-word BIP39 mnemonic: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if err := validator.Validate(line); err != nil {
			palette.Error.Println("Invalid mnemonic:", err)
			continue
		}
		normalized := validator.Normalize(line)
		return []byte(strings.Join(normalized, " ")), nil
	}
}

// ProvidePrivPassphrase is used to prompt for the private passphrase which
// maybe required during upgrades.
func ProvidePrivPassphrase() ([]byte, error) {
	prompt := "Enter the private passphrase of your wallet: "
	for {
		palette.Question.Print(prompt)
		pass, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return nil, err
		}
		fmt.Print("\n")
		pass = bytes.TrimSpace(pass)
		if len(pass) == 0 {
			continue
		}

		return pass, nil
	}
}

// promptList prompts the user with the given prefix, list of valid responses,
// and default list entry to use.  The function will repeat the prompt to the
// user until they enter a valid response.
func promptList(reader *bufio.Reader, prefix string, validResponses []string, defaultEntry string) (string, error) {
	// Setup the prompt according to the parameters.
	validStrings := strings.Join(validResponses, "/")
	var prompt string
	if defaultEntry != "" {
		prompt = fmt.Sprintf("%s (%s) [%s]: ", prefix, validStrings,
			defaultEntry)
	} else {
		prompt = fmt.Sprintf("%s (%s): ", prefix, validStrings)
	}

	// Prompt the user until one of the valid responses is given.
	for {
		fmt.Print(prompt)
		reply, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		reply = strings.TrimSpace(strings.ToLower(reply))
		if reply == "" {
			reply = defaultEntry
		}

		for _, validResponse := range validResponses {
			if reply == validResponse {
				return reply, nil
			}
		}
	}
}

// promptListBool prompts the user for a boolean (yes/no) with the given prefix.
// The function will repeat the prompt to the user until they enter a valid
// response.
func promptListBool(reader *bufio.Reader, prefix string,
	defaultEntry string) (bool, error) { // nolint:unparam

	// Setup the valid responses.
	valid := []string{"n", "no", "y", "yes"}
	response, err := promptList(reader, prefix, valid, defaultEntry)
	if err != nil {
		return false, err
	}
	return response == "yes" || response == "y", nil
}

// promptPass prompts the user for a passphrase with the given prefix.  The
// function will ask the user to confirm the passphrase and will repeat the
// prompts until they enter a matching response.
func promptPass(_ *bufio.Reader, prefix string, confirm bool) ([]byte, error) {
	// Prompt the user until they enter a passphrase.
	prompt := fmt.Sprintf("%s: ", prefix)
	for {
		fmt.Print(prompt)
		pass, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return nil, err
		}
		fmt.Print("\n")
		pass = bytes.TrimSpace(pass)
		if len(pass) == 0 {
			continue
		}

		if !confirm {
			return pass, nil
		}

		fmt.Print("Confirm passphrase: ")
		confirm, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return nil, err
		}
		fmt.Print("\n")
		confirm = bytes.TrimSpace(confirm)
		if !bytes.Equal(pass, confirm) {
			fmt.Println("The entered passphrases do not match")
			continue
		}

		return pass, nil
	}
}

func VerifyMnemonicSaved(mnemonic string) bool {
	reader := bufio.NewReader(os.Stdin)
	words := validator.Normalize(mnemonic)
	if len(words) != 12 {
		palette.Error.Println("Seed phrase must be 12 words")
		return false
	}
	fmt.Print("\033[2J\033[3J\033[H")
	attempts := 0
	for {
		asked := map[int]struct{}{}
		indices := make([]int, 0, 3)
		for len(indices) < 3 {
			i := rand.Intn(12) + 1
			if _, ok := asked[i]; ok {
				continue
			}
			asked[i] = struct{}{}
			indices = append(indices, i)
		}
		palette.Question.Println("Verify your seed phrase. Enter the requested words exactly as saved.")
		correct := true
		for _, idx := range indices {
			palette.Question.Printf("Word #%d: ", idx)
			line, err := reader.ReadString('\n')
			if err != nil {
				return false
			}
			entered := strings.ToLower(strings.TrimSpace(line))
			if entered != words[idx-1] {
				correct = false
			}
		}
		if correct {
			palette.Success.Println("Seed phrase verification passed.")
			return true
		}
		attempts++
		if attempts >= 3 {
			palette.Error.Println("Verification failed. Please back up your seed phrase and try again.")
			return false
		}
		palette.Warning.Println("One or more answers were incorrect. Let's try again.")
	}
}