package prompt

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Prompter struct {
	In        io.Reader
	Out       io.Writer
	AssumeYes bool
}

func (p Prompter) YesNo(question string, def bool) bool {
	if p.AssumeYes {
		return def
	}
	suffix := " [y/N]: "
	if def {
		suffix = " [Y/n]: "
	}
	fmt.Fprint(p.Out, question+suffix)
	answer := p.readLine()
	if answer == "" {
		return def
	}
	answer = strings.ToLower(answer)
	return answer == "y" || answer == "yes"
}

func (p Prompter) Value(question, def string) string {
	fmt.Fprintf(p.Out, "%s [%s]: ", question, def)
	answer := p.readLine()
	if answer == "" {
		return def
	}
	return answer
}

func (p Prompter) Danger(context, token string) error {
	fmt.Fprintf(p.Out, "Dangerous operation: %s\nType %s to continue: ", context, token)
	if p.readLine() != token {
		return fmt.Errorf("confirmation token did not match")
	}
	return nil
}

func (p Prompter) readLine() string {
	scanner := bufio.NewScanner(p.In)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}
