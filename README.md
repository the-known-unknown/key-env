# key-env

`key-env` is a CLI utility that loads environment variables from a `.env` file, resolves secret references through local vault CLIs, and runs a child command with the hydrated environment.

## Current provider support

- `kp://...` via `keepassxc-cli` (implemented)
- `op://...` via 1Password CLI (stubbed, planned next)

## Install requirements

- Go 1.22+ for local builds
- `keepassxc-cli` available on `PATH`

## Usage

```sh
key-env run --env .env.dev --secrets ./secrets.kdbx --password 'my-password' -- <child command>
```

Example:

```sh
key-env run --env .env.dev --secrets ./secrets.kdbx --password 'my-password' -- npm test
```

### Mandatory inputs

- `--env` path to `.env` file
- `--secrets` path to `.kdbx` file
- `-- <child command>` command to execute in hydrated environment

### `.env` value formats

Plain value:

```env
DB_NAME="userdb"
```

KeePassXC reference:

```env
DB_PASSWORD="kp://Services/Databases/Main/Password"
```

`kp://<secret_path>/<credential>` is parsed as:

- `secret_path`: `Services/Databases/Main`
- `credential`: `Password`

For each `kp://` value, `key-env` runs:

```sh
keepassxc-cli show <secrets.kdbx> <secret_path> -q -s -a <credential>
```

and pipes password via stdin.

## Security notes

- Passing passwords through `--password` can expose them in shell history and process listings.
- Prefer secure input methods (stdin/password file/env var) in future revisions.
- Committed env files still expose metadata (variable names and vault paths).

## Development

Run tests:

```sh
go test ./...
```

Build:

```sh
go build ./cmd/keyenv
```
