package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/asimmittal/key-env/internal/envfile"
	"github.com/asimmittal/key-env/internal/runner"
	"github.com/asimmittal/key-env/internal/secrets"
	"github.com/asimmittal/key-env/internal/vault"
	"github.com/asimmittal/key-env/internal/vault/keepassxc"
	"github.com/asimmittal/key-env/internal/vault/onepassword"
)

type Config struct {
	EnvPath     string
	SecretsPath string
	Password    string
	Verbose     bool
	Child       []string
}

func Run(args []string) error {
	if len(args) == 0 || args[0] != "run" {
		return errors.New("usage: key-env run --env <file> [--secrets <kdbx> --password <value>] -- <child command>")
	}

	cfg, err := parseRunArgs(args[1:])
	if err != nil {
		return err
	}

	parsed, err := envfile.ParseFile(cfg.EnvPath)
	if err != nil {
		return err
	}

	if hasType(parsed, envfile.TypeKP) {
		if cfg.SecretsPath == "" {
			return errors.New("env file contains kp:// references; --secrets is required")
		}
		if cfg.Password == "" {
			return errors.New("env file contains kp:// references; --password is required")
		}
		if _, err := os.Stat(cfg.SecretsPath); err != nil {
			return fmt.Errorf("invalid --secrets path %q: %w", cfg.SecretsPath, err)
		}
	}

	providers := map[string]vault.Provider{
		"kp": keepassxc.New(cfg.SecretsPath, cfg.Password),
		"op": onepassword.New(),
	}
	loader := secrets.NewLoader(providers, cfg.Verbose)
	loaded, err := loader.Load(parsed)
	if err != nil {
		return err
	}

	if cfg.Verbose {
		printSummary(loaded)
	}

	finalEnv := secrets.MergeWithCurrentEnv(loaded, os.Environ())
	return runner.RunChild(cfg.Child, finalEnv)
}

func parseRunArgs(args []string) (Config, error) {
	delim := -1
	for i, a := range args {
		if a == "--" {
			delim = i
			break
		}
	}
	if delim == -1 {
		return Config{}, errors.New("missing child command separator `--`; usage: key-env run --env <file> --secrets <kdbx> -- <child command>")
	}

	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	var cfg Config
	fs.StringVar(&cfg.EnvPath, "env", "", "path to env file")
	fs.StringVar(&cfg.SecretsPath, "secrets", "", "path to KeepassXC .kdbx file")
	fs.StringVar(&cfg.Password, "password", "", "vault password (prefer stdin/file in production)")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "print detailed logging")

	if err := fs.Parse(args[:delim]); err != nil {
		return Config{}, err
	}

	cfg.Child = args[delim+1:]
	if cfg.EnvPath == "" {
		return Config{}, errors.New("missing required input: --env")
	}
	if len(cfg.Child) == 0 {
		return Config{}, errors.New("missing required child command after `--`")
	}
	if _, err := os.Stat(cfg.EnvPath); err != nil {
		return Config{}, fmt.Errorf("invalid --env path %q: %w", cfg.EnvPath, err)
	}

	return cfg, nil
}

func hasType(parsed []envfile.ParsedVar, t string) bool {
	for _, item := range parsed {
		if item.Type == t {
			return true
		}
	}
	return false
}

func printSummary(loaded []secrets.LoadedVar) {
	total := len(loaded)
	decrypted := 0
	kpCount := 0
	opCount := 0
	for _, v := range loaded {
		switch v.Type {
		case envfile.TypeKP:
			decrypted++
			kpCount++
		case envfile.TypeOP:
			decrypted++
			opCount++
		}
	}

	pct := func(n int) float64 {
		if decrypted == 0 {
			return 0
		}
		return float64(n) / float64(decrypted) * 100
	}

	fmt.Fprintf(os.Stderr, "\nLoaded %d env vars (%d decrypted)\n", total, decrypted)
	fmt.Fprintf(os.Stderr, "  kp:// %5.1f%%\n", pct(kpCount))
	fmt.Fprintf(os.Stderr, "  op:// %5.1f%%\n", pct(opCount))
	fmt.Fprintln(os.Stderr, "\n--- key-env complete ---")
}
