# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.1, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.
- The README names the spec's security scheme `BasicAuth` correctly.

## v0.1.0, 2026-10-17

- Added `zulip`, a REST client generated with oapi-codegen from Zulip's OpenAPI document (Zulip REST API, spec version 1.0.0, OpenAPI 3.0.1). It covers all 119 paths and 166 operations: messages, users, channels, subscriptions, events, organization settings and administration.
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- A client template override encodes composite form values as JSON for Zulip, preserving union fields and honoring `application/json` field encoding. No source spec patch is required.
- `make generate` downloads the current spec from https://raw.githubusercontent.com/zulip/zulip/main/zerver/openapi/zulip.yaml and regenerates the client. The spec is not committed.
