package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRoutes(t *testing.T) {
	srv := httptest.NewServer(LoPhiRouter())
	client := srv.Client()
	t.Run("/spotify/auth redirects the user to Spotify", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/spotify/auth", nil)
		if err != nil {
			log.Println(err)
		}
		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}
		expect := "https://accounts.spotify.com/"
		receive := res.Request.URL.String()
		if !strings.HasPrefix(receive, expect) {
			t.Errorf("Incorrect URL! Redirected to %q, Expected %q", receive, expect)
		}
	})
	srv.Close()
}
