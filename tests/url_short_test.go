package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"path"
	"testing"
	"url_short/internal/http-server/handlers/url/save"
	"url_short/internal/lib/api"
	"url_short/internal/lib/random"
)

const (
	host = "localhost:8080"
)

func TestURLShort_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(10),
		}).
		WithBasicAuth("denis", "chort123").
		Expect().
		Status(200).
		JSON().
		Object().
		ContainsKey("alias")
}

func TestURLShortener_SaveRedirectRemove(t *testing.T) {
	testCass := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "field URL is not a valid URL",
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
	}

	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("denis", "chort123").
				Expect().Status(http.StatusOK).
				JSON().Object()
			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			testRedirect(t, alias, tc.url)

			reqDel := e.DELETE("/"+path.Join("url", alias)).
				WithBasicAuth("denis", "chort123").
				Expect().Status(http.StatusOK).
				JSON().Object()

			reqDel.Value("status").String().IsEqual("OK")

			testRedirectNotFound(t, alias)
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}
	redirectToURL, err := api.GetRedirect(u.String())

	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectToURL)
}

func testRedirectNotFound(t *testing.T, alias string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}
	_, err := api.GetRedirect(u.String())
	require.ErrorIs(t, err, api.ErrInvalidStatusCode)
}
