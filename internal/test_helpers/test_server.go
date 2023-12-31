package testhelpers

import (
	"altt/internal/routes"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
)

type TestServer struct {
	appPort int
	client  http.Client
}

func NewTestServer(t *testing.T, container *TestContainer) *TestServer {
	appPort, err := freeport.GetFreePort()
	require.NoError(t, err, "failed to get free port for app")
	srv := &TestServer{
		appPort: appPort,
		client:  *http.DefaultClient,
	}

	appHTTPServer := routes.InitAppRouter(
		container.Log,
		container.ServiceBalancer,
		fmt.Sprintf(":%d", srv.appPort),
		container.Conf.DisableMetrics,
	)
	t.Cleanup(func() {
		require.NoError(t, appHTTPServer.Stop())
	})
	go func() {
		require.NoError(t, appHTTPServer.Run())
	}()
	return srv
}

func (ts *TestServer) Request(t *testing.T, method, path string, body interface{}, headers map[string]string) *TestResponse {
	t.Helper()

	var b []byte
	var err error
	if body != nil {
		if headers == nil {
			headers = make(map[string]string)
		}
		headers["Content-Type"] = "application/json"
		b, err = json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
	}

	u := fmt.Sprintf("http://localhost:%d%s", ts.appPort, path)
	req, err := http.NewRequest(method, u, bytes.NewBuffer(b))
	require.NoError(t, err, "failed to construct new request for url %s: %s", u, err)
	if err != nil {
		t.Fatal(err)
	}

	if len(headers) > 0 {
		for headerKey, headerVal := range headers {
			req.Header.Add(headerKey, headerVal)
		}
	}

	res, err := ts.client.Do(req) //nolint:bodyclose
	require.NoError(t, err, "failed to make request to %s: %s", u, err)
	t.Cleanup(func() {
		require.NoError(t, res.Body.Close())
	})
	return &TestResponse{Res: res}
}

func (ts *TestServer) Get(t *testing.T, path string) *TestResponse {
	t.Helper()
	return ts.Request(t, http.MethodGet, path, nil, nil)
}
