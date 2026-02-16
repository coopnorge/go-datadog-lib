package http_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	datadogMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/http"
	"github.com/stretchr/testify/require"
)

func TestResourceNamers(t *testing.T) {
	testCases := []struct {
		name                         string
		rn                           datadogMiddleware.ResourceNamer
		expectedFullPath             string
		expectedFullWithQuery        string
		expectedFullWithQueryAndUser string
	}{
		{
			name:                         "StaticResourceNamer",
			rn:                           datadogMiddleware.StaticResourceNamer("foobar"),
			expectedFullPath:             "foobar",
			expectedFullWithQuery:        "foobar",
			expectedFullWithQueryAndUser: "foobar",
		},
		{
			name:                         "FullURLWithParamsResourceNamer",
			rn:                           datadogMiddleware.FullURLWithParamsResourceNamer(),
			expectedFullPath:             "GET https://www.coop.no/api/some-service/some-endpoint",
			expectedFullWithQuery:        "GET https://www.coop.no/api/some-service/some-endpoint?foo=bar",
			expectedFullWithQueryAndUser: "GET https://bax:xxxxx@www.coop.no/api/some-service/some-endpoint?foo=bar",
		},
		{
			name:                         "FullURLResourceNamer",
			rn:                           datadogMiddleware.FullURLResourceNamer(),
			expectedFullPath:             "GET https://www.coop.no/api/some-service/some-endpoint",
			expectedFullWithQuery:        "GET https://www.coop.no/api/some-service/some-endpoint",
			expectedFullWithQueryAndUser: "GET https://www.coop.no/api/some-service/some-endpoint",
		},
		{
			name:                         "HostResourceNamer",
			rn:                           datadogMiddleware.HostResourceNamer(),
			expectedFullPath:             "www.coop.no",
			expectedFullWithQuery:        "www.coop.no",
			expectedFullWithQueryAndUser: "www.coop.no",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest("GET", "https://www.coop.no/api/some-service/some-endpoint", nil)
			require.NoError(t, err)
			require.Equal(t, tc.expectedFullPath, tc.rn(req))

			q := req.URL.Query()
			q.Add("foo", "bar")
			req.URL.RawQuery = q.Encode()
			require.Equal(t, tc.expectedFullWithQuery, tc.rn(req))

			req.SetBasicAuth("bax", "baz")
			req.URL.User = url.UserPassword("bax", "baz")
			require.Equal(t, tc.expectedFullWithQueryAndUser, tc.rn(req))
		})
	}
}

func TestCustomResourceNamer(t *testing.T) {
	// Example ResourceNamer that can be used to e.g. not log customer-id.
	rn := func(req *http.Request) string {
		u := req.URL
		path := u.Path
		if strings.HasPrefix(path, "/api/some-service/customers/") {
			path = "/api/some-service/customers/:customerid"
		}
		return req.Method + " " + u.Scheme + "://" + u.Host + path
	}

	req, err := http.NewRequest("GET", "https://www.coop.no/api/some-service/some-endpoint", nil)
	require.NoError(t, err)
	require.Equal(t, "GET https://www.coop.no/api/some-service/some-endpoint", rn(req))

	q := req.URL.Query()
	q.Add("foo", "bar")
	req.URL.RawQuery = q.Encode()
	require.Equal(t, "GET https://www.coop.no/api/some-service/some-endpoint", rn(req))

	req, err = http.NewRequest("GET", "https://www.coop.no/api/some-service/customers/1234", nil)
	require.NoError(t, err)
	require.Equal(t, "GET https://www.coop.no/api/some-service/customers/:customerid", rn(req))
}
