package zulip_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psyb0t/go-zulip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testEmail       = "bot@example.invalid"
	testAPIKey      = "EXAMPLE-DO-NOT-USE"
	testAPIPath     = "/api/v1"
	jsonContentType = "application/json"
	contentType     = "Content-Type"
)

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *zulip.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := zulip.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.SetBasicAuth(testEmail, testAPIKey)

			return nil
		},
	)

	client, err := zulip.NewClientWithResponses(srv.URL+testAPIPath, auth)
	require.NoError(t, err)

	return client
}

func assertBasicAuth(t *testing.T, r *http.Request) {
	t.Helper()

	user, pass, ok := r.BasicAuth()
	assert.True(t, ok)
	assert.Equal(t, testEmail, user)
	assert.Equal(t, testAPIKey, pass)
}

func TestGetOwnUser(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, testAPIPath+"/users/me", r.URL.Path)
		assertBasicAuth(t, r)

		w.Header().Set(contentType, jsonContentType)
		_, err := w.Write([]byte(`{
			"result":"success", "msg":"", "user_id":42,
			"email":"bot@example.invalid", "full_name":"Example bot",
			"is_bot":true, "is_active":true
		}`))
		assert.NoError(t, err)
	})

	res, err := client.GetOwnUserWithResponse(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)

	assert.Equal(t, 42, *res.JSON200.UserId)
	assert.Equal(t, testEmail, *res.JSON200.Email)
	assert.Equal(t, "Example bot", *res.JSON200.FullName)
	assert.True(t, *res.JSON200.IsBot)
	assert.True(t, *res.JSON200.IsActive)
}

func TestUnauthorizedResponse(t *testing.T) {
	t.Parallel()

	const response = `{"result":"error","msg":"Invalid API key"}`

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r)
		w.Header().Set(contentType, jsonContentType)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(response))
		assert.NoError(t, err)
	})

	res, err := client.GetOwnUserWithResponse(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode())
	assert.JSONEq(t, response, string(res.Body))
	assert.Nil(t, res.JSON200)
}

func TestSendMessage(t *testing.T) {
	t.Parallel()

	const (
		messageContent = "Hello + &= \u2603"
		messageID      = 73
		messageURL     = "https://chat.example.invalid/#narrow/dm/42/near/73"
	)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, testAPIPath+"/messages", r.URL.Path)
		assert.Equal(t,
			"application/x-www-form-urlencoded",
			r.Header.Get(contentType),
		)
		assertBasicAuth(t, r)
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "direct", r.PostForm.Get("type"))
		assert.JSONEq(t, "[42]", r.PostForm.Get("to"))
		assert.Equal(t, messageContent, r.PostForm.Get("content"))
		assert.Equal(t, "true", r.PostForm.Get("read_by_sender"))

		w.Header().Set(contentType, jsonContentType)
		err := json.NewEncoder(w).Encode(map[string]any{
			"result":       "success",
			"msg":          "",
			"id":           messageID,
			"message_url":  messageURL,
			"message_link": messageURL,
		})
		assert.NoError(t, err)
	})

	var to zulip.SendMessageFormdataBody_To

	require.NoError(t, to.FromSendMessageFormdataBodyTo3([]int{42}))

	readBySender := true

	res, err := client.SendMessageWithFormdataBodyWithResponse(
		context.Background(),
		zulip.SendMessageFormdataRequestBody{
			Type:         zulip.SendMessageFormdataBodyTypeDirect,
			To:           to,
			Content:      messageContent,
			ReadBySender: &readBySender,
		},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	assert.Equal(t, messageID, res.JSON200.Id)
	assert.Equal(t, messageURL, res.JSON200.MessageUrl)
	assert.Equal(t, messageURL, res.JSON200.MessageLink)
}

func TestSendChannelMessage(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r)
		assert.NoError(t, r.ParseForm())
		assert.JSONEq(t, `"Example + &="`, r.PostForm.Get("to"))
		assert.Equal(t, "channel", r.PostForm.Get("type"))
		assert.Equal(t, "Hello + &=", r.PostForm.Get("content"))
		assert.Equal(t, "Example topic", r.PostForm.Get("topic"))

		w.Header().Set(contentType, jsonContentType)
		_, err := w.Write([]byte(`{
			"result":"success", "msg":"", "id":74,
			"message_url":"", "message_link":""
		}`))
		assert.NoError(t, err)
	})

	var to zulip.SendMessageFormdataBody_To

	require.NoError(t, to.FromSendMessageFormdataBodyTo0("Example + &="))

	topic := "Example topic"

	res, err := client.SendMessageWithFormdataBodyWithResponse(
		context.Background(),
		zulip.SendMessageFormdataRequestBody{
			Type:    zulip.SendMessageFormdataBodyTypeChannel,
			To:      to,
			Content: "Hello + &=",
			Topic:   &topic,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	assert.Equal(t, 74, res.JSON200.Id)
}

func TestRegisterQueueMobileTimeout(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, testAPIPath+"/register", r.URL.Path)
		assertBasicAuth(t, r)
		assert.NoError(t, r.ParseForm())
		assert.JSONEq(t, `"mobile"`, r.PostForm.Get("idle_queue_timeout"))

		w.Header().Set(contentType, jsonContentType)
		_, err := w.Write([]byte(`{
			"result":"success", "msg":"", "queue_id":"example-queue"
		}`))
		assert.NoError(t, err)
	})

	var timeout zulip.RegisterQueueFormdataBody_IdleQueueTimeout

	require.NoError(t, timeout.FromRegisterQueueFormdataBodyIdleQueueTimeout1(
		zulip.RegisterQueueFormdataBodyIdleQueueTimeout1Mobile,
	))

	res, err := client.RegisterQueueWithFormdataBodyWithResponse(
		context.Background(),
		zulip.RegisterQueueFormdataRequestBody{IdleQueueTimeout: &timeout},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	require.NotNil(t, res.JSON200.QueueId)
	assert.Equal(t, "example-queue", *res.JSON200.QueueId)
}

