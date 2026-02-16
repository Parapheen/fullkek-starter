package scaffold

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Parapheen/fullkek-starter/internal/stacks"
)

func TestGenerateInitializesGitRepository(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not available in PATH")
	}

	stack, err := stacks.Compose(stacks.DefaultSelection())
	if err != nil {
		t.Fatalf("compose stack: %v", err)
	}

	destination := filepath.Join(t.TempDir(), "my-app")
	generator := DefaultGenerator()

	err = generator.Generate(context.Background(), Options{
		AppName:     "my-app",
		ModulePath:  "example.com/my-app",
		Destination: destination,
		Stack:       stack,
	})
	if err != nil {
		t.Fatalf("generate project: %v", err)
	}

	info, err := os.Stat(filepath.Join(destination, ".git"))
	if err != nil {
		t.Fatalf("stat .git directory: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected .git to be a directory")
	}
}
