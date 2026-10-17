// Package zulip is a Go client for the Zulip REST API,
// generated with oapi-codegen from the OpenAPI document Zulip publishes.
package zulip

//go:generate bash scripts/fetch-spec.sh
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.yaml
