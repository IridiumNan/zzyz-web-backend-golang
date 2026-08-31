# guide-zzyz-web backend golang

This repo is the golang version for [guide-zzyz-web](https://github.com/Guidezzyz/guide-zzyz-web) project which aims at offering useful informations for zzyz school.

You can also check the [system diagram](./docs/architecture_diagram.md) for more informations.

Now is phrase one, check the [target](./docs/phrase_one.md)

Member CLI is written by python. Visit [python-cli](./docs/python-cli.md) for details.

---

## Quick Start

- run golang backend

```bash
# build locally (You should have golang env on you computer)
go build -o ./bin/zzyz ./cmd/web

# start the app
./bin/zzyz
```

- run python member cli

```bash
python memberPyCLI/main.py
```

---

## TEST

```bash
cd test && go test
```

---
