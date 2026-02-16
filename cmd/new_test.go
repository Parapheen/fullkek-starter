package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewCommandIncludesAuthFlag(t *testing.T) {
	t.Parallel()

	root := RootCommand()
	newCmd, _, err := root.Find([]string{"new"})
	if err != nil {
		t.Fatalf("expected to find new command, got: %v", err)
	}

	flag := newCmd.Flags().Lookup("auth")
	if flag == nil {
		t.Fatal("expected --auth flag to be registered")
	}
}

func TestNewNoUIRequiresAppName(t *testing.T) {
	t.Parallel()

	root := RootCommand()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"new", "--no-ui"})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "app name required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRootHelpDoesNotRenderBannerForNonTTY(t *testing.T) {
	t.Parallel()

	root := RootCommand()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("expected help to succeed, got: %v", err)
	}

	if strings.Contains(out.String(), "███████") {
		t.Fatalf("expected no banner in non-interactive output, got: %s", out.String())
	}
}
