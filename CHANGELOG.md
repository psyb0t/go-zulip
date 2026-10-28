# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.2.0, 2026-10-28

- Breaking: the numbered union variants of the `GetEvents` response body (`GetEvents200JSONResponseBodyEvents11` and the other `GetEvents200JSONResponseBodyEvents<N>` types) were regenerated from the updated spec. Many were removed or reshaped, so code that names them must change. Examples: `GetEvents200JSONResponseBodyEvents12.Op` and `.Person`, `GetEvents200JSONResponseBodyEvents14.Presence`, `.Presences`, `.UserId` and `.Email`, and `GetEvents200JSONResponseBodyEvents15.Streams` no longer exist, along with their related `Op`, `Type` and `PresenceStatus` constants.
- Changed field types: `GetEvents200JSONResponseBodyEvents13.Person` is now `*User`, `GetEvents200JSONResponseBodyEvents37.Id` is now `*EventIdSchema`, and `GetEvents200JSONResponseBodyEvents16.Streams` and `GetEvents200JSONResponseBodyEvents22.Attachment` have new types.
- `RegisterQueueResponse.JSON200` and `RegisterQueueResponse.GetJSON200` return a different anonymous struct type, so code that assigns or declares the old type must change.
- Added the Google Meet video call endpoint: `Client.CreateGoogleMeetVideoCall`, `ClientWithResponses.CreateGoogleMeetVideoCallWithResponse`, `NewCreateGoogleMeetVideoCallRequest`, `ParseCreateGoogleMeetVideoCallResponse` and `CreateGoogleMeetVideoCallResponse`. The matching methods were added to `ClientInterface` and `ClientWithResponsesInterface`, so custom implementations of those interfaces must add them.
- Added the event type constant `GetEvents200JSONResponseBodyEvents11TypeHasGoogleMeetToken`.

## v0.1.3, 2026-10-28

- Updated Semgrep in the development image to 1.182.1 and refreshed its locked dependencies. certifi 2026.10.23 drops the emSign root CAs whose intermediate keys were compromised (GHSA-6q48-6c7w-9gw2).
- Dependabot no longer opens pull requests for `mcp` or the OpenTelemetry packages in `requirements-semgrep.txt`. Semgrep pins their versions itself, so bumping one of them on its own cannot install.
- Updated kin-openapi, which the code generator uses, from 0.142.0 to 0.149.0. The newer version fixes several panics on malformed specs reported in GitHub security advisories.
- The security scan merges its two gitleaks runs into one SARIF run. GitHub code scanning rejects a report with two runs from the same tool, so no scan results had been uploaded.

## v0.1.2, 2026-10-27

- The development image is based on Debian (`golang:1.26.10-trixie`) instead of Alpine. Alpine keeps only the latest build of each package, so the exact pins broke the image build whenever a package was updated. System packages now come from the snapshot.debian.org date the base image was built from, which keeps them fixed for a given base image digest. jq 1.8.2, shellcheck 0.11.0 and shfmt 3.13.1 keep their versions, and Semgrep runs in its own virtual environment.
- Raised the minimum Go version in `go.mod` from 1.26.9 to 1.26.10, which fixes the `crypto/x509` vulnerability GO-2026-6809 found by govulncheck.
- The issue-pull workflow runs at hours derived from the repository name instead of the same four hours as every other repository.

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
