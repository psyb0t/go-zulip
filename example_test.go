package zulip_test

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-zulip"
)

// This README example needs a Zulip server and real credentials.
//
//nolint:testableexamples // Needs real credentials, compile check only.
func Example() {
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
