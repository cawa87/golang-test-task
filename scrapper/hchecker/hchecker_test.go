package hchecker

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"
)

func init() {
	checkScheme = ""
}

func TestHealthChecker(t *testing.T) {
	srvA := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srvB := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
	}))
	srvC := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(20 * time.Millisecond)
		w.WriteHeader(http.StatusNotFound)
	}))

	c := New()
	c.Watch([]string{srvA.URL, srvB.URL, srvC.URL}, time.Second)

	time.Sleep(100 * time.Millisecond)

	t.Run("Status", func(t *testing.T) {
		status := c.Status()
		sort.Slice(status, func(i, j int) bool { return status[i].Respond < status[j].Respond })
		require.Len(t, status, 3)

		require.Equal(t, srvA.URL, status[0].Site)
		require.True(t, status[0].Live)
		require.True(t, status[0].Respond >= 0)

		require.Equal(t, srvB.URL, status[1].Site)
		require.True(t, status[1].Live)
		require.True(t, status[1].Respond >= 10e-3)

		require.Equal(t, srvC.URL, status[2].Site)
		require.False(t, status[2].Live)
		require.True(t, status[2].Respond >= 20e-3)
	})

	t.Run("StatusOf", func(t *testing.T) {
		status := c.StatusOf(srvB.URL)
		require.Equal(t, srvB.URL, status.Site)

		t.Run("when unknown site", func(t *testing.T) {
			status := c.StatusOf("unknown")
			require.Equal(t, Status{}, status)
		})
	})

	t.Run("FindMin", func(t *testing.T) {
		status := c.FindMin()
		require.Equal(t, srvA.URL, status.Site)
	})

	t.Run("FindMax", func(t *testing.T) {
		status := c.FindMax()
		require.Equal(t, srvB.URL, status.Site)
	})
}

func TestCheckSite(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	respond, err := checkSite(srv.URL, time.Second)
	require.NoError(t, err)
	require.True(t, respond > 0)

	t.Run("when timeout", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Second)
		}))

		_, err := checkSite(srv.URL, 100*time.Millisecond)
		require.Error(t, err)
	})
}
