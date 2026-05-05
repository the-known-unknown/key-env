# key-env

**Keep secrets encrypted in your password manager -- not in your `.env` files.**

[![CI](https://github.com/the-known-unknown/key-env/actions/workflows/ci.yml/badge.svg)](https://github.com/the-known-unknown/key-env/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/the-known-unknown/key-env)](https://github.com/the-known-unknown/key-env/releases)

---

Your `.env` files today:

```env
CLIENT_SECRET="Jcg5TfdI9X0zHaU03Qx9bGb0rphYh0xIebtpFPTcRT"
CLIENT_NAME="test"
```

Your `.env` files with `key-env`:

```env
CLIENT_SECRET="kp://my-api/Password"
CLIENT_NAME="kp://my-api/Username"
```

The actual values stay encrypted in your password manager. `key-env` resolves them at runtime and passes them to your app. Nothing is written to disk in plaintext.

---

## The problem

`.env` files store API keys, database credentials, and OAuth secrets in plaintext. They get committed to repos, shared over Slack, and copied between machines. Every duplicate is another place a secret can leak.

Agentic coding tools make this worse. Copilot, Cursor, and Claude read and write code on your behalf -- which means they also read your `.env` files. Every AI agent that touches your project increases the surface area of exposure. The result: plaintext secrets scattered across machines, repos, and chat histories, accessible to tools and processes you don't fully control.

---

## How it works

```mermaid
flowchart LR
    EnvFile[".env file"] --> Parser["envfile parser"]
    Parser -->|"kp://"| KP["keepassxc provider"]
    Parser -->|"op://"| OP["onepassword provider"]
    Parser -->|"plain"| Plain["pass-through"]
    KP -->|"shells out"| KPCli["keepassxc-cli show\n(uses --password)"]
    OP -->|"shells out"| OPCli["op read\n(uses ambient auth)"]
    KPCli --> Loader["secrets loader"]
    OPCli --> Loader
    Plain --> Loader
    Loader --> Runner["child process"]
```

When you run your application through `key-env`, it:

1. Reads your `.env` file and identifies secret references (e.g. `kp://...`, `op://...`)
2. Routes each reference to the right provider, which shells out to the corresponding password manager CLI
3. Injects the decrypted values into the environment
4. Launches your application with the fully hydrated environment

Plain values are passed through as-is. You can mix `kp://`, `op://`, and plain values in the same `.env` file.

---

## Set up your local password manager

`key-env` works with [KeePassXC](https://keepassxc.org/) (`kp://`) and [1Password](https://1password.com/) (`op://`). You can use either one, or both at the same time. The corresponding CLI tools (`keepassxc-cli`, `op`) need to be on your `PATH`.

### KeePassXC

[KeePassXC](https://keepassxc.org/) is a free, open-source, offline password manager.

**macOS:**

```sh
brew install keepassxc
```

**Linux:**

```sh
# Ubuntu / Debian
sudo apt install keepassxc

# Fedora
sudo dnf install keepassxc
```

You can also download installers directly from [keepassxc.org](https://keepassxc.org/).

Verify the CLI is available:

```sh
keepassxc-cli --version
```

### 1Password

[1Password](https://1password.com/) is a hosted password manager. `key-env` shells out to its CLI (`op`) to resolve `op://` references.

**Install:**

```sh
brew install 1password-cli
```

For other platforms, see the [official downloads page](https://1password.com/downloads/command-line/).

**Recommended setup -- desktop app integration:**

This is the simplest path for local development. The 1Password desktop app handles authentication via Touch ID / Windows Hello / system auth on demand, so `key-env` doesn't need to manage tokens or sessions.

1. Open the 1Password desktop app
2. Settings → Security → enable Touch ID / Windows Hello / system auth
3. Settings → Developer → toggle on **"Integrate with 1Password CLI"**

Verify the CLI works (you may see a system auth prompt the first time):

```sh
op vault list
```

> **Headless or CI environments:** set `OP_SERVICE_ACCOUNT_TOKEN` in your shell before running `key-env`. The `op` subprocess inherits the variable automatically. See the [official service accounts docs](https://developer.1password.com/docs/service-accounts/get-started) for setup.

---

## Quick start

**1. Download and install the binary**

Grab the latest release for your platform from the [releases page](https://github.com/the-known-unknown/key-env/releases).

```sh
# Example for macOS Apple Silicon — adjust the URL for your platform
curl -LO https://github.com/the-known-unknown/key-env/releases/latest/download/key-env_darwin_arm64.tar.gz
tar xzf key-env_darwin_arm64.tar.gz
```

Move it somewhere on your `PATH`:

```sh
sudo mv key-env /usr/local/bin/
```

Or, if you prefer not to use `sudo`, drop it in a local bin directory:

```sh
mkdir -p ~/.local/bin
mv key-env ~/.local/bin/

# Add to your PATH if it isn't already (add this line to your ~/.zshrc or ~/.bashrc)
export PATH="$HOME/.local/bin:$PATH"
```

Verify it works:

```sh
key-env --version
```

**2. Try it with the included KeePassXC sample**

The repo includes a sample KeePassXC database and `.env` files so you can test immediately:

```sh
git clone https://github.com/the-known-unknown/key-env.git
cd key-env
```

```sh
key-env run \
  --env test/.env.kp \
  --secrets test/keepass-sample-db.kdbx \
  --password '4jFU%i*+Q2qdpFgoHJGK' \
  -- sh -c 'echo $TEST_CLIENT_SECRET'
```

```
Jcg5TfdI9X0zHaU03Qx9bGb0rphYh0xIebtpFPTcRT
```

The secret was resolved from the encrypted KeePassXC database and injected into the shell command.

> **Trying the 1Password flow?** You'll need a 1Password account with the desktop app integration enabled and an item in one of your vaults that the included `test/.env.op` references can resolve. See the [1Password setup section](#1password) above and edit `test/.env.op` to point at your own item.

---

## Examples

The `test/` directory includes sample scripts that print the resolved secrets and run a progress bar, along with two `.env` files that reference each backend:

- `test/.env.kp` -- references the bundled KeePassXC sample database
- `test/.env.op` -- references a 1Password item (edit to point at one of your own)

### KeePassXC examples

**Shell:**

```sh
key-env run \
  --env test/.env.kp \
  --secrets test/keepass-sample-db.kdbx \
  --password '4jFU%i*+Q2qdpFgoHJGK' \
  -- sh -c 'echo "CLIENT_SECRET=$TEST_CLIENT_SECRET"; echo "CLIENT_NAME=$TEST_CLIENT_NAME"'
```

```
CLIENT_SECRET=Jcg5TfdI9X0zHaU03Qx9bGb0rphYh0xIebtpFPTcRT
CLIENT_NAME=test
```

**Python:**

```sh
key-env run \
  --env test/.env.kp \
  --secrets test/keepass-sample-db.kdbx \
  --password '4jFU%i*+Q2qdpFgoHJGK' \
  -- python3 test/sample.py
```

```
Starting Python sample...
Env vars:
TEST_CLIENT_SECRET: Jcg5TfdI9X0zHaU03Qx9bGb0rphYh0xIebtpFPTcRT
TEST_CLIENT_NAME: test

Doing some work...
[██████████] 100%
✔ Done!
```

**Node.js:**

```sh
key-env run \
  --env test/.env.kp \
  --secrets test/keepass-sample-db.kdbx \
  --password '4jFU%i*+Q2qdpFgoHJGK' \
  -- node test/sample.js
```

```
Staring Node.js sample...
Env vars:
TEST_CLIENT_SECRET: Jcg5TfdI9X0zHaU03Qx9bGb0rphYh0xIebtpFPTcRT
TEST_CLIENT_NAME: test
NODE_ENV:

Doing some work...
[██████████] 100%
✔ Node.js sample done!
```

**Node.js with inline env vars:**

`key-env` exec's the child command directly, not through a shell, so the `VAR=value command` syntax won't work after `--`. Use `env` as a wrapper:

```sh
key-env run \
  --env test/.env.kp \
  --secrets test/keepass-sample-db.kdbx \
  --password '4jFU%i*+Q2qdpFgoHJGK' \
  -- env NODE_ENV=development node test/sample.js
```

```
Staring Node.js sample...
Env vars:
TEST_CLIENT_SECRET: Jcg5TfdI9X0zHaU03Qx9bGb0rphYh0xIebtpFPTcRT
TEST_CLIENT_NAME: test
NODE_ENV: development

Doing some work...
[██████████] 100%
✔ Node.js sample done!
```

### 1Password examples

`test/.env.op` references an item in your 1Password account. Edit it to point at a real item in one of your vaults, for example:

```env
TEST_CLIENT_NAME="op://Employee/test_key_env/username"
TEST_CLIENT_SECRET="op://Employee/test_key_env/password"
```

For env files that only contain `op://` references, you don't need to pass `--secrets` or `--password`. With desktop app integration enabled, the `op` CLI authenticates each call via Touch ID / system auth.

**Shell:**

```sh
key-env run \
  --env test/.env.op \
  -- sh -c 'echo "CLIENT_SECRET=$TEST_CLIENT_SECRET"; echo "CLIENT_NAME=$TEST_CLIENT_NAME"'
```

```
CLIENT_SECRET=qwerty12345
CLIENT_NAME=asimmittal
```

> Note, this example only works if you have 1Password cli already setup with secrets at `Employee/test_key_env/username` in your valult. If you have some other values in your local 1Password cli, consider using those to test.

**Node.js:**

```sh
key-env run \
  --env test/.env.op \
  -- node test/sample.js
```

```
Staring Node.js sample...
Env vars:
TEST_CLIENT_SECRET: qwerty12345
TEST_CLIENT_NAME: asimmittal
NODE_ENV:

Doing some work...
[██████████] 100%
✔ Node.js sample done!
```

> **Note:** Each `op://` reference triggers a separate `op read` call, so a `.env` file with two `op://` entries takes about twice as long to resolve as one. With desktop app integration, the first call in a session may also include a Touch ID / system auth prompt.

### Mixing KeePassXC and 1Password

You can mix both backends in the same `.env` file. `key-env` will route each reference to the right provider automatically:

```env
DB_PASSWORD="kp://Services/Databases/Main/Password"
STRIPE_KEY="op://DevSecrets/Stripe/api_key"
APP_PORT="3000"
```

### Handling child process failures

If the child process exits with a non-zero code, `key-env` reports the failure and exits with code 1:

```sh
key-env run \
  --env test/.env.kp \
  --secrets test/keepass-sample-db.kdbx \
  --password '4jFU%i*+Q2qdpFgoHJGK' \
  -- node test/sample_fail.js
```

```
Starting Node.js failure sample...
Env vars:
TEST_CLIENT_SECRET: Jcg5TfdI9X0zHaU03Qx9bGb0rphYh0xIebtpFPTcRT
TEST_CLIENT_NAME: test

Simulating work before failure...
[█████] 100%

✘ Something went wrong — exiting with code 1
key-env: exit status 1
```

Secrets are still resolved and injected normally -- the failure happens in the child process, not in `key-env`. The `key-env: exit status N` line on stderr shows the child's original exit code.

---

## Install

**GitHub Releases:**

Download the latest binary for your platform from the [releases page](https://github.com/the-known-unknown/key-env/releases). Binaries are available for macOS (amd64, arm64) and Linux (amd64, arm64).

**Homebrew (macOS/Linux):**

```sh
brew install the-known-unknown/tap/key-env
```

---

## Build from source

Requires Go 1.22+.

```sh
git clone https://github.com/the-known-unknown/key-env.git
cd key-env
make build
```

This runs the test suite and produces a `key-env` binary in the project root.

---

## Usage

### `.env` file format

Plain values work as-is:

```env
DB_NAME="userdb"
APP_PORT="3000"
```

Secret references use a vault prefix that points into your password manager. You can mix prefixes in the same `.env` file.

**KeePassXC (`kp://`):**

```env
DB_PASSWORD="kp://Services/Databases/Main/Password"
API_KEY="kp://Services/Stripe/SecretKey"
```

The format is `kp://<entry_path>/<credential>`, where:

- `entry_path` is the path to the entry in your KeePassXC database (e.g. `Services/Databases/Main`)
- `credential` is the field to retrieve (e.g. `Password`, `Username`)

**1Password (`op://`):**

```env
STRIPE_KEY="op://DevSecrets/Stripe/api_key"
DB_PASSWORD="op://DevSecrets/Database/password"
```

The format is `op://<vault>/<item>/<field>`, where:

- `vault` is the 1Password vault name (e.g. `DevSecrets`)
- `item` is the item name within that vault (e.g. `Stripe`)
- `field` is the field to retrieve (e.g. `api_key`, `password`, `username`)

See the [1Password secret references syntax docs](https://developer.1password.com/docs/cli/secrets-reference-syntax/) for advanced usage like query parameters.

### Command

```sh
key-env run \
  --env <path-to-env-file> \
  --secrets <path-to-kdbx-file> \
  --password '<vault-password>' \
  -- <your command>
```

| Flag         | Required                | Description                               |
| ------------ | ----------------------- | ----------------------------------------- |
| `--env`      | Yes                     | Path to your `.env` file                  |
| `--secrets`  | If `kp://` refs present | Path to your `.kdbx` KeePassXC database   |
| `--password` | If `kp://` refs present | Password to unlock the KeePassXC database |
| `--verbose`  | No                      | Print per-variable logs and summary stats |
| `--version`  | No                      | Print version and exit                    |
| `--help`     | No                      | Show help message                         |
| `--`         | Yes                     | Separator before the child command        |

For env files that only contain `op://` references, `--secrets` and `--password` can be omitted -- the `op` CLI handles its own authentication.

### Verbose mode

Pass `--verbose` to see which variables were decrypted and a summary of what was loaded:

```sh
key-env run \
  --verbose \
  --env test/.env.kp \
  --secrets test/keepass-sample-db.kdbx \
  --password '4jFU%i*+Q2qdpFgoHJGK' \
  -- sh -c 'echo $TEST_CLIENT_SECRET'
```

```
[key-env | kp]: ✔ TEST_CLIENT_SECRET
[key-env | kp]: ✔ TEST_CLIENT_NAME

Loaded 2 env vars (2 decrypted)
  kp:// 100.0%
  op://   0.0%

--- key-env complete ---
Jcg5TfdI9X0zHaU03Qx9bGb0rphYh0xIebtpFPTcRT
```

---

## Supported vaults

| Prefix  | Provider      | Status      |
| ------- | ------------- | ----------- |
| `kp://` | KeePassXC     | Implemented |
| `op://` | 1Password CLI | Implemented |

See [Set up your local password manager](#set-up-your-local-password-manager) for install instructions. Tested with KeePassXC 2.7.11 and 1Password CLI 2.32.0.

---

## Security notes

- The `--password` flag can expose your vault password in shell history and process listings. In production workflows, prefer passing it via stdin or a secure file.
- Your `.env` file still reveals metadata -- variable names and vault paths -- even though the actual secret values are encrypted. Be mindful when sharing or committing it.
- The decrypted secrets exist in the child process's environment for the duration of its execution. They are not written to disk.

---

## Development

Run tests:

```sh
go test ./...
```

Build:

```sh
make build
```

Clean:

```sh
make clean
```

---

## Architecture

```mermaid
flowchart TD
    CLI["CLI entry point\ncmd/keyenv/main.go"]
    RunCmd["Run command\ninternal/cli/run.go"]
    EnvParser["Env file parser\ninternal/envfile/parser.go"]
    SecretsLoader["Secrets loader\ninternal/secrets/manager.go"]
    VaultInterface["Provider interface\ninternal/vault/provider.go"]
    KeePassXC["KeePassXC provider\ninternal/vault/keepassxc/keepassxc.go"]
    OnePassword["1Password provider\ninternal/vault/onepassword/onepassword.go"]
    Runner["Child process runner\ninternal/runner/exec.go"]

    CLI --> RunCmd
    RunCmd --> EnvParser
    RunCmd --> SecretsLoader
    SecretsLoader --> VaultInterface
    VaultInterface --> KeePassXC
    VaultInterface --> OnePassword
    RunCmd --> Runner
```

**`cmd/keyenv/main.go`** -- Entry point. Parses top-level flags (`--version`, `--help`) and delegates to the `run` command.

**`internal/cli/run.go`** -- Orchestrates the entire flow: parses `run` flags, calls the env parser, wires up vault providers, runs the secrets loader, merges the resolved environment, and exec's the child process. Also prints the `--verbose` summary.

**`internal/envfile/parser.go`** -- Reads a `.env` file line by line and classifies each variable as `plain`, `kp` (KeePassXC), or `op` (1Password) based on its value prefix. Returns a slice of `ParsedVar` structs.

**`internal/secrets/manager.go`** -- Takes the parsed variables and resolves any secret references by dispatching to the appropriate vault provider. Plain values pass through unchanged. Also provides `MergeWithCurrentEnv` to layer the resolved vars on top of the current process environment.

**`internal/vault/provider.go`** -- Defines the `Provider` interface: a single `Resolve(path string) (string, error)` method. All vault backends implement this interface.

**`internal/vault/keepassxc/keepassxc.go`** -- KeePassXC implementation of `Provider`. Parses the `kp://` path into an entry path and credential name, then shells out to `keepassxc-cli show` to retrieve the value.

**`internal/vault/onepassword/onepassword.go`** -- 1Password implementation of `Provider`. Validates the `op://` path into vault/item/field segments, then shells out to `op read 'op://<path>'`. Auth is handled by the `op` CLI itself (desktop app integration, an active `op signin` session, or `OP_SERVICE_ACCOUNT_TOKEN`), which the subprocess inherits from the parent environment.

**`internal/runner/exec.go`** -- Exec's the child command with the fully resolved environment. Wires stdin/stdout/stderr through to the child process.

---

## Releasing

Releases are automated via [GoReleaser](https://goreleaser.com/) and GitHub Actions.

```sh
make release VERSION=0.2.0
```

This tags the commit and pushes the tag. GitHub Actions builds binaries for macOS and Linux, creates a GitHub Release, and updates the Homebrew tap formula.

All releases are listed on the [releases page](https://github.com/the-known-unknown/key-env/releases).
