package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRootHelpIncludesCompletionAndDocs(t *testing.T) {
	var out, err bytes.Buffer
	rt := NewRuntime(&out, &err, strings.NewReader(""))
	if e := Execute(context.Background(), []string{"--help"}, rt); e != nil {
		t.Fatal(e)
	}
	text := out.String()
	for _, want := range []string{"completion", "docs", "client", "apply"} {
		if !strings.Contains(text, want) {
			t.Fatalf("help missing %q:\n%s", want, text)
		}
	}
}
