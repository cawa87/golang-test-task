package server

import (
	"github.com/souz9/golang-test-task/scrapper/hchecker"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

type fakeChecker []hchecker.Status

func (f fakeChecker) Status() []hchecker.Status       { return f }
func (f fakeChecker) StatusOf(string) hchecker.Status { return f[0] }
func (f fakeChecker) FindMin() hchecker.Status        { return f[0] }
func (f fakeChecker) FindMax() hchecker.Status        { return f[len(f)-1] }

func TestServer(t *testing.T) {
	s := Server{}
	hcA := fakeChecker{
		{Site: "one.com", Live: true, Respond: 0.10},
		{Site: "two.org", Live: false, Respond: 1.10}}
	hcB := fakeChecker{}
	hcC := fakeChecker{{}}

	get := func(t *testing.T, url string) string {
		req := httptest.NewRequest("GET", url, nil)
		rec := httptest.NewRecorder()
		s.router().ServeHTTP(rec, req)
		res := rec.Result()

		body, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)

		return string(body)
	}

	t.Run("status", func(t *testing.T) {
		s.HChecker = hcA
		require.JSONEq(t, `[
			{"Site":"one.com", "Live":true,  "Respond":0.10},
			{"Site":"two.org", "Live":false, "Respond":1.10}
		]`, get(t, "/status"))

		t.Run("when no entries", func(t *testing.T) {
			s.HChecker = hcB
			require.Equal(t, `[]`, get(t, "/status"))
		})
	})

	t.Run("statusOf", func(t *testing.T) {
		s.HChecker = hcA
		require.JSONEq(t, `[
			{"Site":"one.com", "Live":true,  "Respond":0.10}
		]`, get(t, "/status?site=one.com"))

		t.Run("when no entries", func(t *testing.T) {
			s.HChecker = hcC
			require.Equal(t, `[]`, get(t, "/status?site=one.com"))
		})
	})

	t.Run("statusMin", func(t *testing.T) {
		s.HChecker = hcA
		require.JSONEq(t, `[
			{"Site":"one.com", "Live":true,  "Respond":0.10}
		]`, get(t, "/status/min"))

		t.Run("when no entries", func(t *testing.T) {
			s.HChecker = hcC
			require.Equal(t, `[]`, get(t, "/status/min"))
		})
	})

	t.Run("statusMax", func(t *testing.T) {
		s.HChecker = hcA
		require.JSONEq(t, `[
			{"Site":"two.org", "Live":false, "Respond":1.10}
		]`, get(t, "/status/max"))

		t.Run("when no entries", func(t *testing.T) {
			s.HChecker = hcC
			require.Equal(t, `[]`, get(t, "/status/max"))
		})
	})
}