func TestGetMessagesWithNarrow(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, testAPIPath+"/messages", r.URL.Path)
		assertBasicAuth(t, r)
		query := r.URL.Query()
		assert.Len(t, query["narrow"], 1)
		assert.JSONEq(t,
			`[{"operator":"channel","operand":"Example + &= "}]`,
			query.Get("narrow"),
		)
		assert.Equal(t, "newest", query.Get("anchor"))
		assert.Equal(t, "1", query.Get("num_before"))
		assert.Equal(t, "0", query.Get("num_after"))

		w.Header().Set(contentType, jsonContentType)
		_, err := w.Write([]byte(`{
			"result":"success", "msg":"", "found_newest":true,
			"messages":[{"id":73,"content":"Hello","type":"stream"}]
		}`))
		assert.NoError(t, err)
	})

	var operand zulip.GetMessagesParams_Narrow_0_Operand

	require.NoError(t, operand.FromGetMessagesParamsNarrow0Operand0(
		"Example + &= ",
	))

	var item zulip.GetMessagesParams_Narrow_Item

	require.NoError(t, item.FromGetMessagesParamsNarrow0(
		zulip.GetMessagesParamsNarrow0{
			Operator: "channel",
			Operand:  operand,
		},
	))
	narrow := []zulip.GetMessagesParams_Narrow_Item{item}
	anchor := "newest"
	numBefore, numAfter := 1, 0

	res, err := client.GetMessagesWithResponse(
		context.Background(),
		&zulip.GetMessagesParams{
			Narrow:    &narrow,
			Anchor:    &anchor,
			NumBefore: &numBefore,
			NumAfter:  &numAfter,
		},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.Len(t, res.JSON200.Messages, 1)
	assert.Equal(t, float64(73), res.JSON200.Messages[0].Id)
	assert.Equal(t, "Hello", res.JSON200.Messages[0].Content)
	assert.True(t, *res.JSON200.FoundNewest)
}

func TestSubscribe(t *testing.T) {
	t.Parallel()

	const subscriptions = `[{"name":"Example + &=","description":"Test"}]`

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, testAPIPath+"/users/me/subscriptions", r.URL.Path)
		assertBasicAuth(t, r)
		assert.NoError(t, r.ParseForm())
		assert.Len(t, r.PostForm["subscriptions"], 1)
		assert.JSONEq(t, subscriptions, r.PostForm.Get("subscriptions"))

		w.Header().Set(contentType, jsonContentType)
		_, err := w.Write([]byte(`{
			"result":"success", "msg":"",
			"subscribed":{"42":["Example + &="]}, "already_subscribed":{}
		}`))
		assert.NoError(t, err)
	})
	body := zulip.SubscribeFormdataRequestBody{}
	require.NoError(t, json.Unmarshal(
		[]byte(`{"subscriptions":`+subscriptions+`}`),
		&body,
	))

	res, err := client.SubscribeWithFormdataBodyWithResponse(
		context.Background(), body,
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.NotNil(t, res.JSON200.Subscribed)
	assert.Equal(t, []string{"Example + &="}, (*res.JSON200.Subscribed)["42"])
}

func TestMalformedResponse(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(contentType, jsonContentType)
		_, err := w.Write([]byte(`{"result":`))
		assert.NoError(t, err)
	})

	res, err := client.GetOwnUserWithResponse(context.Background())

	var syntaxError *json.SyntaxError

	require.ErrorAs(t, err, &syntaxError)
	assert.Nil(t, res)
}

func TestCancelledRequest(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(_ http.ResponseWriter, _ *http.Request) {
		assert.Fail(t, "cancelled request reached server")
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	res, err := client.GetOwnUserWithResponse(ctx)
	require.ErrorIs(t, err, context.Canceled)
	assert.Nil(t, res)
}
