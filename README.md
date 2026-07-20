# flogd

**flogd** (*follow log & do*) watches the output of a command (or a log stream) and
runs an action when a pattern shows up often enough within a time window.

A typical use: tail an auth log, and when the same IP fails to log in *N* times in
*M* seconds, run a command to block it.

## How it works

1. **Tail** — flogd runs a command (e.g. `tail -f /var/log/auth.log`) and reads its
   output line by line.
2. **Match** — each line is tested against a regex. If the regex has a capture
   group, matches are counted *per captured value* (e.g. per IP); otherwise they
   are counted globally.
3. **Do** — when a value reaches `count` matches within `interval` seconds, flogd
   runs the `do` command. If `do` contains `%s`, the captured value is substituted
   in.

## Install

### Prebuilt binary (recommended)

Download the binary for your platform from the latest release. Prebuilt
binaries are published for **Linux amd64/arm64** and **macOS arm64**:

```sh
# Pick the asset matching your OS/arch:
#   flogd-linux-amd64   flogd-linux-arm64   flogd-darwin-arm64
curl -fsSL -o flogd \
  https://github.com/aminehmida/flogd/releases/latest/download/flogd-linux-amd64

chmod +x flogd
sudo mv flogd /usr/local/bin/
```

Not sure which asset you need? Detect it automatically:

```sh
os=$(uname -s | tr '[:upper:]' '[:lower:]')       # linux | darwin
arch=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
curl -fsSL -o flogd \
  "https://github.com/aminehmida/flogd/releases/latest/download/flogd-${os}-${arch}"
chmod +x flogd && sudo mv flogd /usr/local/bin/
```

You can also download assets from the [Releases](../../releases) page directly.

### With Go

```sh
go install github.com/aminehmida/flogd@latest
```

### From source

```sh
git clone https://github.com/aminehmida/flogd.git
cd flogd
go build -o flogd .
```

## Usage

### One-off, from the command line

Block an IP after 5 failed SSH logins within 10 seconds:

```sh
flogd monitor \
  --regex 'Failed password for .* from ([0-9.]+)' \
  --count 5 \
  --interval 10 \
  --do 'echo blocking %s' \
  'tail -f /var/log/auth.log'
```

The positional argument is the command to monitor.

### From a config file

Run one or more monitors defined in YAML:

```sh
flogd monitor --config examples/flogd.yaml
```

### Save a monitor to a config file

`save` appends the given flags as a named entry (default file `./flogd.yaml`):

```sh
flogd save \
  --name ssh-bruteforce \
  --regex 'Failed password for .* from ([0-9.]+)' \
  --count 5 --interval 10 \
  --do 'echo blocking %s' \
  'tail -f /var/log/auth.log'
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--type` | `-t` | `process` | Stream type to monitor (`process` is currently supported) |
| `--regex` | `-r` | `*` | Regex to match against each line |
| `--count` | `-n` | `10` | Matches required before running the action |
| `--interval` | `-i` | `5` | Time window in seconds |
| `--do` | `-d` | — | Command to run on trigger; `%s` is replaced by the capture group |
| `--config` | `-c` | — | Config file to load |
| `--name` | `-m` | — | (`save`) config entry name |
| `--desc` | `-s` | — | (`save`) config entry description |

## Config format

Each entry in the YAML list is one monitor:

```yaml
- name: ssh-bruteforce         # identifier for this monitor
  type: process                # stream type
  description: ""              # optional
  regex: 'Failed password for .* from ([0-9.]+)'
  do: echo blocking %s          # %s = captured group
  count: 5                      # matches before triggering
  interval: 10                  # window in seconds
  command: tail -f /var/log/auth.log
```

See [`examples/flogd.yaml`](examples/flogd.yaml) for a runnable sample that reads
[`examples/logs.txt`](examples/logs.txt).

## Development

This repo uses [mise](https://mise.jdx.dev/) to pin the Go toolchain
(see `mise.toml`):

```sh
mise install      # install the pinned Go version
go test ./...     # run the test suite
go vet ./...      # static checks
go build -o flogd .
```

## License

[MIT](LICENSE) © Amine Hmida
