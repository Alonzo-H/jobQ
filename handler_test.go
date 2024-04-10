package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func newTestClient(t *testing.T) (*httpexpect.Expect, func()) {
    jh := NewJobsHandler()

	server := httptest.NewServer(jh.mux)

	client := httpexpect.WithConfig(httpexpect.Config{
		TestName: t.Name(),
		BaseURL:  server.URL,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

    return client, server.Close
}

func TestJobInfoNotFound(t *testing.T) {
    client, closer := newTestClient(t)
    defer closer()

	client.GET("/jobs/1").Expect().Status(http.StatusNotFound)
}

func TestEnqueueDequeueConclude(t *testing.T) {
    client, closer := newTestClient(t)
    defer closer()

    enqResp := client.POST("/jobs/enqueue").
        WithJSON(map[string]string{
            "type": "TIME_CRITICAL",
        }).
        Expect().
        Status(http.StatusCreated).
        JSON().Object()

    enqId := enqResp.Value("id").Number().Raw()
    enqResp.Value("type").String().IsEqual("TIME_CRITICAL")
    enqResp.Value("status").String().IsEqual("QUEUED")

    getResp := client.GET("/jobs/{jobId}", fmt.Sprintf("%.0f", enqId)).
        Expect().
        Status(http.StatusOK).
        JSON().Object()
    getResp.Value("id").Number().IsEqual(enqId)
    getResp.Value("type").String().IsEqual("TIME_CRITICAL")
    getResp.Value("status").String().IsEqual("QUEUED")

    deqResp := client.POST("/jobs/dequeue").
        Expect().
        Status(http.StatusOK).
        JSON().Object()

    deqResp.Value("id").Number().IsEqual(enqId)
    deqResp.Value("type").String().IsEqual("TIME_CRITICAL")
    deqResp.Value("status").String().IsEqual("IN_PROGRESS")

    getResp = client.GET("/jobs/{jobId}", fmt.Sprintf("%.0f", enqId)).
        Expect().
        Status(http.StatusOK).
        JSON().Object()
    getResp.Value("id").Number().IsEqual(enqId)
    getResp.Value("type").String().IsEqual("TIME_CRITICAL")
    getResp.Value("status").String().IsEqual("IN_PROGRESS")

    client.PUT("/jobs/{jobId}/conclude", fmt.Sprintf("%.0f", enqId)).
        Expect().
        Status(http.StatusOK)

    getResp = client.GET("/jobs/{jobId}", fmt.Sprintf("%.0f", enqId)).
        Expect().
        Status(http.StatusOK).
        JSON().Object()
    getResp.Value("id").Number().IsEqual(enqId)
    getResp.Value("type").String().IsEqual("TIME_CRITICAL")
    getResp.Value("status").String().IsEqual("CONCLUDED")
}

func Test3EnqueueConclude2Dequeue(t *testing.T) {
    client, closer := newTestClient(t)
    defer closer()

    enqResp := client.POST("/jobs/enqueue").
        WithJSON(map[string]string{
            "type": "TIME_CRITICAL",
        }).
        Expect().
        Status(http.StatusCreated).
        JSON().Object()
    enqId1 := enqResp.Value("id").Number().Raw()

    enqResp = client.POST("/jobs/enqueue").
        WithJSON(map[string]string{
            "type": "NOT_TIME_CRITICAL",
        }).
        Expect().
        Status(http.StatusCreated).
        JSON().Object()
    enqId2 := enqResp.Value("id").Number().Raw()

    enqResp = client.POST("/jobs/enqueue").
        WithJSON(map[string]string{
            "type": "NOT_TIME_CRITICAL",
        }).
        Expect().
        Status(http.StatusCreated).
        JSON().Object()
    enqId3 := enqResp.Value("id").Number().Raw()

    client.PUT("/jobs/{jobId}/conclude", fmt.Sprintf("%.0f", enqId2)).
        Expect().
        Status(http.StatusOK)

    deqResp := client.POST("/jobs/dequeue").
        Expect().
        Status(http.StatusOK).
        JSON().Object()

    deqResp.Value("id").Number().IsEqual(enqId1)
    deqResp.Value("type").String().IsEqual("TIME_CRITICAL")
    deqResp.Value("status").String().IsEqual("IN_PROGRESS")

    deqResp = client.POST("/jobs/dequeue").
        Expect().
        Status(http.StatusOK).
        JSON().Object()

    deqResp.Value("id").Number().IsEqual(enqId3)
    deqResp.Value("type").String().IsEqual("NOT_TIME_CRITICAL")
    deqResp.Value("status").String().IsEqual("IN_PROGRESS")
}

func TestEnqueueEnqueueCancelDequeue(t *testing.T) {
    client, closer := newTestClient(t)
    defer closer()

    enqResp := client.POST("/jobs/enqueue").
        WithJSON(map[string]string{
            "type": "TIME_CRITICAL",
        }).
        Expect().
        Status(http.StatusCreated).
        JSON().Object()
    enqId1 := enqResp.Value("id").Number().Raw()

    enqResp = client.POST("/jobs/enqueue").
        WithJSON(map[string]string{
            "type": "TIME_CRITICAL",
        }).
        Expect().
        Status(http.StatusCreated).
        JSON().Object()
    enqId2 := enqResp.Value("id").Number().Raw()

    client.PUT("/jobs/{jobId}/cancel", fmt.Sprintf("%.0f", enqId1)).
        Expect().
        Status(http.StatusOK)

    deqResp := client.POST("/jobs/dequeue").
        Expect().
        Status(http.StatusOK).
        JSON().Object()

    deqResp.Value("id").Number().IsEqual(enqId2)
    deqResp.Value("type").String().IsEqual("TIME_CRITICAL")
    deqResp.Value("status").String().IsEqual("IN_PROGRESS")
}
