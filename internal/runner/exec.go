package runner

import (
	"errors"
	"os"
	"os/exec"
)

func RunChild(child []string, env []string) error {
	if len(child) == 0 {
		return errors.New("missing child command")
	}

	cmd := exec.Command(child[0], child[1:]...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
