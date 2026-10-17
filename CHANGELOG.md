# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.0, 2026-10-17

- Added `zulip`, a REST client generated with oapi-codegen from Zulip's OpenAPI document (Zulip REST API, spec version 1.0.0, OpenAPI 3.0.1). It covers all 119 paths and 166 operations: messages, users, channels, subscriptions, events, organization settings and administration.
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- A client template override encodes composite form values as JSON for Zulip, preserving union fields and honoring `application/json` field encoding. No source spec patch is required.
- `make generate` downloads the current spec from https://raw.githubusercontent.com/zulip/zulip/main/zerver/openapi/zulip.yaml and regenerates the client. The spec is not committed.
