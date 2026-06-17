# intrastate

> One-line description of what intrastate does.

intrastate is a Go command-line tool. _(Replace this section with the
project's purpose once it is defined.)_

## Install

```sh
make install        # builds ./bin/intrastate and installs to ~/.local/bin
```

Or build locally:

```sh
make build          # ./bin/intrastate
```

## Usage

```sh
intrastate version              # build version, commit, date
intrastate version --as=json    # same, as a JSON envelope
```

Every command accepts the global `--as text|json` flag. Under `--as=json`
stdout carries a single terminal envelope discriminated by a `type`
field (`ok` | `failed`); under `--as=text` it is human-readable.

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for layout, build/test commands,
and the output/error contract every verb follows.

```sh
make check          # fmt-check + vet + lint + test (local mirror of CI)
make test           # race + coverage
```

## License

[MIT](LICENSE)
