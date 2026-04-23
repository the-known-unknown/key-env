# key-env

[![CI](https://github.com/the-known-unknown/key-env/actions/workflows/ci.yml/badge.svg)](https://github.com/the-known-unknown/key-env/actions/workflows/ci.yml)
[![Release](https://github.com/the-known-unknown/key-env/actions/workflows/release.yml/badge.svg)](https://github.com/the-known-unknown/key-env/releases)

## The problem

`.env` files are everywhere in local development. They store API keys, database credentials, OAuth secrets, and other sensitive values that your application needs to run. Typically these secrets sit in plaintext, committed to repos or shared over Slack.

Modern software makes this worse. Applications depend on dozens of third-party services, each with their own credentials. Agentic coding tools — Copilot, Cursor, Claude — now read and write code on your behalf, which means they also read your `.env` files. Every new dependency and every AI agent that touches your project increases the surface area of exposure.

The result: plaintext secrets scattered across machines, repos, and chat histories, accessible to tools and processes you don't fully control.

## What key-env does

`key-env` keeps your secrets encrypted in a local password manager and out of your `.env` files entirely.

Instead of storing actual secret values, your `.env` file contains **references** that point to entries in your password manager:

```env
CLIENT_SECRET="kp://my-api/Password"
CLIENT_NAME="kp://my-api/Username"
```

When you run your application through `key-env`, it:

1. Reads your `.env` file and identifies secret references (e.g. `kp://...`)
2. Resolves each reference by querying your local password manager
3. Injects the decrypted values into the environment
4. Launches your application with the fully hydrated environment

Your secrets never leave the password manager until the moment they're needed, and they never touch disk in plaintext.

### Supported providers


| Prefix  | Provider      | Status      |
| ------- | ------------- | ----------- |
| `kp://` | KeePassXC     | Implemented |
| `op://` | 1Password CLI | Planned     |


## Install

**Via Homebrew (macOS):**

```sh
brew install the-known-unknown/tap/key-env
```

**From source:**

```sh
git clone https://github.com/the-known-unknown/key-env.git
cd key-env
make build
```

**From GitHub Releases:**

Download the latest binary for your platform from the [releases page](https://github.com/the-known-unknown/key-env/releases).

## Dependencies

- **Go 1.22+** — required to build from source
- **KeePassXC 2.7+** — the `keepassxc-cli` binary must be available on your `PATH` (install via `brew install keepassxc` on macOS). Tested with `2.7.11`.

## Usage

### What an `.env` file looks like

Plain values work as-is:

```env
DB_NAME="userdb"
APP_PORT="3000"
```

Secret references use the `kp://` prefix to point into your KeePassXC database:

```env
DB_PASSWORD="kp://Services/Databases/Main/Password"
API_KEY="kp://Services/Stripe/SecretKey"
```

The format is `kp://<entry_path>/<credential>`, where:

- `entry_path` is the path to the entry in your KeePassXC database (e.g. `Services/Databases/Main`)
- `credential` is the field to retrieve (e.g. `Password`, `Username`)

You can mix plain values and secret references in the same file.

### How to use it

```sh
./key-env run \
  --env <path-to-env-file> \
  --secrets <path-to-kdbx-file> \
  --password '<vault-password>' \
  -- <your command>
```

For example, to run `npm test` with secrets loaded:

```sh
./key-env run \
  --env .env.dev \
  --secrets ./secrets.kdbx \
  --password 'my-vault-password' \
  -- npm test
```


| Flag         | Required | Description                               |
| ------------ | -------- | ----------------------------------------- |
| `--env`      | Yes      | Path to your `.env` file                  |
| `--secrets`  | Yes      | Path to your `.kdbx` KeePassXC database   |
| `--password` | Yes      | Password to unlock the KeePassXC database |
| `--`         | Yes      | Separator before the child command        |


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

### Try it out

A sample KeePassXC database and env file are included in the `test/` directory:

```sh
$ make clean
$ make build
$ ./key-env run \
    --env test/.env.sample \
    --secrets test/keepass-sample-db.kdbx \
    --password '4jFU%i*+Q2qdpFgoHJGK' \
    -- sh -c 'echo $TEST_CLIENT_SECRET'

>> Output:
Jcg5TfdI9X0zHaU03Qx9bGb0rphYh0xIebtpFPTcRT
```

### Running sample scripts

The `test/` directory also includes sample Python and Node.js scripts that print the resolved secrets and show a progress bar:

**Python:**

```sh
$ ./key-env run \
    --env test/.env.sample \
    --secrets test/keepass-sample-db.kdbx \
    --password '4jFU%i*+Q2qdpFgoHJGK' \
    -- python3 test/sample.py
```

```
Env vars:
TEST_CLIENT_SECRET: Jcg5TfdI9X0zHaU03Qx9bGb0rphYh0xIebtpFPTcRT
TEST_CLIENT_NAME: test
--------------------------------

Doing some work...
[██████████] 100%
✔ Done!
```

**Node.js:**

```sh
$ ./key-env run \
    --env test/.env.sample \
    --secrets test/keepass-sample-db.kdbx \
    --password '4jFU%i*+Q2qdpFgoHJGK' \
    -- node test/sample.js
```

```
Env vars:
TEST_CLIENT_SECRET: Jcg5TfdI9X0zHaU03Qx9bGb0rphYh0xIebtpFPTcRT
TEST_CLIENT_NAME: test
--------------------------------

Doing some work...
[██████████] 100%
✔ Done!
```

## Releasing

Releases are automated via [GoReleaser](https://goreleaser.com/) and GitHub Actions. To cut a new release:

```sh
make release VERSION=0.2.0
```

This tags the commit and pushes the tag. GitHub Actions then:

1. Runs tests
2. Builds binaries for macOS (amd64 + arm64) and Linux (amd64 + arm64)
3. Creates a GitHub Release with downloadable archives
4. Updates the Homebrew tap formula

All releases are listed on the [releases page](https://github.com/the-known-unknown/key-env/releases).

Commit messages are used to generate release notes. Use conventional prefixes for cleaner changelogs:

| Prefix   | Category      |
|----------|---------------|
| `feat:`  | Features      |
| `fix:`   | Bug fixes     |
| `docs:`  | Documentation |

## Security notes

- The `--password` flag can expose your vault password in shell history and process listings. In production workflows, prefer passing it via stdin or a secure file.
- Your `.env` file still reveals metadata — variable names and vault paths — even though the actual secret values are encrypted. Be mindful when sharing or committing it.
- The decrypted secrets exist in the child process's environment for the duration of its execution. They are not written to disk.

