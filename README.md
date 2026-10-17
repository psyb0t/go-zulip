# go-zulip

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-zulip.svg)](https://pkg.go.dev/github.com/psyb0t/go-zulip)
[![CI](https://github.com/psyb0t/go-zulip/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-zulip/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-zulip/badges/version.svg)](https://github.com/psyb0t/go-zulip/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-zulip/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-zulip/badges/importers.svg)](https://github.com/psyb0t/go-zulip/blob/badges/importers.md)

Unofficial Go client for the [Zulip REST API](https://zulip.com/api/), generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the OpenAPI document Zulip publishes. I am not affiliated with Zulip.

The package is `zulip`, in the repository root. Everything in `client.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-zulip
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`. Responses decode into `JSON200` and the other status fields. The full API, with the spec's descriptions as doc comments, is in the [Go reference](https://pkg.go.dev/github.com/psyb0t/go-zulip) and in `go doc github.com/psyb0t/go-zulip`.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-zulip"
)

func main() {
	auth := zulip.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.SetBasicAuth(
				os.Getenv("ZULIP_EMAIL"),
				os.Getenv("ZULIP_API_KEY"),
			)

			return nil
		},
	)

	client, err := zulip.NewClientWithResponses(
		os.Getenv("ZULIP_SERVER"),
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.GetOwnUserWithResponse(context.Background())
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

Set `ZULIP_SERVER` to your server's API base URL, for example `https://your-server.example/api/v1`. The same code is in `example_test.go`, so `make test` checks that it compiles.

Method names come from the spec's operation IDs: `GetOwnUserWithResponse` reads the current user, `GetMessagesWithResponse` fetches message history and `SendMessageWithFormdataBodyWithResponse` sends a message. The client also covers users, channels, subscriptions, events, organization settings and administration.

## Auth

Zulip's `BasicAuth` security scheme uses HTTP Basic auth: the account or bot email address is the username and its API key is the password. Pass `WithRequestEditorFn` as shown above. The client does not attach credentials automatically.

Use the API key for the intended Zulip organization, and send it only to that organization's trusted HTTPS server. A Zulip Cloud base URL looks like `https://your-organization.zulipchat.com/api/v1`; a self-hosted server uses its own origin followed by `/api/v1`.

Check `StatusCode()` and the decoded status fields before using a response. A non-2xx API response is not necessarily a Go error; the error return covers request, transport and decoding failures.

## Regenerating

The spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` downloads `https://raw.githubusercontent.com/zulip/zulip/main/zerver/openapi/zulip.yaml` to `api/openapi.yaml`.
- oapi-codegen runs with `config.yaml` across the whole spec. The config maps `number` to `float64` instead of the default `float32`.
- `scripts/client.tmpl`, selected through `user-templates` in the config, encodes composite form values as JSON and honors form fields marked `application/json`, including quoted strings. This fixes dropped union fields and indexed arrays in the default form encoder. The source spec is unchanged.

`.github/workflows/regen.yml` does the same every day in CI. When the new spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude edits the code and the checks run again, up to three times, and a review step reads the result. A passing change is committed with a CHANGELOG entry and tagged as the next release.

## Agent integrations

The [skill](.agents/skills/go-zulip) teaches a coding agent how to use go-zulip when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-zulip@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-zulip@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-zulip
```

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE) for repository code. Generated material from the Zulip specification and the adapted oapi-codegen template retain [Apache 2.0](LICENSE-APACHE-2.0) notices; see [NOTICE](NOTICE). The specification itself is not redistributed here.

The spec warns that server-provided GIPHY, Tenor and KLIPY API keys are for Zulip clients only. Their provider terms apply to those integrations.
